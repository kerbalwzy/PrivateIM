package main

import (
	"log"
	"net/http"
)

func init() {
	log.SetPrefix("AuthCenter ")
}

func main() {
	http.HandleFunc("/login", Login)
	log.Fatal(http.ListenAndServe(":4444", nil))

}
