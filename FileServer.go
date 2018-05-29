package main

import (
	"github.com/lzbj/FileServer/cmd"
	"os"
)

var (
	REVISION     = ""
	LIB_REVISION = ""
	GO_VERSION   = ""
)

func main() {
	cmd.Main(os.Args)
}
