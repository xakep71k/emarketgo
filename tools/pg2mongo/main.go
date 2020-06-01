package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	host     = "tatyana"
	port     = 5432
	user     = "postgres"
	password = "1"
	dbname   = "emarketrordb"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("open: %v\n", err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("ping: %v\n", err)
	}

	rows, err := db.Query(`SELECT title,gallery,price,max_count,enable,description FROM products`)

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var title *string
		var galleryData []byte
		var price *float64
		var maxCount *int
		var enable *bool
		var description *string

		err = rows.Scan(&title, &galleryData, &price, &maxCount, &enable, &description)
		if err != nil {
			log.Fatal(err)
		}

		var gallery []string
		err = json.Unmarshal(galleryData, &gallery)
		if err != nil {
			log.Fatal(err)
		}

		initString := func() *string {
			s := ""
			return &s
		}

		if title == nil {
			title = initString()
		}

		if price == nil {
			f := 0.0
			price = &f
		}

		if maxCount == nil {
			i := 0
			maxCount = &i
		}

		if enable == nil {
			b := false
			enable = &b
		}

		if description == nil {
			description = initString()
		}

		if *enable {
			fmt.Printf("%v %v %v %v\n", *title, gallery, *price, *maxCount)
			for _, file := range gallery {
				image := "/home/alek/emarket_files/gallery/" + file
				info, err := os.Stat(image)
				if err != nil {
					log.Fatalln(err)
				}

				if info.Size() >= 1024*1024 {
					log.Fatalf("size too big %v\n", image)
				}
			}
		}
	}
}
