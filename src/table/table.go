package table

type Table struct {
	ID, PlayersNum int
}

func NewTable(id, playersNum int) Table {
	return Table{id, playersNum}
}
