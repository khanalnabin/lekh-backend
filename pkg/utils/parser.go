package utils

import (
	"errors"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenMetadata struct {
	UserID  primitive.ObjectID
	Expires int64
}

func extractToken(c *fiber.Ctx) string {
	bearToken := c.Get("Authorization")
	ts := strings.Split(bearToken, " ")
	if len(ts) == 2 {
		return ts[1]
	}
	return ""
}

func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := extractToken(c)
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractTokenMetadata(c *fiber.Ctx) (*TokenMetadata, error) {
	token, err := verifyToken(c)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID, ok := claims["id"].(string)
		if !ok {
			return nil, errors.New("invalid token")
		}
		objectID, _ := primitive.ObjectIDFromHex(userID)
		expires := int64(claims["expires"].(float64))

		return &TokenMetadata{
			UserID:  objectID,
			Expires: expires,
		}, nil
	}
	return nil, err
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}
