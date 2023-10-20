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

func UserProfile(c *fiber.Ctx) error {
	return nil
}

func UserFollow(c *fiber.Ctx) error {

	friendID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	now := time.Now().Unix()

	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	expires := claims.Expires

	if now > expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check the expiraton time of your token",
		})
	}

	db, err := database.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	userID := claims.UserID

	coll := db.Collection("users")
	var user models.User
	var friend models.User
	err = coll.FindOneAndUpdate(context.Background(), bson.M{"_id": friendID}, bson.M{"$push": bson.M{"followers": userID}}).Decode(&friend)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	err = coll.FindOneAndUpdate(context.Background(), bson.M{"_id": userID}, bson.M{"$push": bson.M{"following": friendID}}).Decode(&user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"error":  false,
		"user":   user,
		"friend": friend,
	})
}
func UserUnfollow(c *fiber.Ctx) error {

	friendID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	now := time.Now().Unix()

	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	expires := claims.Expires

	if now > expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check the expiraton time of your token",
		})
	}

	db, err := database.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	userID := claims.UserID

	coll := db.Collection("users")
	var user models.User
	var friend models.User
	err = coll.FindOneAndUpdate(context.Background(), bson.M{"_id": friendID}, bson.M{"$pull": bson.M{"followers": userID}}).Decode(&friend)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	err = coll.FindOneAndUpdate(context.Background(), bson.M{"_id": userID}, bson.M{"$pull": bson.M{"following": friendID}}).Decode(&user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"error":  false,
		"user":   user,
		"friend": friend,
	})
}

func UsersFollowing(c *fiber.Ctx) error { return nil }

func UsersFollowers(c *fiber.Ctx) error { return nil }

func Users(c *fiber.Ctx) error { return nil }
