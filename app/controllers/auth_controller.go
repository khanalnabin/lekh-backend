package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/nabinkhanal/lekh-backend/app/models"
	"gitlab.com/nabinkhanal/lekh-backend/pkg/utils"
	"gitlab.com/nabinkhanal/lekh-backend/platform/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UserRegister(c *fiber.Ctx) error {
	register := &models.Register{}

	if err := c.BodyParser(register); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	db := database.DbConn
	user := &models.User{}
	user.Id = primitive.NewObjectID()
	user.Name = register.Name
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Email = register.Email
	user.Username = register.Username
	hash, err := utils.GenerateHash(register.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	user.PasswordHash = hash
	_, err = db.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	user.PasswordHash = ""
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"user":  user,
	})
}
func UserLogin(c *fiber.Ctx) error {

	login := &models.Login{}

	if err := c.BodyParser(login); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	db := database.DbConn
	collection := db.Collection("users")

	var user models.User
	filter := bson.M{"username": login.Username}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "user with the given username is not found",
		})
	}

	passwordMatch := utils.ComparePassword(user.PasswordHash, login.Password)
	if !passwordMatch {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "wrong user email or password",
		})
	}
	tokens, err := utils.GenerateNewTokens(user.Id.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"tokens": fiber.Map{
			"access":  tokens.Access,
			"refresh": tokens.Refresh,
		},
	})
}

func UserLogout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Success",
	})
}
