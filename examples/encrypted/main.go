package main

import (
	"github.com/conejoninja/bundle"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {

	b, e := bundle.LoadBundle("./filename.bundle", []byte("12345678901234"))

	if e != nil {
		log.Fatalf("Error loading bundle: %s", e)
	}
	for k := range b.Assets {
		data, err := b.Asset(k)
		if err == nil {
			err = ioutil.WriteFile(k, data, 0644)
			if err != nil {
				fmt.Println("Unable to create File", err)
				return
			}
		}
	}
}
