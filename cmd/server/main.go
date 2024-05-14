package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func Run() error {
	fmt.Println("Starting up our application...")
	app := fiber.New()
	app.Get("/", greet)

	app.Listen(":3000")
	return nil
}

func main() {

	if err := Run(); err != nil {
		fmt.Println(err)
	}
}

func greet(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
