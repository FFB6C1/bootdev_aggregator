package main

func getCommands() map[string]func(*state, command) error {
	return map[string]func(*state, command) error{
		"login": handlerLogin,
	}
}

func registerCommands(cmds commands) {
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)
}
