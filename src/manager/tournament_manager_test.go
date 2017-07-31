package manager

import (
	//"fmt"
	"github.com/chemi123/poker_director/src/table"
	"net/http"
	"testing"
)

// TODO: 全部はテストできていないうまいやり方を見つける必要あり
func TestParseJSONRequest(t *testing.T) {
	req := http.Request{Method: "GET"}
	_, err := parseJSONRequest(&req)
	if err == nil {
		t.Fatal("parseJSONRequest only takes POST Method. But it takes other Methods.")
	}
}

func TestSetTableAsRequested(t *testing.T) {
	tm := TournamentManager{tables: []table.Table{table.Table{ID: 1, PlayersNum: 5}}}

	// PlayersNumが5->4になっていることを確認
	err := tm.setTableAsRequested(1, 4)
	if err != nil {
		t.Fatal(err)
	}

	if tm.tables[0].PlayersNum != 4 {
		t.Fatalf("PlayersNum should be 4, but it is %v", tm.tables[0].PlayersNum)
	}

	// 存在しないテーブル(この場合ID=2)をセットしようとした場合は異常ケースとしてerrorが返ってくることの確認
	err = tm.setTableAsRequested(2, 1)
	if err == nil {
		t.Fatal("Should output error when set to non existing table")
	}
}

func TestBalanceTable(t *testing.T) {
	tm := TournamentManager{tables: []table.Table{table.Table{ID: 1, PlayersNum: 5},
		table.Table{ID: 2, PlayersNum: 6},
		table.Table{ID: 3, PlayersNum: 8}}}

	tm.balanceTable()

	// TODO: もっとエレガントに書く
	if tm.tables[0].PlayersNum != 6 {
		t.Fatal("table ID: 1 should have 6 players as a result of balanceTable. It has %s players", tm.tables[0].PlayersNum)
	}

	if tm.tables[1].PlayersNum != 6 {
		t.Fatal("table ID: 2 should have 6 players as a result of balanceTable. It has %s players", tm.tables[1].PlayersNum)
	}

	if tm.tables[2].PlayersNum != 7 {
		t.Fatal("table ID: 3 should have 7 players as a result of balanceTable. It has %s players", tm.tables[2].PlayersNum)
	}
}

// TODO
func TestHandleDealerRequest(t *testing.T) {
}
