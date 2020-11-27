package cli

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

type Opts struct {
	WEBRoot string
	Listen  string
	DBFile  string
}

func Parse() Opts {
	if len(os.Args) != 6 {
		log.Fatalf("Usage: %s --web-root <path> --listen <ip:port> --data <path>\n", os.Args[0])
		os.Exit(1)
	}

	webRootOpt := flag.String("web-root", "", "<path>")
	listenOpt := flag.String("listen", "", "<ip:port>")
	dataOpt := flag.String("data", "", "<path>")
	flag.Parse()

	if webRootOpt == nil || *webRootOpt == "" {
		log.Fatalln("web root not specified")
	}

	if listenOpt == nil || *listenOpt == "" {
		log.Fatalln("listen ip:port not specified")
	}

	if dataOpt == nil || *dataOpt == "" {
		log.Fatalln("listen ip:port not specified")
	}

	abs := func(file string) string {
		fullpath, err := filepath.Abs(file)

		if err != nil {
			log.Fatalln(err)
		}

		return fullpath
	}

	return Opts{
		DBFile:  abs(*dataOpt),
		WEBRoot: abs(*webRootOpt),
		Listen:  *listenOpt,
	}
}
