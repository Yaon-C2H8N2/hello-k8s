package utils

import (
	"api/internal/models"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func NeedsAuth(c *gin.Context) (*models.User, error) {
	bearer := c.GetHeader("Authorization")
	bearer = bearer[7:]

	token, err := ValidateToken(bearer)
	if err != nil || token == nil {
		if err != nil {
			c.JSON(500, gin.H{
				"error":   "Invalid token",
				"message": err.Error(),
			})
		} else {
			c.JSON(401, gin.H{
				"error":   "Invalid token",
				"message": "Token is not valid",
			})
		}
		return nil, err
	}
	var userId = token.Claims.(jwt.MapClaims)["user_id"]

	conn := GetConnections()
	defer conn.Close(context.Background())

	sql := `
		SELECT users.id, users.first_name, users.last_name, users.age
		FROM users
		WHERE users.id = $1 
	`

	rows := DoRequest(conn, sql, userId)
	if rows.Next() {
		var user = &models.User{}
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Age)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   "Failed to scan user",
				"message": err.Error(),
			})
			return nil, err
		}

		return user, nil
	} else {
		c.JSON(401, gin.H{
			"error":   "Invalid token",
			"message": "User not found",
		})
		return nil, err
	}
}
