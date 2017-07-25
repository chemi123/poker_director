package main

import (
	"github.com/chemi123/poker_director/src/manager"
	"net/http"
)

func main() {
	handler := &manager.TournamentManager{}
	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}
