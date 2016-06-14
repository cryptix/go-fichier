package main

import (
	"log"
	"os"

	"github.com/cryptix/go-fichier"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <fname>", os.Args[0])
	}
	fname := os.Args[1]

	f, err := os.Open(fname)
	check(err)

	lnkDl, lnkRm, err := fichier.UploadFile(fname, f)
	check(err)

	log.Println("DL link:", lnkDl)
	log.Println("remove link:", lnkRm)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
