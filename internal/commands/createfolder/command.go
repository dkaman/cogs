package createfolder

import (
	"context"
	"fmt"
	"io"

	"github.com/dkaman/cogs/internal/commands"

	"github.com/dkaman/discogs-golang"
	"github.com/dkaman/reg"
	"github.com/olekukonko/tablewriter"
)

func init() {
	reg.Register[commands.Command]("create-folder", New())
}

type folderCreateResult struct {
	name   string
	folder *discogs.Folder
	err    error
}

type createFolderCommand struct {
	requiredParams []string
	targets        []string
	results        []folderCreateResult
}

func New() (cmd *createFolderCommand) {
	cmd = &createFolderCommand{
		requiredParams: []string{
			commands.GlobalUser,
		},
	}

	return
}

func (c *createFolderCommand) Configure(app *commands.App, args []string) (err error) {
	if len(args) == 0 {
		err = fmt.Errorf("create-folder called with no arguments")
		return
	}

	c.targets = args

	return app.Config.Validate(c.requiredParams)
}

func (c *createFolderCommand) Run(app *commands.App) (err error) {
	var username string

	err = app.Config.Get(commands.GlobalUser, &username)
	if err != nil {
		return
	}

	for _, t := range c.targets {
		folder, err := app.Client.Collection.CreateFolder(context.TODO(), username, t)
		c.results = append(c.results, folderCreateResult{
			name:   t,
			folder: folder,
			err:    err,
		})
	}

	return
}

func (c *createFolderCommand) Print(app *commands.App, w io.Writer) (err error) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"name", "id", "count", "resource-url", "error"})

	for _, r := range c.results {
		if r.err == nil {
			idString := fmt.Sprintf("%d", r.folder.ID)
			countString := fmt.Sprintf("%d", r.folder.Count)
			table.Append([]string{r.folder.Name, idString, countString, r.folder.ResourceURL, ""})
		} else {
			table.Append([]string{r.name, "", "", "", r.err.Error()})
		}
	}

	table.Render()

	return
}
