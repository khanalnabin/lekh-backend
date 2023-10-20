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
	"go.mongodb.org/mongo-driver/mongo"
)

func PostsFeed(c *fiber.Ctx) error {

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
	userID := claims.UserID

	db, err := database.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	coll := db.Collection("posts")

	var user models.User
	err = db.Collection("users").FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	following := make([]string, 0)
	for _, val := range user.Following {
		following = append(following, val.Hex())
	}
	var posts []models.Post
	filter := bson.M{"$or": bson.A{bson.M{"created_by": userID}, bson.M{"visibility": 0}, bson.M{"$and": bson.A{bson.M{"visibility": 2}, bson.M{"created_by": bson.M{"$in": following}}}}}}
	cursor, err := coll.Find(context.Background(), filter)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if err = cursor.All(context.Background(), &posts); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "",
		"posts": posts,
	})
}
func PostsOwn(c *fiber.Ctx) error {

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
	userID := claims.UserID

	db, err := database.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	coll := db.Collection("posts")
	var posts []models.Post
	cursor, err := coll.Find(context.Background(), bson.M{"created_by": userID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	if err = cursor.All(context.Background(), &posts); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "",
		"posts": posts,
	})
}

func PostNew(c *fiber.Ctx) error {
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
	userID := claims.UserID
	post := &models.Post{}

	if err := c.BodyParser(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	post.Creator = userID
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	post.Creator = userID
	post.ID = primitive.NewObjectID()

	db, err := database.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	coll := db.Collection("posts")
	if _, err = coll.InsertOne(context.Background(), *post); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"post":  *post,
	})
}

func PostEdit(c *fiber.Ctx) error {

	postID := c.Params("id")
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
	userID := claims.UserID
	post := &models.Post{}

	if err := c.BodyParser(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	post.UpdatedAt = time.Now()

	db, err := database.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	defer db.Client().Disconnect(context.Background())

	coll := db.Collection("posts")
	objectID, _ := primitive.ObjectIDFromHex(postID)
	updatedPost := bson.M{"updated_at": post.UpdatedAt, "content": post.Content, "image_link": post.Image, "visibility": post.PostVisibility}
	var p models.Post
	coll.FindOneAndUpdate(context.Background(), bson.M{"_id": objectID, "created_by": userID}, bson.M{"$set": updatedPost}).Decode(&p)
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"post":  p,
	})
}

func PostGet(c *fiber.Ctx) error {
	postID := c.Params("id")
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
	coll := db.Collection("posts")
	objectID, _ := primitive.ObjectIDFromHex(postID)

	var post models.Post

	err = coll.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&post)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "the post does not exist",
		})
	}
	if post.Creator == claims.UserID || post.PostVisibility == models.PUBLIC || (post.PostVisibility == models.FOLLOWERS && isFollowing(post.Creator, claims.UserID, db.Collection("users"))) {
		return c.JSON(fiber.Map{
			"error": false,
			"msg":   nil,
			"post":  post,
		})
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": true,
		"msg":   "the user is unauthorized to view the post",
	})
}

func PostDelete(c *fiber.Ctx) error {

	postIDString := c.Params("id")
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
	defer db.Client().Disconnect(context.Background())

	coll := db.Collection("posts")
	postID, _ := primitive.ObjectIDFromHex(postIDString)

	var p models.Post
	err = coll.FindOneAndDelete(context.Background(), bson.M{"_id": postID, "created_by": claims.UserID}).Decode(&p)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "post not found",
		})
	}
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "",
		"post":  p,
	})
}

func isFollowing(creator primitive.ObjectID, user primitive.ObjectID, coll *mongo.Collection) bool {
	var c models.User
	err := coll.FindOne(context.Background(), bson.M{"_id": creator}).Decode(&c)
	if err != nil {
		return false
	}
	for _, follower := range c.Followers {
		if follower == user {
			return true
		}
	}
	return false
}
