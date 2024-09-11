package request

import "github.com/gofiber/fiber/v2"

func SuccessRouter(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false, "message": message, "data": data})
}

func ErrorRouter(c *fiber.Ctx, message string, err string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": true, "message": message, "err": err})
}

func DocumentCreateRouter(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"error": false, "message": "Document created successfully"})
}

func InvalidBodyParserRouter(c *fiber.Ctx, err string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid Body Passed", "err": err})
}

func ServerErrorRouter(c *fiber.Ctx, err string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": false, "message": "Internal server error", "err": err})
}

func InvalidUserRouter(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": true, "message": "Invalid User | Please login again"})
}
