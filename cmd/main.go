package main

import (
	"fmt"
	"github.com/dkaman/cogs/internal/config"
)

func main() {
	c, err := config.New(
		config.WithJSONConfigFile("./config.json"),
		config.WithEnvVars(),
	)
	if err != nil {
		fmt.Printf("error with config: %s", err)
		return
	}

	var o int
	c.Get("option", &o)

	fmt.Printf("c: %v\n", o)
}
