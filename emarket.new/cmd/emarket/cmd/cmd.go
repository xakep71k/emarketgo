package cmd

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

type CmdArgs struct {
	WEBRoot  string
	Listen   string
	DataFile string
}

func Parse() CmdArgs {
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

	dataFile, err := filepath.Abs(*dataOpt)

	if err != nil {
		log.Fatalln(err)
	}

	webRoot, err := filepath.Abs(*webRootOpt)

	if err != nil {
		log.Fatalln(err)
	}

	return CmdArgs{
		DataFile: dataFile,
		Listen:   *listenOpt,
		WEBRoot:  webRoot,
	}
}
