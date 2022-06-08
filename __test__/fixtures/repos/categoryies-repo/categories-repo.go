package fixtures_categoryies_repo

import (
	"github.com/google/uuid"
	"github.com/mecamon/shoppingify-server/models"
	"time"
)

var User = models.User{
	Name:      "User for category",
	Lastname:  "Lastname",
	Email:     "some@email.com",
	Password:  "PasswordForRegCat1234",
	IsActive:  true,
	IsVisitor: false,
	LoginCode: uuid.NewString(),
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var UserForSearchCat = models.User{
	Name:      "User for search",
	Lastname:  "Lastname",
	Email:     "search@email.com",
	Password:  "PasswordForSearchCat1234",
	IsActive:  true,
	IsVisitor: false,
	LoginCode: uuid.NewString(),
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var UserForGetAll = models.User{
	Name:      "User for get all",
	Lastname:  "Lastname",
	Email:     "get@all.com",
	Password:  "PasswordValid123",
	IsActive:  true,
	IsVisitor: false,
	LoginCode: uuid.NewString(),
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var UserForGetAllWithItemName = models.User{
	Name:      "User for get with item",
	Lastname:  "Lastname",
	Email:     "get@withItem.com",
	Password:  "PasswordValid123",
	IsActive:  true,
	IsVisitor: false,
	LoginCode: uuid.NewString(),
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var Cat1 = models.Category{
	Name:      "Repo-Cat-1",
	UserID:    0,
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var Cat2 = models.Category{
	Name:      "Repo-Cat2",
	UserID:    0,
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var Cat1ForItemsName = models.Category{
	Name:      "Repo-Fruits",
	UserID:    0,
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var Cat2ForItemsName = models.Category{
	Name:      "Repo-Meats",
	UserID:    0,
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var Item1ForItemsName = models.Item{
	Name:       "Item1 for Items Name",
	Note:       "One note",
	CategoryID: 0,
	ImageURL:   "",
	CreatedAt:  time.Now().Unix(),
	UpdatedAt:  time.Now().Unix(),
}

var Item2ForItemsName = models.Item{
	Name:       "Item2 for Items Name",
	Note:       "One note",
	CategoryID: 0,
	ImageURL:   "",
	CreatedAt:  time.Now().Unix(),
	UpdatedAt:  time.Now().Unix(),
}
