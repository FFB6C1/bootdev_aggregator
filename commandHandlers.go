package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		err := fmt.Errorf("Login: Insufficient arguments.")
		return err
	}

	err := s.config.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Println("Welcome, " + cmd.arguments[0] + "! You're now logged in.")
	return nil

}
