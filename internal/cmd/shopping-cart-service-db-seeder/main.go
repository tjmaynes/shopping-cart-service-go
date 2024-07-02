package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"

	driver "github.com/tjmaynes/shopping-cart-service-go/internal/driver"
	"github.com/tjmaynes/shopping-cart-service-go/internal/pkg/item"
)

// SeedData ..
func SeedData(jsonSource string, dbConn *sql.DB) []uuid.UUID {
	cartRepository := item.NewRepository(dbConn)
	ctx := context.Background()
	cartService := item.NewService(cartRepository)

	jsonFile, err := os.Open(jsonSource)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer jsonFile.Close()

	jsonBytes, _ := io.ReadAll(jsonFile)

	var items []item.ItemDTO
	err = json.Unmarshal(jsonBytes, &items)
	if err != nil {
		panic(err)
	}

	var ids []uuid.UUID
	for _, rawItem := range items {
		item, err := cartService.AddItem(ctx, &rawItem)
		if err != nil {
			panic(err)
		}
		ids = append(ids, item.ID)
	}

	return ids
}

func main() {
	var (
		dbSource       = flag.String("db-source", "./db/my.db", "Database url connection string.")
		seedDataSource = flag.String("seed-data-source", "./db/seed.json.", "JSON Source, such as ./db/seed.json.")
	)

	flag.Parse()

	dbConn, err := driver.ConnectDB(*dbSource)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	ids := SeedData(*seedDataSource, dbConn)
	fmt.Printf("Added %d entries!", len(ids))
}
