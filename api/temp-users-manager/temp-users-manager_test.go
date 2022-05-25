//go:build !integration
// +build !integration

package temp_users_manager

import (
	"github.com/mecamon/shoppingify-server/models"
	"testing"
)

func TestTempUsersManager_Tracker(t *testing.T) {
	crtlChan := make(chan int)
	tmpUsers := TempUsersManager{
		UsersCh:             make(chan models.User),
		TotalOfUserDisabled: []models.User{},
	}

	usersSlice := []models.User{
		{Name: "Carlos"},
		//{Name: "Alberto"},
		//{Name: "Alberto"},
	}

	go tmpUsers.Tracker(crtlChan)

	for _, user := range usersSlice {
		tmpUsers.UsersCh <- user
	}

	if len(tmpUsers.TotalOfUserDisabled) != len(usersSlice) {
		t.Errorf("total of users lenght is %d and it should be %d", len(tmpUsers.TotalOfUserDisabled), len(usersSlice))
	}

}
