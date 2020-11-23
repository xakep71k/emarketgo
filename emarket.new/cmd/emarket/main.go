package main

import (
	"emarket/cmd/emarket/cli"
	db "emarket/internal/emarket/file"
	"emarket/internal/emarket/http/html/page"
	"fmt"
	"log"
)

func main() {
	opts := cli.Parse()
	magazStorage := db.NewMagazineStorage(opts.DataFile)
	magazines, err := magazStorage.Find()

	if err != nil {
		log.Fatalln(err)
	}

	for _, p := range page.MagazineList(magazines) {
		fmt.Println(string(p))
	}
}
