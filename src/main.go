package main

import (
	// 相対パスよくないけど暫定で使う
	"github.com/chemi123/poker_director/src/manager"
	"net/http"
)

func main() {
	handler := &manager.TournamentManager{}
	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}
