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
	Tables []table.Table
}

type TableRequest struct {
	ID         int `json:"id"`
	PlayersNum int `json:"playersnum"`
}

func (td *TournamentDirector) setTableAsRequested(tableReq TableRequest) {
	for i, _ := range td.Tables {
		if td.Tables[i].ID == tableReq.ID {
			td.Tables[i].PlayersNum = tableReq.PlayersNum
		}
	}
}

func tableBalance() {
	// table balanceするよ
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

	// TODO: tableIdを付ける規則は改めて修正する必要がある
	//       でないと例えばtableIdが1, 2, 3, 4とあって2がクローズされた後にまたtableが追加されたら既に存在しているtableId=4が再度出来上がる
	//       バグを仕込む可能性が高そうな箇所である
	if tableReq.ID == 0 {
		tableId := len(td.Tables) + 1
		td.Tables = append(td.Tables, table.NewTable(tableId, tableReq.PlayersNum))

		// TODO: ここでクライアント側にtableIdを返す処理が必要
		fmt.Fprintf(w, "Your tableId is %v\n", tableId)
	} else {
		td.setTableAsRequested(tableReq)
		tableBalance()
	}

	// debug
	for _, v := range td.Tables {
		fmt.Fprintf(w, "td: {table_id:%v, players_num:%v}\n", v.ID, v.PlayersNum)
	}
}
