package routes

import (
	"kitten-server/api/controllers"
	"kitten-server/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupHTTPRoutes(app *fiber.App) {
	api := app.Group("/api/v1");

	api.Get("/username/:username", controllers.CheckUserNameExistsModel)
	api.Post("/register", controllers.CreateUserName);
	api.Post("/login",controllers.LoginUserName)

	api.Use(middleware.GetJWTConfig());

	api.Get("/details", controllers.GetUserDetailsModel);
	api.Post("/start", controllers.CreateGameBoard);
	api.Patch("/status/:id", controllers.StoreGameMovesModel);
	api.Get("/userGames",controllers.GetUserGamesModel);
	api.Get("/game/:id", controllers.GetGameByIdModel);
	api.Get("/leaderboard",controllers.GetLeaderBoardModel)
}