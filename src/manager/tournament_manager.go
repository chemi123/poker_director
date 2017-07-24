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
	tables               []table.Table
	tournamentDirectorID int
	requestedJson        *simplejson.Json
}

// リクエストで指定されたテーブルに値をbodyのjson通りにセットする
// もし指定のTable IDが存在しなければエラーを返す
// この関数が正しく動作するためにはTable IDは全てユニークでなければならない
func (tm *TournamentManager) setTableAsRequested(requestedTableId int, requestedPlayersNum int) error {
	for i, _ := range tm.tables {
		if tm.tables[i].ID == requestedTableId {
			tm.tables[i].PlayersNum = requestedPlayersNum
			return nil
		}
	}
	return errors.New("Requested Table ID does not exist")
}

// tableBalance
func (tm *TournamentManager) tableBalance() {
	minTable, maxTable := &tm.tables[0], &tm.tables[0]

	// TODO: 大した計算量ないからひとまずは愚直に計算する
	//       もっと効率化はできるがやるならバグに気をつけないといけないから費用対効果としては微妙かも。十分速いし
	for {
		for i, _ := range tm.tables {
			if minTable.PlayersNum > tm.tables[i].PlayersNum {
				minTable = &tm.tables[i]
			}

			if maxTable.PlayersNum < tm.tables[i].PlayersNum {
				maxTable = &tm.tables[i]
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
	tableId, err := tm.requestedJson.Get("ID").Int()
	if err != nil {
		fmt.Println("Failed to assert key \"ID\"")
		return
	}

	playersNum, err := tm.requestedJson.Get("PlayersNum").Int()
	if err != nil {
		fmt.Println("Failed to assert key \"PlayersNum\"")
		return
	}

	// TODO: tableIdを付ける規則は改めて修正する必要がある
	//       でないと例えばtableIdが1, 2, 3, 4とあって2がクローズされた後にまたtableが追加されたら既に存在しているtableId=4が再度出来上がる
	//       バグを仕込む可能性が高そうな箇所である
	if tableId == 0 {
		tableId = len(tm.tables) + 1
		tm.tables = append(tm.tables, table.NewTable(tableId, playersNum))

		// TODO: ここでクライアント側にtableIdを返す処理が必要
		fmt.Fprintf(w, "Your tableId is %v\n", tableId)
	} else if len(tm.tables) > 0 {
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

// TODO: これから実装
func (tm *TournamentManager) handleTournamentDirectorRequest() {
	if tm.tournamentDirectorID == 0 {
	}
}

func (tm *TournamentManager) ServeHTTP(w http.ResponseWriter, httpReq *http.Request) {
	// TDとDealerのAPIは統一するべきかよく考える必要がある
	var err error
	tm.requestedJson, err = parseJsonRequest(httpReq)
	if err != nil {
		fmt.Fprintln(w, "Failed to parse json request")
		return
	}

	isTdRequest, err := tm.requestedJson.Get("IsTDRequest").Bool()
	if err != nil {
		fmt.Fprintln(w, "Failed to assert key \"IsTDRequest\"")
	}

	if isTdRequest == true {
		tm.handleTournamentDirectorRequest()
	} else {
		tm.handleDealerRequest(w)
	}

	// debug
	for _, v := range tm.tables {
		fmt.Fprintf(w, "tm: {table_id:%v, players_num:%v}\n", v.ID, v.PlayersNum)
	}
}
