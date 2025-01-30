package requests

import "api/internal/models"

type RegisterRequest struct {
	models.User
	Username string `json:"username"`
	Password string `json:"password"`
}
