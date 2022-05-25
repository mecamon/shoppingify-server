package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	IsActive  bool   `json:"is_active"`
	IsVisitor bool   `json:"is_visitor"`
	LoginCode string `json:"login_code"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ErrorMap map[string]string

type CustomClaims struct {
	TokenType string
	ID        int `json:"id"`
	*jwt.RegisteredClaims
}
