package main

import (
	"fmt"

	"github.com/ffb6c1/bootdev_aggregator/internal/config"
)

func main() {
	userConfig := config.Read()
	userConfig.SetUser("chesca")
	fmt.Println(config.Read())
}
