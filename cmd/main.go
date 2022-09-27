package main

import (
	"airtonomy-assessment/server"
	"fmt"
	"net/http"
)

const PORT = "8080"

func main() {
	s := server.NewServer()
	fmt.Printf("Running server on port %s\n", PORT)
	http.ListenAndServe(":"+PORT, s)
}
