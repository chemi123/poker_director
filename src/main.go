package main

import (
	// 相対パスよくないけど暫定で使う
	"./director"
	"net/http"
)

func main() {
	handler := &director.TournamentDirector{}
	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}
