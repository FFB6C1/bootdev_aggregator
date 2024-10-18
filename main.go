package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ffb6c1/bootdev_aggregator/internal/config"
	"github.com/ffb6c1/bootdev_aggregator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	userConfig, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", userConfig.DBurl)
	if err != nil {
		log.Fatal(err)
	}

	userState := state{
		db:     database.New(db),
		config: &userConfig,
	}

	cmds := commands{
		cmds: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	if len(os.Args) < 2 {
		err := fmt.Errorf("Insufficient Arguments. Please use a command name.")
		log.Fatal(err)
	}

	newCommand := command{
		name:      os.Args[1],
		arguments: os.Args[2:],
	}

	err = cmds.run(&userState, newCommand)
	if err != nil {
		log.Fatal(err)
	}

}
