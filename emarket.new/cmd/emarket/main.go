package main

import (
	"emarket/cmd/emarket/cli"
	db "emarket/internal/emarket/file"
	"emarket/internal/emarket/html/page"
	"fmt"
	"log"
)

func main() {
	opts := cli.Parse()
	magazService := db.NewMagazineService(opts.DataFile)
	magazines, err := magazService.Find()

	if err != nil {
		log.Fatalln(err)
	}

	for _, p := range page.MagazineList(magazines) {
		fmt.Println(string(p))
	}
}
