package main

import (
	"./director"
	"net/http"
)

func main() {
	handler := &director.TournamentDirector{}
	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}
