package main

import (
	"flag"
	"fmt"

	galleryapi "github.com/Aeroxee/gallery-api"
)

func main() {
	hostname := flag.String("hostname", "127.0.0.1", "Type your hostname")
	port := flag.String("port", "8000", "Type your port number")

	flag.Parse()
	addr := fmt.Sprintf("%v:%v", *hostname, *port)

	galleryapi.Run(addr)
}
