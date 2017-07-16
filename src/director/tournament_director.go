package director

import (
	// 相対パスよくないけど暫定で使う
	"../table"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TournamentDirector struct {
	Tables []table.Table
}

// リクエストで指定されたテーブルに値をbodyのjson通りにセットする
// もし指定のTable IDが存在しなければエラーを返す
func (td *TournamentDirector) setTableAsRequested(tableReq table.TableRequest) error {
	for i, _ := range td.Tables {
		if td.Tables[i].ID == tableReq.ID {
			td.Tables[i].PlayersNum = tableReq.PlayersNum
			return nil
		}
	}
	return errors.New("Requested Table ID does not exist")
}

// tableBalance
func (td *TournamentDirector) tableBalance() {
	minTable, maxTable := &td.Tables[0], &td.Tables[0]

	// TODO: 大した計算量ないからまずは愚直に計算するがこのやり方はカッコ悪いので後で効率化を図る修正をする
	for {
		for i, _ := range td.Tables {
			if minTable.PlayersNum > td.Tables[i].PlayersNum {
				minTable = &td.Tables[i]
			}

			if maxTable.PlayersNum < td.Tables[i].PlayersNum {
				maxTable = &td.Tables[i]
			}
		}

		if (maxTable.PlayersNum - minTable.PlayersNum) < 2 {
			break
		}

		minTable.PlayersNum += 1
		maxTable.PlayersNum -= 1
	}
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
	} else if len(td.Tables) > 0 {
		if err = td.setTableAsRequested(tableReq); err != nil {
			fmt.Fprintln(w, "Requested ID does not exist")
			return
		}
		td.tableBalance()
	} else {
		w.Write([]byte("No table is set yet\n"))
		return
	}

	// debug
	for _, v := range td.Tables {
		fmt.Fprintf(w, "td: {table_id:%v, players_num:%v}\n", v.ID, v.PlayersNum)
	}
}
