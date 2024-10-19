package main

func registerCommands(cmds commands) {
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", checkLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", checkLoggedIn(handlerFollow))
	cmds.register("following", checkLoggedIn(handlerFollowing))
	cmds.register("unfollow", checkLoggedIn(handlerUnfollow))
}
