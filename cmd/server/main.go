package main

import (
	"fibre/internal/comment"
	"fibre/internal/db"
	_ "fibre/internal/db"
	"fibre/internal/transport/http"
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "golang.org/x/crypto/sha3"
	_ "golang.org/x/text/language"
)

func Run() error {
	fmt.Println("Starting up our application...")

	store, err := db.NewDatabase()

	if err != nil {
		msg := fmt.Errorf("could not connect to database: %w", err)
		fmt.Println(msg)
		return msg
	}

	//if err := store.MigrateDB(); err != nil {
	//	return err
	//}

	//if err := store.Ping(context.Background()); err != nil {
	//	return err
	//}

	cmtService := comment.NewService(store)

	httpHandler := http.NewHandler(cmtService)

	if err := httpHandler.Serve(); err != nil {
		return err
	}

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

//func getComment(c *fiber.Ctx) error {
//	return c.JSON()
//}
