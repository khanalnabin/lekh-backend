package main

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"gitlab.com/nabinkhanal/lekh-backend/pkg/configs"
	"gitlab.com/nabinkhanal/lekh-backend/pkg/middlewares"
	"gitlab.com/nabinkhanal/lekh-backend/pkg/routes"
	"gitlab.com/nabinkhanal/lekh-backend/pkg/utils"
)

func main() {

	config := configs.FiberConfig()
	app := fiber.New(config)
	middlewares.FiberMiddleware(app)
	routes.PublicRoutes(app)
	routes.PrivateRoutes(app)
	routes.NotFoundRoute(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})
	utils.Start(app)
}
