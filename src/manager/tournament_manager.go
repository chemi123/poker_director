package manager

import (
	"errors"
	"fmt" // for debug
	"github.com/bitly/go-simplejson"
	"github.com/chemi123/poker_director/src/table"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	PLAYERSNUMKEY string = "PlayersNum"
	IDKEY         string = "ID"
	NEWTABLEKEY   string = "NewTable"
)

type TournamentManager struct {
	tables        []table.Table
	requestedJSON *simplejson.Json
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

func (tm *TournamentManager) balanceTable() {
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

func parseJSONRequest(httpReq *http.Request) (*simplejson.Json, error) {
	if httpReq.Method != http.MethodPost {
		return nil, errors.New("Only takes POST request\n")
	}

	reqBody, err := ioutil.ReadAll(httpReq.Body)
	if err != nil {
		return nil, errors.New("Failed to read request body\n")
	}

	return simplejson.NewJson(reqBody)
}

func (tm *TournamentManager) handleDealerRequest() {
	tableID, err := tm.requestedJSON.Get(IDKEY).Int()
	if err != nil {
		log.Println(err)
		return
	}

	playersNum, err := tm.requestedJSON.Get(PLAYERSNUMKEY).Int()
	if err != nil {
		log.Println(err)
		return
	}

	newTable, err := tm.requestedJSON.Get(NEWTABLEKEY).Bool()
	if err != nil {
		log.Println(err)
		return
	}

	// TODO: tableIDを付ける規則は改めて修正する必要がある
	//       でないと例えばtableIDが1, 2, 3, 4とあって2がクローズされた後にまたtableが追加されたら既に存在しているtableID=4が再度出来上がる
	//       バグを仕込む可能性が高そうな箇所である
	if newTable {
		tableID = len(tm.tables) + 1
		tm.tables = append(tm.tables, table.NewTable(tableID, playersNum))

		// TODO: ここでクライアント側にtableIDを返す処理が必要
		log.Printf("Your tableID is %v\n", tableID)
	} else if len(tm.tables) > 0 {
		if err = tm.setTableAsRequested(tableID, playersNum); err != nil {
			log.Println(err)
			return
		}
		tm.balanceTable()
	} else {
		log.Println("No table is set yet")
		return
	}
}

func (tm *TournamentManager) ServeHTTP(w http.ResponseWriter, httpReq *http.Request) {
	log.SetOutput(w)
	var err error
	tm.requestedJSON, err = parseJSONRequest(httpReq)
	if err != nil {
		log.Println(err)
		return
	}

	tm.handleDealerRequest()

	// debug
	for _, v := range tm.tables {
		fmt.Fprintf(w, "tm: {table_id:%v, players_num:%v}\n", v.ID, v.PlayersNum)
	}
}
