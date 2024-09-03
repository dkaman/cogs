package collection

import (
	"fmt"
	"io"
	// "log/slog"

	"github.com/dkaman/cogs/internal/commands"
	"github.com/dkaman/cogs/internal/config"
	// "github.com/dkaman/discogs-golang"

	"github.com/dkaman/reg"
	flag "github.com/spf13/pflag"
)

func init() {
	reg.Register[commands.Command]("collection", New())
}

var (
	FlagUsername string
)

type collectionCommand struct {
	flags  *flag.FlagSet
	config *config.Config
}

func New() (cmd *collectionCommand) {
	cmd = &collectionCommand{
		flags: flag.NewFlagSet("collection", flag.ContinueOnError),
	}

	cmd.flags.StringVar(&FlagUsername, "username", "", "the username that owns the target collection")

	return
}

func (c *collectionCommand) Configure(app *commands.App, args []string) (err error) {
	fmt.Printf("hello from configure\n")
	c.config = app.Config
	err = c.flags.Parse(args)
	if err != nil {
		return
	}
	return app.Config.Merge(config.WithFlags(c.flags))
}

func (c *collectionCommand) Run(app *commands.App) (err error) {
	fmt.Printf("hello from run\n")
	fmt.Printf("config: %v\n", app.Config.JSON())
	return
}

func (c *collectionCommand) Print(app *commands.App, w io.Writer) (err error) {
	s := []byte("hello from print\n")

	_, err = w.Write(s)
	if err != nil {
		return err
	}

	return nil
}
