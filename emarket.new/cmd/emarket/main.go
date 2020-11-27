package main

import (
	"emarket/cmd/emarket/cli"
	db "emarket/internal/emarket/db/fs"
	"emarket/internal/emarket/http"
	"fmt"
	"log"
	"os"
)

func main() {
	opts := cli.Parse()
	magazStorage := db.NewMagazineStorage(opts.DBFile)
	handler := http.NewEMarketHandler(opts.WEBRoot, magazStorage)
	srv := http.NewServer(handler, opts.Listen)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("started %v\n", dir)
	log.Fatal(srv.ListenAndServe())
}
