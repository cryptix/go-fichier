package main

import (
	"log"
	"os"

	"github.com/cryptix/go-fichier"
)

func main() {
	c, err := fichier.NewClient("user", os.Getenv("1FICHIER_PASS"))
	check(err)

	resp, err := c.GetInfo()
	check(err)

	log.Printf("GetInfo: %+v", resp)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
