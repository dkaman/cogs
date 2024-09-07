package main

import (
	"os"
	"fmt"

	"github.com/dkaman/cogs/internal/commands"

	_ "github.com/dkaman/cogs/internal/commands/listfolders"
	_ "github.com/dkaman/cogs/internal/commands/createfolder"
)

func main() {
	if err := commands.Root(os.Args[1:]); err != nil {
		fmt.Printf("error running command: %s\n", err)
		os.Exit(1)
	}
}
