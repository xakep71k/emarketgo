package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	host     = "tatyana"
	port     = 5432
	user     = "postgres"
	password = "1"
	dbname   = "emarketrordb"
)

var client *mongo.Client

type Product struct {
	ID          string `bson:"_id,omitempty"`
	Title       string `bson:"title"`
	Price       int    `bson:"price"`
	Gallery     []byte `bson:"thumb"`
	Enable      bool   `bson:"enable"`
	Description string `bson:"description"`
	Quantity    int    `bson:"quantity"`
	OldID       int    `bson:"oldid"`
	OldImgName  string `bson:"oldimgfile"`
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

	client, err := NewMongoClient()
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

		product := Product{
			Title:       *title,
			Price:       int(*price),
			Gallery:     blobs[0],
			Enable:      *enable,
			Description: *description,
			Quantity:    *maxCount,
			OldID:       *oldid,
			ID:          fmt.Sprint(uuid.New()),
			OldImgName:  gallery[0],
		}

		collection := client.Database("emarket").Collection("products")
		_, err = collection.InsertOne(context.Background(), product)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func NewMongoClient() (*mongo.Client, error) {
	ip := "127.0.0.1"
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + ip + ":27017"))
	if err != nil {
		fmt.Printf("new client: %v\n", err)
	}

	ctx := DefaultContext()
	err = client.Connect(ctx)
	if err != nil {
		fmt.Printf("connect: %v\n", err)
		return nil, err
	}

	/*
		defer func() {
			if err := client.Disconnect(ctx); err != nil {
				fmt.Printf("disconnect: %v\n", err)
			}
		}()
	*/

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Printf("ping: %v\n", err)
		return nil, err
	}

	return client, nil
}

func DefaultContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}

func Tx(callback func(ctx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	ctx := DefaultContext()

	s, err := client.StartSession()
	if err != nil {
		return nil, err
	}

	defer s.EndSession(ctx)

	opts := options.Transaction().SetReadConcern(readconcern.Snapshot())
	return s.WithTransaction(ctx, callback, opts)
}
