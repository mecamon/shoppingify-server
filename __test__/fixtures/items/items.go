package fixtures_items

import (
	"github.com/google/uuid"
	"github.com/mecamon/shoppingify-server/models"
	"time"
)

var User = models.User{
	Name:      "user from item",
	Lastname:  "user from item",
	Email:     "userfrom@item.com",
	Password:  "PasswordValid123",
	IsActive:  true,
	IsVisitor: false,
	LoginCode: uuid.NewString(),
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var UserForGet = models.User{
	Name:      "user get",
	Lastname:  "lastname",
	Email:     "userget@item.com",
	Password:  "PasswordValid123",
	IsActive:  true,
	IsVisitor: false,
	LoginCode: uuid.NewString(),
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var UserForGetOne = models.User{
	Name:      "user get one",
	Lastname:  "lastname",
	Email:     "userget@one.com",
	Password:  "PasswordValid123",
	IsActive:  true,
	IsVisitor: false,
	LoginCode: uuid.NewString(),
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var GenericCat = models.Category{
	Name:      "CatForItemsGet",
	UserID:    0,
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var GenericItem = models.Item{
	Name:      "ItemForItemGet",
	Note:      "I am a note",
	ImageURL:  "",
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var Cat = models.Category{
	Name:      "Cat for items handler",
	UserID:    0,
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}
