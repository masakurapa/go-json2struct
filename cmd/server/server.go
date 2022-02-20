package main

import (
	"fmt"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("../../web"))
	http.Handle("/", fs)
	fmt.Println("server is listening on port 8000")
	http.ListenAndServe(":8000", nil)
}
