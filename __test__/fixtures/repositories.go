package fixtures

import (
	"github.com/mecamon/shoppingify-server/models"
)

var UserForRegisterCat = models.User{
	Name:     "User for category",
	Lastname: "Lastname",
	Email:    "some@email.com",
	Password: "PasswordForRegCat1234",
}

var UserForSearchCat = models.User{
	Name:     "User for search",
	Lastname: "Lastname",
	Email:    "search@email.com",
	Password: "PasswordForSearchCat1234",
}
