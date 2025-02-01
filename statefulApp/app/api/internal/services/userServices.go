package services

import (
	"api/internal/models"
	"api/internal/models/requests"
	"api/pkg/utils"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
)

func Authenticate(c *gin.Context) {
	authRequestBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "Failed to read request body",
			"message": err.Error(),
		})
		return
	}

	var authRequest = &requests.AuthRequest{}
	err = json.Unmarshal(authRequestBytes, authRequest)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "Failed to unmarshal request body",
			"message": err.Error(),
		})
		return
	}

	//todo : check for jwt in header and validate it

	conn := utils.GetConnections()
	defer conn.Close(context.Background())

	sql := `
		SELECT users.id, users.first_name, users.last_name, users.age
		FROM users
		WHERE login = $1 AND password = $2
	`

	rows := utils.DoRequest(conn, sql, authRequest.Username, authRequest.Password)
	if rows.Next() {
		var user = &models.User{}
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Age)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   "Failed to scan user",
				"message": err.Error(),
			})
			return
		}

		token, err := utils.GenerateToken(*user)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   "Failed to generate token",
				"message": err.Error(),
			})
			return
		}

		c.Header("Set-Cookie", "token=Bearer "+token)
		c.JSON(200, gin.H{
			"authenticated": true,
			"token":         token,
		})

	} else {
		c.JSON(401, gin.H{
			"authenticated": false,
		})
	}
}

func Register(c *gin.Context) {
	registerRequestBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "Failed to read request body",
			"message": err.Error(),
		})
		return
	}

	var registerRequest = &requests.RegisterRequest{}
	err = json.Unmarshal(registerRequestBytes, registerRequest)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "Failed to unmarshal request body",
			"message": err.Error(),
		})
		return
	}

	conn := utils.GetConnections()
	defer conn.Close(context.Background())

	sql := `
		INSERT INTO users (first_name, last_name, age, login, password)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, first_name, last_name, age;
	`

	rows := utils.DoRequest(conn, sql, registerRequest.FirstName, registerRequest.LastName, registerRequest.Age, registerRequest.Username, registerRequest.Password)
	if rows.Next() {
		var user = &models.User{}
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Age)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   "Failed to scan user",
				"message": err.Error(),
			})
			return
		}

		token, err := utils.GenerateToken(*user)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   "Failed to generate token",
				"message": err.Error(),
			})
			return
		}

		c.Header("Set-Cookie", "token=Bearer "+token)
		c.JSON(200, gin.H{
			"authenticated": true,
			"token":         token,
		})

	} else {
		c.JSON(401, gin.H{
			"authenticated": false,
		})
	}
}
