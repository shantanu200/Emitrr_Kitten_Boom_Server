package controllers

import (
	"kitten-server/api/middleware"
	"kitten-server/api/request"
	redisfunction "kitten-server/functions/redisFunction"

	"github.com/gofiber/fiber/v2"
)

type UserPayload struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func CheckUserNameExistsModel(c *fiber.Ctx) error {
	userName := c.Params("username")

	if userName == "" {
		return request.ErrorRouter(c, "Please provide valid username", "")
	}

	exists, err := redisfunction.CheckUserNameExists(userName)

	if err != nil {
		return request.ErrorRouter(c, "Unable to check user exists", err.Error())
	}

	if exists {
		return request.SuccessRouter(c, "User already exists", fiber.Map{"exists": true})
	}
	return request.SuccessRouter(c, "User does not exists", fiber.Map{"exists": false})
}

func LoginUserName(c *fiber.Ctx) error {
	var payload UserPayload

	if err := c.BodyParser(&payload); err != nil {
		return request.InvalidBodyParserRouter(c, err.Error())
	}

	token, err := redisfunction.LoginUserName(payload.UserName, payload.Password)

	if err != nil {
		return request.ErrorRouter(c, "Unable to login user", err.Error())
	}

	return request.SuccessRouter(c, "User logged in successfully", fiber.Map{"accessToken": token, "userName": payload.UserName})
}

func CreateUserName(c *fiber.Ctx) error {
	var payload UserPayload

	if err := c.BodyParser(&payload); err != nil {
		return request.InvalidBodyParserRouter(c, err.Error())
	}

	token, err := redisfunction.CreateUserNameHandler(payload.UserName, payload.Password)

	if err != nil {
		return request.ErrorRouter(c, "Unable to create user", err.Error())
	}

	return request.SuccessRouter(c, "User logged in successfully", fiber.Map{"accessToken": token, "userName": payload.UserName})
}

func GetUserDetailsModel(c *fiber.Ctx) error {
	userName, err := middleware.GetUserId(c)

	if err != nil {
		return request.InvalidUserRouter(c)
	}

	user, err := redisfunction.GetUserDetails(userName)

	if err != nil {
		return request.ErrorRouter(c, "Unable to get user details", err.Error())
	}

	return request.SuccessRouter(c, "User details fetched successfully", user)
}
