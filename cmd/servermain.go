package cmd

import (
	"github.com/urfave/cli"
)

var serverFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "address",
		Value: ":" + globalMinioPort,
		Usage: "Bind to a specific ADDRESS:PORT, ADDRESS can be an IP or hostname.",
	},
}

var serverCmd = cli.Command{
	Name:   "server",
	Usage:  "Start object storage server.",
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
ENVIRONMENT VARIABLES:
  ACCESS:
     MINIO_ACCESS_KEY: Custom username or access key of minimum 3 characters in length.
     MINIO_SECRET_KEY: Custom password or secret key of minimum 8 characters in length.

  ENDPOINTS:
     MINIO_ENDPOINTS: List of all endpoints delimited by ' '.

  BROWSER:
     MINIO_BROWSER: To disable web browser access, set this value to "off".

  CACHE:
     MINIO_CACHE_DRIVES: List of mounted drives or directories delimited by ";".
     MINIO_CACHE_EXCLUDE: List of cache exclusion patterns delimited by ";".
     MINIO_CACHE_EXPIRY: Cache expiry duration in days.
	
  DOMAIN:
     MINIO_DOMAIN: To enable virtual-host-style requests, set this value to Minio host domain name.

  WORM:
     MINIO_WORM: To turn on Write-Once-Read-Many in server, set this value to "on".

EXAMPLES:
  1. Start minio server on "/home/shared" directory.
     $ {{.HelpName}} /home/shared

  2. Start minio server bound to a specific ADDRESS:PORT.
     $ {{.HelpName}} --address 192.168.1.101:9000 /home/shared

  3. Start minio server and enable virtual-host-style requests.
     $ export MINIO_DOMAIN=mydomain.com
     $ {{.HelpName}} --address mydomain.com:9000 /mnt/export

  4. Start minio server on 64 disks server with endpoints through environment variable.
     $ export MINIO_ENDPOINTS=/mnt/export{1...64}
     $ {{.HelpName}}

  5. Start distributed minio server on an 8 node setup with 8 drives each. Run following command on all the 8 nodes.
     $ export MINIO_ACCESS_KEY=minio
     $ export MINIO_SECRET_KEY=miniostorage
     $ {{.HelpName}} http://node{1...8}.example.com/mnt/export/{1...8}

  6. Start minio server with edge caching enabled.
     $ export MINIO_CACHE_DRIVES="/mnt/drive1;/mnt/drive2;/mnt/drive3;/mnt/drive4"
     $ export MINIO_CACHE_EXCLUDE="bucket1/*;*.png"
     $ export MINIO_CACHE_EXPIRY=40
     $ {{.HelpName}} /home/shared
`,
}

func serverMain(ctx *cli.Context) {

}
