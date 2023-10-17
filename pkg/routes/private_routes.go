package routes

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/nabinkhanal/lekh-backend/app/controllers"
	"gitlab.com/nabinkhanal/lekh-backend/pkg/middlewares"
)

func PrivateRoutes(a *fiber.App) {
	route := a.Group("/api/v1")
	route.Post("/user/logout", middlewares.JWTProtected(), controllers.UserLogout)

}
