package controllers

import (
	"kitten-server/api/middleware"
	"kitten-server/api/request"
	redisfunction "kitten-server/functions/redisFunction"

	"github.com/gofiber/fiber/v2"
)

func CreateGameBoard(c *fiber.Ctx) error {
	userName, err := middleware.GetUserId(c)

	if userName == "" || err != nil {
		return request.InvalidUserRouter(c)
	}

	gameId, err := redisfunction.CreateGameBoard(userName)

	if err != nil {
		return request.ErrorRouter(c, "Unable to create game", err.Error())
	}

	return request.SuccessRouter(c, "Game created successfully", fiber.Map{"gameId": gameId})
}

func GetGameByIdModel(c *fiber.Ctx) error {

	userName, err := middleware.GetUserId(c)

	if userName == "" || err != nil {
		return request.InvalidUserRouter(c)
	}

	gameKey := c.Params("id")

	if gameKey == "" {
		return request.ErrorRouter(c, "Please provide valid game id", "")
	}

	requestGameKey := "game-" + userName + "-" + gameKey

	result, err := redisfunction.GetGameBoard(requestGameKey)

	if err != nil {
		return request.ErrorRouter(c, "Unable to get game", err.Error())
	}

	return request.SuccessRouter(c, "Game fetched successfully", result)
}

func StoreGameMovesModel(c *fiber.Ctx) error {

	userName, err := middleware.GetUserId(c)

	if userName == "" || err != nil {
		return request.InvalidUserRouter(c)
	}

	gameKey := c.Params("id")

	if gameKey == "" {
		return request.ErrorRouter(c, "Please provide valid game id", "")
	}

	var payload redisfunction.GameBoard

	if err := c.BodyParser(&payload); err != nil {
		return request.InvalidBodyParserRouter(c, err.Error())
	}

	requestId := "game-" + userName + "-" + gameKey

	if err := redisfunction.StoreGameMoves(requestId, payload, userName); err != nil {
		return request.ErrorRouter(c, "Unable to store game moves", err.Error())
	}

	return request.SuccessRouter(c, "Game moves stored successfully", nil)
}

func GetUserGamesModel(c *fiber.Ctx) error {
	userName, err := middleware.GetUserId(c)

	if err != nil || userName == "" {
		return request.InvalidUserRouter(c)
	}

	result, err := redisfunction.GetUserGames(userName)

	if err != nil {
		return request.ErrorRouter(c, "Unable to get user games", err.Error())
	}

	return request.SuccessRouter(c, "User games fetched successfully", result)
}

func GetLeaderBoardModel(c *fiber.Ctx) error {

	result,err := redisfunction.GetLeaderBoard()

	if err != nil {
		return request.ErrorRouter(c, "Unable to get leaderboard", err.Error())
	}

	return request.SuccessRouter(c, "Leaderboard fetched successfully", result);
}
