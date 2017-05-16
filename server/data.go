package server

import (
	"time"
)

var Users = map[string]*User{
	"1": {"1", "Adam", "_@anh.io", time.Now()},
}

var nextUser = 2

// GetUser returns our fake user.
func GetUser(id string) *User {
	if user, ok := Users[id]; ok {
		return user
	}
	return nil
}
