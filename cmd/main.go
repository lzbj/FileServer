package cmd

import (
	"github.com/minio/mc/pkg/console"
	"github.com/minio/minio/pkg/trie"
	"github.com/minio/minio/pkg/words"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
	"sort"
)

// global flags for file server.
var globalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "config-dir, C",
		Value: getConfigDir(),
		Usage: func() string {
			usage := "Path to configuration directory."
			if getConfigDir() == "" {
				usage = usage + "  This option must be set."
			}
			return usage
		}(),
	},
	cli.BoolFlag{
		Name:  "quiet",
		Usage: "Disable startup information.",
	},
	cli.BoolFlag{
		Name:  "json",
		Usage: "Output server logs and startup information in json format.",
	},
}

var helpTemplate = `NAME:
  {{.Name}} - {{.Usage}}

DESCRIPTION:
  {{.Description}}

USAGE:
  {{.HelpName}} {{if .VisibleFlags}}[FLAGS] {{end}}COMMAND{{if .VisibleFlags}}{{end}} [ARGS...]

COMMANDS:
  {{range .VisibleCommands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
  {{end}}{{if .VisibleFlags}}
FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}{{end}}
VERSION:
  ` +
	`{{ "\n"}}`

// Main main for file server.
func Main(args []string) {
	appName := filepath.Base(args[0])

	// Run the app - exit on error.
	if err := newApp(appName).Run(args); err != nil {
		os.Exit(1)
	}
}

func newApp(name string) *cli.App {
	commands := []cli.Command{}

	commandsTree := trie.NewTrie()

	// registerCommand registers a cli command.
	registerCommand := func(command cli.Command) {
		commands = append(commands, command)
		commandsTree.Insert(command.Name)
	}

	findClosestCommands := func(command string) []string {
		var closestCommands []string
		for _, value := range commandsTree.PrefixMatch(command) {
			closestCommands = append(closestCommands, value.(string))
		}

		sort.Strings(closestCommands)
		// Suggest other close commands - allow missed, wrongly added and
		// even transposed characters
		for _, value := range commandsTree.Walk(commandsTree.Root()) {
			if sort.SearchStrings(closestCommands, value.(string)) < len(closestCommands) {
				continue
			}
			// 2 is arbitrary and represents the max
			// allowed number of typed errors
			if words.DamerauLevenshteinDistance(command, value.(string)) < 2 {
				closestCommands = append(closestCommands, value.(string))
			}
		}

		return closestCommands
	}

	// Register all commands.
	registerCommand(serverCmd)

	// Set up app.
	cli.HelpFlag = cli.BoolFlag{
		Name:  "help, h",
		Usage: "Show help.",
	}

	app := cli.NewApp()
	app.Name = name
	app.Author = "File Server Author"
	app.Usage = "File Upload Storage Server."
	app.Description = `File Upload Storage Server`
	app.Flags = globalFlags
	app.HideVersion = true
	app.Commands = commands
	app.CustomAppHelpTemplate = helpTemplate
	app.CommandNotFound = func(ctx *cli.Context, command string) {
		console.Printf("‘%s’ is not a supported sub-command. See ‘fileserver --help’.\n", command)
		closestCommands := findClosestCommands(command)
		if len(closestCommands) > 0 {
			console.Println()
			console.Println("Did you mean one of these?")
			for _, cmd := range closestCommands {
				console.Printf("\t‘%s’\n", cmd)
			}
		}

		os.Exit(1)
	}

	return app
}
