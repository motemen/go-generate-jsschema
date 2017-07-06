package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	gen "github.com/motemen/go-generate-jsschema"
)

func main() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)

	g := gen.Generator{}
	err := g.FromArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(g.Schema)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
