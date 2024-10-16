package main

import (
	"fmt"

	"github.com/ffb6c1/bootdev_aggregator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name      string
	arguments []string
}

// Commands struct and methods
type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if _, ok := c.cmds[cmd.name]; !ok {
		err := fmt.Errorf("%s is not a valid command.", cmd.name)
		return err
	}
	err := c.cmds[cmd.name](s, cmd)
	return err
}
