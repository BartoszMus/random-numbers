package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	os.Setenv("PORT", "8000")
	http.HandleFunc("/random/mean", GetResponses)
	log.Printf("You can see your responses at port %s", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
