package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/conejoninja/bundle"
	"strings"
)

func main() {

	b, e := bundle.LoadBundle("./web.bundle")

	if e != nil {
		log.Fatalf("Error loading bundle: %s", e)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path[1:] == "gopherize.me.png" {
			gopherImg, err := b.Asset("gopherize.me.png")
			if err != nil {
				fmt.Fprintf(w, "Asset gopherize.me.png does not exists: %s!", err)
				return
			}
			fmt.Fprint(w, string(gopherImg))
			return
		}

		htmlByte, err := b.Asset("index.html")
		if err != nil {
			fmt.Fprintf(w, "Asset index.html does not exists: %s!", err)
			return
		}
		fmt.Fprint(w, strings.Replace(string(htmlByte), "{SERVER}", "http://" + ln.Addr().String(), -1))
	})

	fmt.Println("Open your browser at http://" + ln.Addr().String())
	log.Fatal(http.Serve(ln, nil))

}

