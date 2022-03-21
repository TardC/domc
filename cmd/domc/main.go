package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tardc/domc"
)

func main() {
	f1 := flag.String("h1", "", "one HTML filename")
	f2 := flag.String("h2", "", "another HTML filename")
	flag.Parse()

	html1, err := os.ReadFile(*f1)
	if err != nil {
		log.Fatalln(err)
	}
	html2, err := os.ReadFile(*f2)
	if err != nil {
		log.Fatalln(err)
	}

	nodes1, err := domc.ParseHTML(string(html1))
	if err != nil {
		log.Fatalln(err)
	}
	nodes2, err := domc.ParseHTML(string(html2))
	if err != nil {
		log.Fatalln(err)
	}

	v1 := domc.NewVector(nodes1)
	v2 := domc.NewVector(nodes2)

	fmt.Println(v1.IsSimilar(v2))
}
