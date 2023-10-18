package routes

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/nabinkhanal/lekh-backend/app/controllers"
)

func PublicRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	route.Post("/auth/register", controllers.UserRegister)
	route.Post("/auth/login", controllers.UserLogin)

}
