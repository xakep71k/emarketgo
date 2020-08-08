package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	host     = "tatyana"
	port     = 5432
	user     = "postgres"
	password = "1"
	dbname   = "emarketrordb"
)

type Product struct {
	ID          string `bson:"_id,omitempty" json:"id"`
	Title       string `bson:"title" json:"title"`
	Price       int    `bson:"price" json:"price"`
	Thumb       []byte `bson:"thumb" json:"thumb"`
	Enable      bool   `bson:"enable" json:"enable"`
	Description string `bson:"description" json:"description"`
	Quantity    int    `bson:"quantity" json:"quantity"`
	OldID       int    `bson:"oldid" json:"oldid"`
	OldImgName  string `bson:"oldimgfile" json:"oldimgfile"`
	PageNum     int    `bson:"-" json:"-"`
}

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

	rows, err := db.Query(`SELECT id,title,gallery,price,max_count,enable,description FROM products`)

	if err != nil {
		log.Fatal(err)
	}

	var products []*Product

	for rows.Next() {
		var title *string
		var galleryData []byte
		var price *float64
		var maxCount *int
		var enable *bool
		var description *string
		var oldid *int

		err = rows.Scan(&oldid, &title, &galleryData, &price, &maxCount, &enable, &description)
		if err != nil {
			log.Fatal(err)
		}

		var gallery []string
		err = json.Unmarshal(galleryData, &gallery)
		if err != nil {
			log.Fatal(err)
		}

		if oldid == nil {
			log.Fatal("id cannot be nil")
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

		fmt.Printf("%v %v %v %v %v %v\n", *oldid, *title, gallery, *price, *maxCount, *enable)
		blobs := make([][]byte, 0)
		for _, file := range gallery {
			image := "/home/alek/emarket_files/gallery/" + file
			info, err := os.Stat(image)
			if err != nil {
				log.Fatalln(err)
			}

			if info.Size() >= 1024*1024 {
				log.Fatalf("size too big %v\n", image)
			}

			blob, err := ioutil.ReadFile(image)
			blobs = append(blobs, blob)
			if err != nil {
				log.Fatal(err)
			}
		}

		if len(blobs) != 1 {
			log.Fatal("more then 1 image")
		}

		product := &Product{
			Title:       *title,
			Price:       int(*price),
			Thumb:       blobs[0],
			Enable:      *enable,
			Description: *description,
			Quantity:    *maxCount,
			OldID:       *oldid,
			ID:          fmt.Sprint(uuid.New()),
			OldImgName:  gallery[0],
		}

		products = append(products, product)
	}

	data, err := json.MarshalIndent(products, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
}
