package main

import (
	"emarket/internal/magazine"
	"emarket/internal/mux"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	ListenAddr string
	WebRoot    string
	DataPath   string
}

func main() {
	cfg, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	magazineC, err := prepareMagazineController(cfg)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Handler:      mux.New(magazineC),
		Addr:         cfg.ListenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("started at %v\n", dir)
	log.Fatal(srv.ListenAndServe())
}

func readConfig() (cfg Config, err error) {
	if len(os.Args) != 6 {
		fmt.Printf("Usage: %s --web-root <path> --listen <ip:port> --data <path>\n", os.Args[0])
		os.Exit(1)
	}

	webRootOpt := flag.String("web-root", "", "<path>")
	listenOpt := flag.String("listen", "", "<ip:port>")
	dataOpt := flag.String("data", "", "<path>")
	flag.Parse()

	if webRootOpt == nil || *webRootOpt == "" {
		err = errors.New("web root not specified")
		return
	}

	if listenOpt == nil || *listenOpt == "" {
		err = errors.New("listen ip:port not specified")
		return
	}

	if dataOpt == nil || *dataOpt == "" {
		err = errors.New("listen ip:port not specified")
		return
	}

	var dataPath, webRoot string
	dataPath, err = filepath.Abs(*dataOpt)

	if err != nil {
		return
	}

	webRoot, err = filepath.Abs(*webRootOpt)

	return Config{
		WebRoot:    webRoot,
		DataPath:   dataPath,
		ListenAddr: *listenOpt,
	}, err
}

func prepareMagazineController(cfg Config) (c *magazine.Controller, err error) {
	var m *magazine.Model
	var v *magazine.View

	m, err = magazine.NewModel(cfg.DataPath)
	if err != nil {
		return
	}

	v, err = magazine.NewView(cfg.WebRoot)
	if err != nil {
		return
	}

	err = v.PrepareStaticContent(m.GetSortedRecords())
	if err != nil {
		return
	}

	c, err = magazine.NewController(m, v)
	if err != nil {
		return
	}

	return
}
