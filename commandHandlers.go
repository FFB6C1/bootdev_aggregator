package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ffb6c1/bootdev_aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		err := fmt.Errorf("Login: Insufficient arguments.")
		return err
	}

	if _, err := s.db.GetUser(context.Background(), cmd.arguments[0]); err != nil {
		err := fmt.Errorf("Login: User does not exist.")
		return err
	}

	err := s.config.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Println("Welcome, " + cmd.arguments[0] + "! You're now logged in.")
	return nil

}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		err := fmt.Errorf("Register: Please include a username.")
		return err
	}

	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	}
	if _, err := s.db.CreateUser(context.Background(), args); err != nil {
		return err
	}
	fmt.Println("User successfully created: " + cmd.arguments[0])
	return handlerLogin(s, cmd)
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.Reset(context.Background()); err != nil {
		return err
	}
	fmt.Println("Database reset.")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Println("handlerList err:", err)
		os.Exit(1)
	}
	for _, user := range users {
		if user.Name == s.config.CurrentUserName {
			fmt.Println("* " + user.Name + " (current)")
		} else {
			fmt.Println("* " + user.Name)
		}
	}
	return nil
}
