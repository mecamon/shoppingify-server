package temp_users_manager

import (
	"github.com/mecamon/shoppingify-server/models"
)

type TempUsersManager struct {
	UsersCh             chan models.User
	TotalOfUserDisabled []models.User
}

func (t *TempUsersManager) Tracker(ctrl chan int) {
	for c := range t.UsersCh {
		t.disable(c)
	}
	ctrl <- 1
	close(ctrl)
}

func (t *TempUsersManager) disable(user models.User) {
	t.TotalOfUserDisabled = append(t.TotalOfUserDisabled, user)
}
