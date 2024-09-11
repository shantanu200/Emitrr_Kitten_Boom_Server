package redisfunction

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"kitten-server/internals"
	"time"

	"github.com/google/uuid"
)

type GameBoard struct {
	ID          string   `json:"id"`
	Moves       []string `json:"moves"`
	Deck        []string `json:"deck"`
	Status      string   `json:"status"`
	DefuseCount int64    `json:"defuseCount"`
	IsGameOver  bool     `json:"isGameOver"`
	CreatedAt   string   `json:"createdAt"`
}

func GameBoardExists(gameKey string) (bool, error) {
	exists, err := internals.RDB.Exists(context.TODO(), gameKey).Result()

	if err != nil || exists == 0 {
		return false, errors.New("gameboard not exists")
	}

	return true, nil
}

func GetGameBoard(gameKey string) (GameBoard, error) {
	game, err := internals.RDB.HGetAll(context.TODO(), gameKey).Result()

	if err != nil {
		return GameBoard{}, err
	}

	var _game = GameBoard{
		ID:        game["Id"],
		Status:    game["Status"],
		CreatedAt: game["CreatedAt"],
	}

	var deck []string

	if err := json.Unmarshal([]byte(game["Deck"]), &deck); err != nil {
		return GameBoard{}, err
	}

	_game.Deck = deck

	var moves []string

	if err := json.Unmarshal([]byte(game["Moves"]), &moves); err != nil {
		return GameBoard{}, err
	}

	_game.Moves = moves

	if game["isGameOver"] == "true" {
		_game.IsGameOver = true
	} else {
		_game.IsGameOver = false
	}

	return _game, nil

}

func CreateGameBoard(userName string) (string, error) {
	parentKey := "game-" + userName

	randomId := uuid.New().String()

	gameKey := "game-" + userName + "-" + randomId

	timeStamp := time.Now().Format(time.RFC3339)

	internals.RDB.HSet(context.TODO(), gameKey, "Id", randomId)
	internals.RDB.HSet(context.TODO(), gameKey, "Moves", "[]")
	internals.RDB.HSet(context.TODO(), gameKey, "Deck", "[]")
	internals.RDB.HSet(context.TODO(), gameKey, "DefuseCount", 0)
	internals.RDB.HSet(context.TODO(), gameKey, "Status", "ONGOING") // Status will be ONGOING, WON, LOST
	internals.RDB.HSet(context.TODO(), gameKey, "CreatedAt", timeStamp)

	if err := internals.RDB.HSet(context.TODO(), parentKey, gameKey, gameKey).Err(); err != nil {
		return "", err
	}

	if err := internals.RDB.HIncrBy(context.TODO(), userName, "totalGamePlayed", 1).Err(); err != nil {
		return "", err
	}

	return randomId, nil
}

func StoreGameMoves(gameKey string, gameBoard GameBoard, userName string) error {
	exists, err := GameBoardExists(gameKey)

	if err != nil || !exists {
		return err
	}

	_deck, err := json.Marshal(gameBoard.Deck)

	if err != nil {
		return err
	}

	_moves, err := json.Marshal(gameBoard.Moves)

	if err != nil {
		return err
	}

	internals.RDB.HSet(context.TODO(), gameKey, "Deck", _deck)
	internals.RDB.HSet(context.TODO(), gameKey, "Moves", _moves)
	internals.RDB.HSet(context.TODO(), gameKey, "DefuseCount", gameBoard.DefuseCount)
	internals.RDB.HSet(context.TODO(), gameKey, "Status", gameBoard.Status)

	fmt.Println(gameBoard.Status)

	if gameBoard.Status == "WON" {
		internals.RDB.HIncrBy(context.TODO(), userName, "totalGameWon", 1)
		internals.RDB.HSet(context.TODO(), gameKey, "IsGameOver", "true")
	} else if gameBoard.Status == "LOST" {
		internals.RDB.HIncrBy(context.TODO(), userName, "totalGameLost", 1)
		internals.RDB.HSet(context.TODO(), gameKey, "IsGameOver", "true")
	}

	if gameBoard.Status != "ONGOING" {
		if err := UpdatePlayerRanking(userName); err != nil {
			return err
		}
	}

	return nil
}

func GetUserGames(userName string) ([]GameBoard, error) {
	key := "game-" + userName

	games, err := internals.RDB.HGetAll(context.TODO(), key).Result()

	if err != nil {
		return nil, err
	}

	var _games []GameBoard

	for _, value := range games {
		_game, err := internals.RDB.HGetAll(context.TODO(), value).Result()

		if err != nil {
			return nil, err
		}

		game := GameBoard{
			ID:        _game["Id"],
			Status:    _game["Status"],
			CreatedAt: _game["CreatedAt"],
		}

		if _game["IsGameOver"] == "true" {
			game.IsGameOver = true
		} else {
			game.IsGameOver = false
		}

		_games = append(_games, game)
	}

	return _games, nil
}
