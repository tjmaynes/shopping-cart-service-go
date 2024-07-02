package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"

	"github.com/icrowley/fake"
	"github.com/tjmaynes/shopping-cart-service-go/internal/pkg/item"
)

// GenerateSeedData ..
func GenerateSeedData(jsonDestination string, data []item.Item) {
	json, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	err = ioutil.WriteFile(jsonDestination, json, 0644)
	if err != nil {
		panic(err)
	}
}

func generateCartItems(itemCount int, manufacturerCount int) []item.Item {
	var manufacturers []string
	for i := 0; i < manufacturerCount; i++ {
		manufacturers = append(manufacturers, fake.Brand())
	}

	var items []item.Item
	for i := 0; i < itemCount; i++ {
		items = append(items, item.Item{
			Name:         fake.ProductName(),
			Price:        item.Decimal(int64(rand.Intn(100) + 100)),
			Manufacturer: manufacturers[rand.Intn(4)],
		})
	}

	return items
}

func main() {
	var (
		seedDataDestination = flag.String("seed-data-destination", "./db/seed.json", "JSON Destination, such as ./db/seed.json.")
		itemCount           = flag.Int("item-count", 50, "Number of seed items, such as 10, 20, 50, etc.")
		manufacturerCount   = flag.Int("manufacturer-count", 5, "Number of unique manufacturers, such as 10, 20, 50, etc.")
	)
	flag.Parse()

	GenerateSeedData(*seedDataDestination, generateCartItems(*itemCount, *manufacturerCount))
}
