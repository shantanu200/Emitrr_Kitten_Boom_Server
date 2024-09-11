package main

import (
	"kitten-server/api/routes"
	"kitten-server/internals"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	go internals.InitRedis()
	app := fiber.New();

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, https://66e1f2178181b869a2888abf--peaceful-boba-ad5cb8.netlify.app",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Content-Type, Authorization, Origin, Accept",
		AllowCredentials: true,
	}))
	app.Use(healthcheck.New())
	app.Use(compress.New())
	app.Use(logger.New())
	
	routes.SetupHTTPRoutes(app);

	log.Println("Server started on port 9000");
	if err := app.Listen(":9000"); err != nil {
		panic(err)
	}
}