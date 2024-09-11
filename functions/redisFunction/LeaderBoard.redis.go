package redisfunction

import (
	"context"
	"kitten-server/internals"
	"math"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func UpdatePlayerRanking(userName string) error {

	totalWin, err := internals.RDB.HGet(context.TODO(), userName, "totalGameWon").Result()

	if err != nil {
		return err
	}

	totalGamePlayed, err := internals.RDB.HGet(context.TODO(), userName, "totalGamePlayed").Result()

	if err != nil {
		return err
	}

	_totalWin, err := strconv.ParseInt(totalWin, 10, 64)

	if err != nil {
		return err
	}

	_totalGamePlayed, err := strconv.ParseInt(totalGamePlayed, 10, 64)

	if err != nil {
		return err
	}

	score := (float64(_totalWin) / float64(_totalGamePlayed)) * 100

	score = math.Round(score) / 100

	if err := internals.RDB.ZAdd(context.TODO(), "leaderboard", redis.Z{Score: float64(score), Member: userName}).Err(); err != nil {
		return err
	}

	return nil
}

func GetUserRank(userName string) (int64, error) {
	rank, err := internals.RDB.ZRevRank(context.TODO(), "leaderboard", userName).Result()

	if err != nil {
		return 0, err
	}

	return rank + 1, nil
}

func GetLeaderBoard() ([]string,error) {
	leaderboard, err := internals.RDB.ZRevRangeByScore(context.TODO(), "leaderboard", &redis.ZRangeBy{
		Min: "0",
		Max: "+inf",
	}).Result()

	if err != nil {
		return nil,err
	}



	return leaderboard, nil
}
