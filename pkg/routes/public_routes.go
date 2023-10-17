package routes

import (
	"gitlab.com/nabinkhanal/lekh-backend/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App){
	route := a.Group("/api/v1")

	route.Post("/user/register", controllers.UserRegister)
	route.Post("/user/login", controllers.UserLogin)

}
