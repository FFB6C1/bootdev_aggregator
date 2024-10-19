package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ffb6c1/bootdev_aggregator/internal/database"
)

func checkLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	returnedFunc := func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			fmt.Println("Not currently logged in. Please use 'login [username]' to log in.")
			os.Exit(1)
		}
		return handler(s, cmd, user)
	}
	return returnedFunc
}
