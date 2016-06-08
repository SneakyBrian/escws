package main

//go:generate esc -o static.go static

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	useLocalPtr := flag.Bool("useLocal", false, "use local filesystem for debugging purposes")
	portNumPtr := flag.Uint("port", 8181, "port number for web server")
	flag.Parse()
	http.Handle("/static/", http.FileServer(FS(*useLocalPtr)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *portNumPtr), nil))
}
