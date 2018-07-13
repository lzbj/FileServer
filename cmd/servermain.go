package cmd

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lzbj/FileServer/lib/server"
	"github.com/lzbj/FileServer/lib/stats"
	"github.com/lzbj/FileServer/lib/status"
	"github.com/lzbj/FileServer/lib/storage"
	"github.com/minio/minio/cmd/logger"
	"github.com/urfave/cli"
	"net/http"
)

var serverFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "address",
		Value: server.GlobalPort,
		Usage: "Bind to a specific ADDRESS:PORT, ADDRESS can be an IP or hostname.",
	},
	cli.StringFlag{
		Name:  "fspath",
		Value: server.GlobalFSPath,
		Usage: "Bind the file storage server storage backend to a file system path",
	},
	cli.StringFlag{
		Name:  "backendtype",
		Value: server.StorageDefaultBackEndType,
		Usage: "Determine the file storage server storage backend to a file system path or S3",
	},
	cli.StringFlag{
		Name:  "cachepath",
		Value: server.GlobalCacheFSPath,
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

	logger.Disable = false

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
	var err error

	// Do some checks, such as the address.
	handleCmdArgs(ctx)

	// Do some env checks.
	handleEnvArgs()

	// Do some initialization works.
	err = initConfig(ctx)
	if err != nil {

	}

	initStorage(ctx)

	initDBSys(ctx)

	initDBCache(ctx)

	initMonitorSys(ctx)

	// configure server.

	var handler http.Handler
	handler, err = configureServerHandler()
	if err != nil {
		logger.Fatal(err, "Unable to configure one of server's RPC services")
	}

	httpServer := server.NewServer(ctx.String("address"), handler)
	go func() {
		server.GlobalHTTPServerErrorCh <- httpServer.Start()
	}()

	if err != nil {
		httpServer.Shutdown()
	}
	handleSignals()
}

// handle command args and do some checks.
func handleCmdArgs(ctx *cli.Context) {
	fmt.Println(ctx.String("address"))

}

// handle ENV args if necessary
func handleEnvArgs() {

}

// Do some configuration initialization and conflicts and schema checks.
func initConfig(ctx *cli.Context) error {
	fpath := ctx.String("fspath")
	if len(fpath) != 0 {
		server.GlobalFSPath = fpath
	}
	return nil
}

func initDBCache(ctx *cli.Context) error {
	//Not implemented yet.
	return nil
}

func initDBSys(ctx *cli.Context) error {
	//Not implemented yet.
	return nil
}

func initMonitorSys(ctx *cli.Context) error {
	//Not implemented yet.
	return nil
}

func initStorage(ctx *cli.Context) error {
	var err error
	if server.GlobalFSInitalized {
		logger.Info("global file storage already initialized")
		return errors.New("already initialized")
	}
	logger.Info("start to initialize the global FS storage...")
	fmt.Println(server.GlobalFSPath)
	server.GlobalBackEndFSSys, err = storage.NewFStorage(server.GlobalFSPath)
	if err != nil {
		logger.Info("error %s", err)
		return errors.New(err.Error())
	}
	return nil
}

func configureServerHandler() (http.Handler, error) {
	logger.Info("Start to configure server handlers...")
	router := mux.NewRouter().SkipClean(true)

	//Register upload api router
	server.RegisterStorageServerRouter(router)
	server.RegisterStorageServerRouterDownload(router)

	//Register status router
	status.RegisteStatusRouter(router)

	//Register statistics api router
	stats.RegisteStatusRouter(router)
	logger.Info("Configure server handlers finished .")
	return server.RegisterHandlers(router, server.GlobalHandlers...), nil

}

func handleSignals() {
	for {
		select {
		case signal := <-server.GlobalHTTPServerErrorCh:
			switch signal {

			}

		case osSignal := <-server.GlobalOSSignalCh:

			switch osSignal {

			}
		}
	}
}
