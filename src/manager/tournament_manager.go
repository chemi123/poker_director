package manager

import (
	// 相対パスよくないけど暫定で使う
	"../table"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
)

type TournamentManager struct {
	Tables               []table.Table
	TournamentDirectorID int
	RequestJson          *simplejson.Json
}

// リクエストで指定されたテーブルに値をbodyのjson通りにセットする
// もし指定のTable IDが存在しなければエラーを返す
// この関数が正しく動作するためにはTable IDは全てユニークでなければならない
func (tm *TournamentManager) setTableAsRequested(requestedTableId int, requestedPlayersNum int) error {
	for i, _ := range tm.Tables {
		if tm.Tables[i].ID == requestedTableId {
			tm.Tables[i].PlayersNum = requestedPlayersNum
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

func parseJsonRequest(httpReq *http.Request) (*simplejson.Json, error) {
	if httpReq.Method != http.MethodPost {
		return nil, errors.New("Only takes POST request\n")
	}

	reqBody, err := ioutil.ReadAll(httpReq.Body)
	if err != nil {
		return nil, errors.New("Failed to read request body\n")
	}

	return simplejson.NewJson(reqBody)
}

func (tm *TournamentManager) handleDealerRequest(w http.ResponseWriter) {
	tableId, err := tm.RequestJson.Get("ID").Int()
	if err != nil {
		fmt.Println("Failed to assert key \"ID\"")
		return
	}

	playersNum, err := tm.RequestJson.Get("PlayersNum").Int()
	if err != nil {
		fmt.Println("Failed to assert key \"PlayersNum\"")
		return
	}

	// TODO: tableIdを付ける規則は改めて修正する必要がある
	//       でないと例えばtableIdが1, 2, 3, 4とあって2がクローズされた後にまたtableが追加されたら既に存在しているtableId=4が再度出来上がる
	//       バグを仕込む可能性が高そうな箇所である
	if tableId == 0 {
		tableId = len(tm.Tables) + 1
		tm.Tables = append(tm.Tables, table.NewTable(tableId, playersNum))

		// TODO: ここでクライアント側にtableIdを返す処理が必要
		fmt.Fprintf(w, "Your tableId is %v\n", tableId)
	} else if len(tm.Tables) > 0 {
		if err = tm.setTableAsRequested(tableId, playersNum); err != nil {
			fmt.Fprintln(w, "Requested ID does not exist")
			return
		}
		tm.tableBalance()
	} else {
		w.Write([]byte("No table is set yet\n"))
		return
	}
}

func (tm *TournamentManager) ServeHTTP(w http.ResponseWriter, httpReq *http.Request) {
	// TDとDealerのAPIは統一するべきかよく考える必要がある
	var err error
	tm.RequestJson, err = parseJsonRequest(httpReq)
	if err != nil {
		fmt.Fprintln(w, "Failed to parse json request")
		return
	}

	isTdRequest, err := tm.RequestJson.Get("IsTDRequest").Bool()
	if err != nil {
		fmt.Fprintln(w, "Failed to assert key \"IsTDRequest\"")
	}

	if isTdRequest == true {
		//tm.handleTournamentDirectorRequest()
	} else {
		tm.handleDealerRequest(w)
	}

	// debug
	for _, v := range tm.Tables {
		fmt.Fprintf(w, "tm: {table_id:%v, players_num:%v}\n", v.ID, v.PlayersNum)
	}
}
