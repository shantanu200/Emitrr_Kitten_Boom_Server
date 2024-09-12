package redisfunction

import (
	"context"
	"kitten-server/internals"
	"math"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func UpdatePlayerRanking(userName string) error {
	ctx := context.TODO()

	pipe := internals.RDB.Pipeline()
	totalWinCmd := pipe.HGet(ctx, userName, "totalGameWon")
	totalGamePlayedCmd := pipe.HGet(ctx, userName, "totalGamePlayed")

	totalGamesPlayedCmd := pipe.Get(ctx, "totalGamesPlayed")
	totalPlayersCmd := pipe.Get(ctx, "totalPlayers")
	_, err := pipe.Exec(ctx)
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

	winRate := float64(totalWin) / float64(totalGamePlayed)

	totalGamesPlayed, err := strconv.ParseInt(totalGamesPlayedCmd.Val(), 10, 64)
	if err != nil {
		return err
	}

	totalPlayers, err := strconv.ParseInt(totalPlayersCmd.Val(), 10, 64)
	if err != nil {
		return err
	}

	averageGamesPlayed := float64(totalGamesPlayed) / float64(totalPlayers)

	weightFactor := float64(totalGamePlayed) / averageGamesPlayed
	if weightFactor > 1 {
		weightFactor = 1
	}

	normalizedScore := winRate * weightFactor
	normalizedScore = math.Round(normalizedScore*100) / 100

	if err := internals.RDB.ZAdd(ctx, "leaderboard", redis.Z{
		Score:  normalizedScore,
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
