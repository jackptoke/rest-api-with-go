package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()
	app.Get("/", greet)

	app.Listen(":3000")
}

func greet(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
