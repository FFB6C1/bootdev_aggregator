package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ffb6c1/bootdev_aggregator/internal/config"
)

func main() {
	userConfig, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	userState := state{
		config: &userConfig,
	}
	cmds := commands{
		cmds: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)
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
