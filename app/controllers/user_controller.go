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
	requestID, _ := primitive.ObjectIDFromHex(c.Params("id"))
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
	var requestedUser models.User
	err = db.Collection("users").FindOne(context.Background(), bson.M{"_id": requestID}).Decode(&requestedUser)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	requestedUser.PasswordHash = ""
	var posts []models.Post
	var filter bson.M
	if utils.In(userID, requestedUser.Followers) {
		filter = bson.M{"created_by": requestID, "$or": bson.A{bson.M{"created_by": userID}, bson.M{"visibility": 0}, bson.M{"visibility": 2}}}
	} else {
		filter = bson.M{"created_by": requestID, "$or": bson.A{bson.M{"created_by": userID}, bson.M{"visibility": 0}}}
	}
	cursor, err := db.Collection("posts").Find(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	err = cursor.All(context.Background(), &posts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"error": false,
		"user":  requestedUser,
		"posts": posts,
	})

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
	err = coll.FindOneAndUpdate(context.Background(), bson.M{"_id": friendID}, bson.M{"$addToSet": bson.M{"followers": userID}}).Decode(&friend)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	err = coll.FindOneAndUpdate(context.Background(), bson.M{"_id": userID}, bson.M{"$addToSet": bson.M{"following": friendID}}).Decode(&user)

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

func UsersFollowing(c *fiber.Ctx) error {

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
	var currentUser models.User
	db.Collection("users").FindOne(context.Background(), bson.M{"_id": claims.UserID}).Decode(&currentUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	var users []models.User
	filter := bson.M{"_id": bson.M{"$in": currentUser.Following}}
	cursor, err := db.Collection("users").Find(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	err = cursor.All(context.Background(), &users)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	for i := 0; i < len(users); i++ {
		users[i].PasswordHash = ""
	}
	return c.JSON(fiber.Map{
		"error": false,
		"users": users,
	})
}

func UsersFollowers(c *fiber.Ctx) error {

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
	var currentUser models.User
	db.Collection("users").FindOne(context.Background(), bson.M{"_id": claims.UserID}).Decode(&currentUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	var users []models.User
	filter := bson.M{"_id": bson.M{"$in": currentUser.Followers}}
	cursor, err := db.Collection("users").Find(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	err = cursor.All(context.Background(), &users)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	for i := 0; i < len(users); i++ {
		users[i].PasswordHash = ""
	}
	return c.JSON(fiber.Map{
		"error": false,
		"users": users,
	})
}

func Users(c *fiber.Ctx) error {

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
	var users []models.User
	cursor, err := db.Collection("users").Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	err = cursor.All(context.Background(), &users)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	for i := 0; i < len(users); i++ {
		users[i].PasswordHash = ""
	}
	return c.JSON(fiber.Map{
		"error": false,
		"users": users,
	})
}
