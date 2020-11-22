package main

import (
	"emarket/cmd/emarket/utils"
	"emarket/internal/emarket/file"
	"emarket/internal/emarket/html/page"
	"fmt"
	"log"
)

func main() {
	cmdArgs := utils.ParseCmdLine()
	magazService := file.NewMagazineService(cmdArgs.DataFile)
	magazines, err := magazService.Find()

	if err != nil {
		log.Fatalln(err)
	}

	for _, p := range page.MagazineList(magazines) {
		fmt.Println(string(p))
	}
}
