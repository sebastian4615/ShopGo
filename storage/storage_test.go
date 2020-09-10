package storage

import (
	"testing"

	uuid "github.com/nu7hatch/gouuid"
)

func TestAddDeleteUser(t *testing.T) {
	randomName, err := uuid.NewV4()
	if err != nil {
		t.Error("Cannot generate user name!")
	}
	randomNameStr := randomName.String()
	_, err = AddNewUser(randomNameStr, "")
	if err == nil {
		t.Error("Create user with empty password pass")
	}
	newPass := "abc"
	newUser, err := AddNewUser(randomNameStr, newPass)
	if err != nil {
		t.Error("Create user failed")
	}
	defer DeleteUser(newUser.Id)
	if newUser.Name != randomNameStr || newUser.Password != newPass {
		t.Error("Create user failed")
	}
	userFromDb, err := GetUserByName(newUser.Name)
	if err != nil {
		t.Error("Get new user from DB fail")
	}
	if userFromDb.Id != newUser.Id || userFromDb.Password != newUser.Password || userFromDb.Name != newUser.Name {
		t.Error("User DB is different than created")
	}
}
