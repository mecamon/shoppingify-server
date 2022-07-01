package fixtures_items_repo

import (
	"github.com/google/uuid"
	"github.com/mecamon/shoppingify-server/models"
	"time"
)

var User = models.User{
	Name:      "Items test user",
	Lastname:  "Items test user",
	Email:     "itemsrepo@user.com",
	Password:  "PasswordValid0987",
	IsActive:  true,
	IsVisitor: false,
	LoginCode: uuid.NewString(),
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var UserForGetAllByID = models.User{
	Name:      "Get all by catID",
	Lastname:  "Lastname",
	Email:     "getall@bycatid.com",
	Password:  "PasswordValid0987",
	IsActive:  true,
	IsVisitor: false,
	LoginCode: uuid.NewString(),
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var Cat = models.Category{
	Name:      "Cat For Item Reg",
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var CatForGetByID = models.Category{
	Name:      "Cat Get by Id",
	UserID:    0,
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var Item1 = models.Item{
	Name:      "Item 1 item repo",
	Note:      "Item one note",
	ImageURL:  "",
	IsActive:  true,
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var Item2 = models.Item{
	Name:      "Item 2 item repo",
	Note:      "Item two note",
	ImageURL:  "",
	IsActive:  true,
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}

var GenericItem = models.Item{
	Name:      "Item",
	Note:      "Item one note",
	ImageURL:  "",
	IsActive:  true,
	CreatedAt: time.Now().Unix(),
	UpdatedAt: time.Now().Unix(),
}
