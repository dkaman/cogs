package collection

import (
	"context"
	"fmt"
	"io"

	"github.com/dkaman/cogs/internal/commands"
	"github.com/dkaman/cogs/internal/config"

	"github.com/dkaman/discogs-golang"
	"github.com/dkaman/reg"
	"github.com/olekukonko/tablewriter"
	flag "github.com/spf13/pflag"
)

func init() {
	reg.Register[commands.Command]("list-folders", New())
}

var (
	FlagUsername string
)

type listFoldersCommand struct {
	flags          *flag.FlagSet
	requiredParams []string
	folders        []discogs.Folder
}

func New() (cmd *listFoldersCommand) {
	cmd = &listFoldersCommand{
		flags: flag.NewFlagSet("collection", flag.ContinueOnError),
		requiredParams: []string{
			commands.GlobalUser,
		},
	}

	cmd.flags.StringVar(&FlagUsername, "username", "", "the username that owns the target collection")

	return
}

func (c *listFoldersCommand) Configure(app *commands.App, args []string) (err error) {
	err = c.flags.Parse(args)
	if err != nil {
		return
	}

	err = app.Config.Merge(config.WithFlags(c.flags))
	if err != nil {
		return
	}

	return app.Config.Validate(c.requiredParams)
}

func (c *listFoldersCommand) Run(app *commands.App) (err error) {
	var username string

	err = app.Config.Get(commands.GlobalUser, &username)
	if err != nil {
		return
	}

	folders, err := app.Client.Collection.ListFolders(context.TODO(), username)
	if err != nil {
		return
	}

	c.folders = folders

	return
}

func (c *listFoldersCommand) Print(app *commands.App, w io.Writer) (err error) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"name", "id", "count", "resource-url"})

	for _, f := range c.folders {
		idString := fmt.Sprintf("%d", f.ID)
		countString := fmt.Sprintf("%d", f.Count)
		table.Append([]string{f.Name, idString, countString, f.ResourceURL})
	}

	table.Render()

	return
}
