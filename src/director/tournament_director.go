package director

import (
	// 相対パスよくないけど暫定で使う
	"../table"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TournamentDirector struct {
	tables []table.Table
}

type TableRequest struct {
	ID         int  `json:"id"`
	PlayersNum int  `json:"playersnum"`
	NewTable   bool `json:"newtable"`
}

func (td *TournamentDirector) setTableAsRequested(tableReq TableRequest) {
}

func (td *TournamentDirector) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.Write([]byte("Only takes POST request"))
		return
	}

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.Write([]byte("Failed to read request body\n"))
		return
	}

	var tableReq TableRequest
	if err = json.Unmarshal(reqBody, &tableReq); err != nil {
		w.Write([]byte("Request Body is not a proper json format\n"))
		return
	}

	if tableReq.NewTable == true {
		td.tables = append(td.tables, table.NewTable(tableReq.ID, tableReq.PlayersNum))
	} else {
		td.setTableAsRequested(tableReq)
	}

	// debug
	for _, v := range td.tables {
		fmt.Fprintf(w, "td: {table_id:%v, players_num:%v}\n", v.ID, v.PlayersNum)
	}
}
