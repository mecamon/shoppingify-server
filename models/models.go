package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID        int64  `json:"id"`
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

type Item struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Note       string `json:"note"`
	CategoryID int64  `json:"category_id"`
	ImageURL   string `json:"image_url"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

type Category struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	UserID    int64  `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type CategoryDTO struct {
	ID   int64
	Name string
}

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ErrorMap map[string]string

type CustomClaims struct {
	TokenType string
	ID        int64 `json:"id"`
	*jwt.RegisteredClaims
}
