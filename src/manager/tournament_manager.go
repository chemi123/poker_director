package manager

import (
	// 相対パスよくないけど暫定で使う
	"../request"
	"../table"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TournamentManager struct {
	Tables               []table.Table
	TournamentDirectorID int
}

// リクエストで指定されたテーブルに値をbodyのjson通りにセットする
// もし指定のTable IDが存在しなければエラーを返す
// この関数が正しく動作するためにはTable IDは全てユニークでなければならない
func (tm *TournamentManager) setTableAsRequested(tableReq request.DealerRequest) error {
	for i, _ := range tm.Tables {
		if tm.Tables[i].ID == tableReq.ID {
			tm.Tables[i].PlayersNum = tableReq.PlayersNum
			return nil
		}
	}
	return errors.New("Requested Table ID does not exist")
}

// tableBalance
func (tm *TournamentManager) tableBalance() {
	minTable, maxTable := &tm.Tables[0], &tm.Tables[0]

	// TODO: 大した計算量ないからひとまずは愚直に計算する
	//       もっと効率化はできるがやるならバグに気をつけないといけないから費用対効果としては微妙かも。十分速いし
	for {
		for i, _ := range tm.Tables {
			if minTable.PlayersNum > tm.Tables[i].PlayersNum {
				minTable = &tm.Tables[i]
			}

			if maxTable.PlayersNum < tm.Tables[i].PlayersNum {
				maxTable = &tm.Tables[i]
			}
		}

		if (maxTable.PlayersNum - minTable.PlayersNum) < 2 {
			break
		}

		minTable.PlayersNum += 1
		maxTable.PlayersNum -= 1
	}
}

func (tm *TournamentManager) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.Write([]byte("Only takes POST request"))
		return
	}

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.Write([]byte("Failed to read request body\n"))
		return
	}

	var tableReq request.DealerRequest
	if err = json.Unmarshal(reqBody, &tableReq); err != nil {
		w.Write([]byte("Request Body is not a proper json format\n"))
		return
	}

	// TODO: tableIDを付ける規則は改めて修正する必要がある
	//       でないと例えばtableIDが1, 2, 3, 4とあって2がクローズされた後にまたtableが追加されたら既に存在しているtableID=4が再度出来上がる
	//       バグを仕込む可能性が高そうな箇所である
	if tableReq.ID == 0 {
		tableID := len(tm.Tables) + 1
		tm.Tables = append(tm.Tables, table.NewTable(tableID, tableReq.PlayersNum))

		// TODO: ここでクライアント側にtableIDを返す処理が必要
		fmt.Fprintf(w, "Your tableID is %v\n", tableID)
	} else if len(tm.Tables) > 0 {
		if err = tm.setTableAsRequested(tableReq); err != nil {
			fmt.Fprintln(w, "Requested ID does not exist")
			return
		}
		tm.tableBalance()
	} else {
		w.Write([]byte("No table is set yet\n"))
		return
	}

	// debug
	for _, v := range tm.Tables {
		fmt.Fprintf(w, "tm: {table_id:%v, players_num:%v}\n", v.ID, v.PlayersNum)
	}
}
