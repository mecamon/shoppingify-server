package models

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
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

type ItemDTO struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Note       string `json:"note"`
	ImageURL   string `json:"image_url"`
	CategoryID int64  `json:"category_id"`
}

type CategoriesGroup struct {
	CategoryID   int64     `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Items        []ItemDTO `json:"items"`
}

type Category struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	UserID    int64  `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type CategoryDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type List struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	IsCompleted bool   `json:"is_completed"`
	IsCancelled bool   `json:"is_cancelled"`
	UserID      int64  `json:"user_id"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	CompletedAt int64  `json:"completed_at"`
}

type ListDTO struct {
	ID          int64             `json:"id"`
	Name        string            `json:"name"`
	Date        time.Time         `json:"date"`
	Items       []SelectedItemDTO `json:"items"`
	IsCompleted bool              `json:"is_completed"`
	IsCancelled bool              `json:"is_cancelled"`
}

type OldListDTO struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	IsCompleted bool      `json:"is_completed"`
	IsCancelled bool      `json:"is_cancelled"`
}

type SelectedItem struct {
	ID          int64 `json:"id"`
	ItemID      int64 `json:"item_id"`
	Quantity    int   `json:"quantity"`
	IsCompleted bool  `json:"is_completed"`
	ListID      int64 `json:"list_id"`
	CreatedAt   int64 `json:"created_at"`
	UpdatedAt   int64 `json:"updated_at"`
}

type UpdateSelItemDTO struct {
	ItemID   int64 `json:"item_id"`
	Quantity int   `json:"quantity"`
}

type SelectedItemDTO struct {
	ID       int64  `json:"id"`
	ItemID   int64  `json:"item_id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type TopCategory struct {
	ID          int64 `json:"id"`
	UserID      int64 `json:"user_id"`
	CategoryID  int64 `json:"category_id"`
	SumQuantity int   `json:"sum_quantity"`
}

type TopCategoryDTO struct {
	ID          int64  `json:"id"`
	CategoryID  int64  `json:"category_id"`
	Name        string `json:"name"`
	SumQuantity int    `json:"sum_quantity"`
	Percentage  int    `json:"percentage"`
}

type TopItem struct {
	ID          int64 `json:"id"`
	UserID      int64 `json:"user_id"`
	ItemID      int64 `json:"item_id"`
	SumQuantity int   `json:"sum_quantity"`
}

type TopItemDTO struct {
	ID          int64  `json:"id"`
	ItemID      int64  `json:"item_id"`
	Name        string `json:"name"`
	SumQuantity int    `json:"sum_quantity"`
	Percentage  int    `json:"percentage"`
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
