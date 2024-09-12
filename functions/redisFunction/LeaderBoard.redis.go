package redisfunction

import (
	"context"
	"kitten-server/internals"
	"math"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func UpdatePlayerRanking(userName string) error {
	pipe := internals.RDB.Pipeline()
	totalWinCmd := pipe.HGet(context.TODO(), userName, "totalGameWon")
	totalGamePlayedCmd := pipe.HGet(context.TODO(), userName, "totalGamePlayed")
	_, err := pipe.Exec(context.TODO())

	if err != nil {
		return err
	}

	totalWin, err := strconv.ParseInt(totalWinCmd.Val(), 10, 64)
	if err != nil {
		return err
	}

	totalGamePlayed, err := strconv.ParseInt(totalGamePlayedCmd.Val(), 10, 64)
	if err != nil {
		return err
	}

	if totalGamePlayed == 0 {
		return nil
	}

	score := (float64(totalWin) / float64(totalGamePlayed)) * 100
	score = math.Round(score*100) / 100 // Round to 2 decimal places

	if err := internals.RDB.ZAdd(context.TODO(), "leaderboard", redis.Z{
		Score:  score,
		Member: userName,
	}).Err(); err != nil {
		return err
	}

	return nil
}

func GetUserRank(userName string) (int64, error) {
	rank, err := internals.RDB.ZRevRank(context.TODO(), "leaderboard", userName).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return rank + 1, nil
}
func GetLeaderBoard() ([]string, error) {
	leaderboard, err := internals.RDB.ZRevRangeByScore(context.TODO(), "leaderboard", &redis.ZRangeBy{
		Min: "0",
		Max: "+inf",
	}).Result()

	if err != nil {
		return nil, err
	}

	return leaderboard, nil
}
