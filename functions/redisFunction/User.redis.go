package redisfunction

import (
	"context"
	"errors"
	"kitten-server/api/utils"
	"kitten-server/internals"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserNamePayload struct {
	UserName        string `json:"username"`
	Password        string `json:"password"`
	TotalGamePlayed int64  `json:"totalGamePlayed"`
	TotalGameWon    int64  `json:"totalGameWon"`
	TotalGameLost   int64  `json:"totalGameLost"`
	CreatedAt       string `json:"createdAt"`
	LeaderBoardRank int64  `json:"leaderBoardRank"`
}

func CheckUserNameExists(userName string) (bool, error) {
	exists, err := internals.RDB.Exists(context.TODO(), userName).Result()

	if err != redis.Nil && exists != 0 {
		return true, nil
	}

	return false, nil
}

func LoginUserName(userName string, password string) (string, error) {
	user, err := internals.RDB.Exists(context.TODO(), userName).Result()

	if err != nil || user == 0 {
		return "", errors.New("user not exists")
	}

	userPassword, err := internals.RDB.HGet(context.TODO(), userName, "password").Result()

	if err != nil {
		return "", err
	}

	if userPassword != password {
		return "", errors.New("invalid password")
	}

	token, err := utils.GenerateAccessToken(userName)

	if err != nil {
		return "", err
	}

	return token, nil
}

func CreateUserNameHandler(userName string, password string) (string, error) {
	internals.RDB.HSet(context.TODO(), userName, "username", userName)
	internals.RDB.HSet(context.TODO(), userName, "password", password)
	internals.RDB.HSet(context.TODO(), userName, "totalGamePlayed", 0)
	internals.RDB.HSet(context.TODO(), userName, "totalGameWon", 0)
	internals.RDB.HSet(context.TODO(), userName, "totalGameLost", 0)
	internals.RDB.HSet(context.TODO(), userName, "createdAt", time.Now().Format(time.RFC3339))

	token, err := utils.GenerateAccessToken(userName)

	if err != nil {
		return "", err
	}

	return token, nil

}

func GetUserDetails(userName string) (*UserNamePayload, error) {
	user, err := internals.RDB.HGetAll(context.TODO(), userName).Result()

	if err == redis.Nil || user == nil {
		return nil, errors.New("user not exists")
	}

	totalGamesPlayed, err := strconv.ParseInt(user["totalGamePlayed"], 10, 64)

	if err != nil {
		return nil, err
	}

	totalGameWon, err := strconv.ParseInt(user["totalGameWon"], 10, 64)

	if err != nil {
		return nil, err
	}

	totalGameLost, err := strconv.ParseInt(user["totalGameLost"], 10, 64)

	if err != nil {
		return nil, err
	}

	rank, err := GetUserRank(userName)

	if err != redis.Nil && err != nil {
		return nil, err
	}

	_user := UserNamePayload{
		UserName:        user["username"],
		Password:        user["password"],
		TotalGamePlayed: totalGamesPlayed,
		TotalGameWon:    totalGameWon,
		TotalGameLost:   totalGameLost,
		LeaderBoardRank: rank,
		CreatedAt:       user["createdAt"],
	}

	return &_user, nil
}
