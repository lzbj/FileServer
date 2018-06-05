package cmd

import (
	"github.com/urfave/cli"
	"github.com/minio/minio/cmd/logger"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/lzbj/FileServer/lib/server"
	"github.com/lzbj/FileServer/lib/status"
	"github.com/lzbj/FileServer/lib/stats"
)

var serverFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "address",
		Value: ":" + globalPort,
		Usage: "Bind to a specific ADDRESS:PORT, ADDRESS can be an IP or hostname.",
	},
	cli.StringFlag{
		Name:  "fspath",
		Value: globalFSPath,
		Usage: "Bind the file storage server storage backend to a file system path",
	},
	cli.StringFlag{
		Name:  "backendtype",
		Value: storageDefaultBackEndType,
		Usage: "Determine the file storage server storage backend to a file system path or S3",
	},
	cli.StringFlag{
		Name:  "cachepath",
		Value: globalCacheFSPath,
		Usage: " ",
	},
}

var serverCmd = cli.Command{
	Name:   "server",
	Usage:  "Start file storage server.",
	Flags:  append(serverFlags, globalFlags...),
	Action: serverMain,
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} {{if .VisibleFlags}}[FLAGS] {{end}}DIR1 [DIR2..]
  {{.HelpName}} {{if .VisibleFlags}}[FLAGS] {{end}}DIR{1...64}

DIR:
  DIR points to a directory on a filesystem. When you want to combine
  multiple drives into a single large system, pass one directory per
  filesystem separated by space. You may also use a '...' convention
  to abbreviate the directory arguments. Remote directories in a
  distributed setup are encoded as HTTP(s) URIs.
{{if .VisibleFlags}}
FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}{{end}}
`,
}

func serverMain(ctx *cli.Context) {
	if ctx.Args().First() == "help" {
		cli.ShowCommandHelpAndExit(ctx, "server", 1)
	}

	logger.Disable = true

	// Get "json" flag from command line argument and
	// enable json and quite modes if jason flag is turned on.
	jsonFlag := ctx.IsSet("json") || ctx.GlobalIsSet("json")
	if jsonFlag {
		logger.EnableJSON()
	}

	// Get quiet flag from command line argument.
	quietFlag := ctx.IsSet("quiet") || ctx.GlobalIsSet("quiet")
	if quietFlag {
		logger.EnableQuiet()
	}

	handleCmdArgs(ctx)

	handleEnvArgs()

	initConfig()

	// configure server.
	var err error
	var handler http.Handler
	handler, err = configureServerHandler()
}

func handleCmdArgs(ctx *cli.Context) {

}

func handleEnvArgs() {

}

func initConfig() {

}

func configureServerHandler() (http.Handler, error) {
	router := mux.NewRouter().SkipClean(true)

	//Register upload api router

	server.RegisterStorageServerRouter(router)

	//Register status router
	status.RegisteStatusRouter(router)

	//Register statistics api router
	stats.RegisteStatusRouter(router)

	return nil, nil
}
