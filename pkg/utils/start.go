package utils

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/nabinkhanal/lekh-backend/platform/database"
)

func Start(a *fiber.App) {
	idleConnClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// server shutdown process
		db := database.DbConn
		db.Collection("users").Indexes().DropAll(context.Background())
		log.Println("Closing Database Connection...")
		db.Client().Disconnect(context.Background())

		if err := a.Shutdown(); err != nil {
			log.Println("Oops... Server is not shutting down! Reason:", err)
		}
		close(idleConnClosed)
	}()
	if err := a.Listen(":" + os.Getenv("SERVER_PORT")); err != nil {
		log.Println("Oops... Server is not running! Reason:", err)
	}
	<-idleConnClosed
}
