package routes

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/nabinkhanal/lekh-backend/app/controllers"
	"gitlab.com/nabinkhanal/lekh-backend/pkg/middlewares"
)

func PrivateRoutes(a *fiber.App) {
	route := a.Group("/api/v1")
	route.Post("/auth/logout", middlewares.JWTProtected(), controllers.UserLogout)

	route.Get("/user/profile/:id", middlewares.JWTProtected(), controllers.UserProfile)

	route.Get("/users/followers", middlewares.JWTProtected(), controllers.UsersFollowers)
	route.Get("/users/following", middlewares.JWTProtected(), controllers.UsersFollowing)
	route.Get("/users", middlewares.JWTProtected(), controllers.Users)

	route.Post("/user/follow/:id", middlewares.JWTProtected(), controllers.UserFollow)
	route.Post("/user/unfollow/:id", middlewares.JWTProtected(), controllers.UserUnfollow)

	route.Get("/posts/feed", middlewares.JWTProtected(), controllers.PostsFeed)
	route.Get("/posts/all", middlewares.JWTProtected(), controllers.PostsOwn)
	route.Get("/post/:id", middlewares.JWTProtected(), controllers.PostGet)
	route.Put("/post/:id", middlewares.JWTProtected(), controllers.PostEdit)
	route.Delete("/post/:id", middlewares.JWTProtected(), controllers.PostDelete)
	route.Post("/post", middlewares.JWTProtected(), controllers.PostNew)
}
