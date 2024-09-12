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

	if err != redis.Nil {
		return false, nil
	}

	return exists != 0, nil
}

func LoginUserName(userName string, password string) (string, error) {
	pipe := internals.RDB.Pipeline()
	_userName := pipe.Exists(context.TODO(), userName)
	_password := pipe.HGet(context.TODO(), userName, "password")

	_, err := pipe.Exec(context.TODO())

	if err != nil {
		return "", err
	}

	if _userName.Val() == 0 {
		return "", errors.New("user not found")
	}

	if _password.Val() != password {
		return "", errors.New("invalid password")
	}

	token, err := utils.GenerateAccessToken(userName)

	if err != nil {
		return "", err
	}

	return token, nil
}

func CreateUserNameHandler(userName string, password string) (string, error) {
	pipe := internals.RDB.Pipeline();

	pipe.HSet(context.TODO(), userName, "username", userName)
	pipe.HSet(context.TODO(), userName, "password", password)
	pipe.HSet(context.TODO(), userName, "totalGamePlayed", 0)
	pipe.HSet(context.TODO(), userName, "totalGameWon", 0)
	pipe.HSet(context.TODO(), userName, "totalGameLost", 0)
	pipe.HSet(context.TODO(), userName, "createdAt", time.Now().Format(time.RFC3339))

	_,err := pipe.Exec(context.TODO());

	if err != nil {
		return "",err
	}

	token, err := utils.GenerateAccessToken(userName)

	if err != nil {
		return "", err
	}

	return token, nil

}

func GetUserDetails(userName string) (*UserNamePayload, error) {
	pipe := internals.RDB.Pipeline()

	userCmd := pipe.HGetAll(context.TODO(), userName)
	
	_, err := pipe.Exec(context.TODO())

	if err != nil {
		return nil, err
	}

	user := userCmd.Val()
	if len(user) == 0 {
		return nil, errors.New("user not exists")
	}

	rankCmd, err := GetUserRank(userName);

	if err != nil {
		return nil, err
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



	userDetails := UserNamePayload{
		UserName:        user["username"],
		Password:        user["password"],
		TotalGamePlayed: totalGamesPlayed,
		TotalGameWon:    totalGameWon,
		TotalGameLost:   totalGameLost,
		LeaderBoardRank: rankCmd,
		CreatedAt:       user["createdAt"],
	}

	return &userDetails, nil
}