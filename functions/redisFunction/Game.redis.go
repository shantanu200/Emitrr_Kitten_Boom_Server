package redisfunction

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"kitten-server/internals"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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

	if err := internals.RDB.ZAdd(context.TODO(), parentKey, redis.Z{Score: float64(time.Now().Unix()), Member: gameKey}).Err(); err != nil {
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

type PaginationResponse struct {
	TotalPages     int         `json:"totalPages"`
	TotalDocuments int         `json:"totalDocuments"`
	CurrentPage    int         `json:"currentPage"`
	Games          []GameBoard `json:"games"`
}

func GetUserGames(userName string, page int, limit int) (*PaginationResponse, error) {
	key := "game-" + userName
	ctx := context.TODO()

	// Get the total number of documents
	totalDocument := internals.RDB.ZCount(ctx, key, "-inf", "+inf").Val()

	if totalDocument == 0 {
		fmt.Println("No Documents found")
		return &PaginationResponse{
			TotalPages:     0,
			TotalDocuments: 0,
			CurrentPage:    page,
			Games:          []GameBoard{},
		}, nil
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalDocument) / float64(limit)))
	offset := (page - 1) * limit

	// Fetch the game IDs for the current page
	gameIDs, err := internals.RDB.ZRevRange(ctx, key, int64(offset), int64(offset+limit-1)).Result()
	if err != nil {
		return nil, err
	}

	// Use a WaitGroup to handle concurrent requests
	var wg sync.WaitGroup
	var mu sync.Mutex
	var _games []GameBoard
	errCh := make(chan error, len(gameIDs))

	// Fetch each game's details concurrently
	for _, gameID := range gameIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()

			// Fetch game details concurrently
			gameData, err := internals.RDB.HGetAll(ctx, id).Result()
			if err != nil {
				errCh <- err
				return
			}

			// Build GameBoard struct
			game := GameBoard{
				ID:        gameData["Id"],
				Status:    gameData["Status"],
				CreatedAt: gameData["CreatedAt"],
			}

			game.IsGameOver = gameData["IsGameOver"] == "true"

			// Lock the slice while appending to avoid race conditions
			mu.Lock()
			_games = append(_games, game)
			mu.Unlock()
		}(gameID)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errCh)

	// Check if there were any errors during fetching
	if len(errCh) > 0 {
		return nil, <-errCh
	}

	// Return the paginated response
	return &PaginationResponse{
		TotalPages:     totalPages,
		TotalDocuments: int(totalDocument),
		CurrentPage:    page,
		Games:          _games,
	}, nil
}
