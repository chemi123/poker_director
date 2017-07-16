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

func (td *TournamentDirector) setTableAsRequested(tableReq table.TableRequest) {
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

	var tableReq table.TableRequest
	if err = json.Unmarshal(reqBody, &tableReq); err != nil {
		w.Write([]byte("Request Body is not a proper json format\n"))
		return
	}

	// TODO: tableIDを付ける規則は改めて修正する必要がある
	//       でないと例えばtableIDが1, 2, 3, 4とあって2がクローズされた後にまたtableが追加されたら既に存在しているtableID=4が再度出来上がる
	//       バグを仕込む可能性が高そうな箇所である
	if tableReq.ID == 0 {
		tableID := len(td.Tables) + 1
		td.Tables = append(td.Tables, table.NewTable(tableID, tableReq.PlayersNum))

		// TODO: ここでクライアント側にtableIDを返す処理が必要
		fmt.Fprintf(w, "Your tableID is %v\n", tableID)
	} else {
		td.setTableAsRequested(tableReq)
		tableBalance()
	}

	// debug
	for _, v := range td.Tables {
		fmt.Fprintf(w, "td: {table_id:%v, players_num:%v}\n", v.ID, v.PlayersNum)
	}
}
