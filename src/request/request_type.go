package request

// Dealerからのリクエスト(Tableに関するもの)
type DealerRequest struct {
	ID         int `json:"id"`
	PlayersNum int `json:"playersnum"`
}

// Tournament Directorからのリクエスト(テーブルバランス等の指示)
type TournamentDirectorRequest struct {
	ID int `json:"id"`
}
