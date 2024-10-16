package main

func getCommands() map[string]func(*state, command) error {
	return map[string]func(*state, command) error{
		"login": handlerLogin,
	}
}
