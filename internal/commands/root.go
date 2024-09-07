package commands

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/dkaman/cogs/internal/config"

	discogs "github.com/dkaman/discogs-golang"
	"github.com/dkaman/reg"
)

const (
	GlobalUser = "global.user"
)

var (
	api    *discogs.Client
	env    *config.Config
	logger *slog.Logger

	flagConfigFile string
	flagUser       string

	ErrNoSubcommand = errors.New("no subcommand detected")
)

type Command interface {
	Configure(*App, []string) error
	Run(*App) error
	Print(*App, io.Writer) error
}

type App struct {
	Config *config.Config
	Client *discogs.Client
	Logger *slog.Logger
}

func listSubcommands() {
	fmt.Printf("available subcommands: %v\n", reg.Drivers[Command]())
}

func init() {
	flag.StringVar(&flagConfigFile, "config", config.DefaultConfigPath(), "path to a json config file")
	flag.StringVar(&flagUser, "user", "", "username for all discogs api requests")
	flag.Parse()
}

func initApp() (app *App, err error) {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	env, err = config.New(
		config.WithJSONConfigFile(flagConfigFile),
		config.WithEnvVars(),
	)
	if err != nil {
		return
	}

	if flagUser != "" {
		env.PutString(GlobalUser, flagUser)
	}

	var pat string
	err = env.Get("global.pat", &pat)

	api, err = discogs.NewClient(pat, nil)
	if err != nil {
		logger.Error("discogs api client error", "error", err)
		os.Exit(1)
	}

	app = &App{
		Config: env,
		Client: api,
		Logger: logger,
	}

	return
}

// lol this one's goofy
func initCmdLine() (subcommand string, subcommandArgs []string, err error) {
	args := flag.Args()
	if len(args) < 1 {
		err = ErrNoSubcommand
		return
	}

	if len(args) < 2 {
		subcommand = args[0]
		return
	}

	if len(args) >= 2 {
		subcommand = args[0]
		subcommandArgs = args[1:]
		return
	}

	return
}

func Root(args []string) (err error) {
	app, err := initApp()
	if err != nil {
		return
	}

	subcommand, subcommandArgs, err := initCmdLine()
	if errors.Is(err, ErrNoSubcommand) {
		listSubcommands()
		return nil
	} else if err != nil {
		return
	}

	cmd, err := reg.Open[Command](subcommand)
	if err != nil {
		return
	}

	err = cmd.Configure(app, subcommandArgs)
	if err != nil {
		return
	}

	err = cmd.Run(app)
	if err != nil {
		return
	}

	err = cmd.Print(app, os.Stdout)
	if err != nil {
		return
	}

	return
}
