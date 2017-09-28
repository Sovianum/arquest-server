package interdb

import (
	"github.com/hashicorp/go-memdb"
	"time"
)

const (
	requestTableName = "request"
	meetingTableName = "meeting"

	idIdx = "id"
)

type UserLink struct {
	Id1 uint
	Id2 uint
	Ts time.Time
}

func GetSchema() *memdb.DBSchema {
	var requestTable = &memdb.TableSchema{
		Name: requestTableName,
		Indexes: map[string]*memdb.IndexSchema {
			idIdx: {
				Name: idIdx,
				Unique: true,
				Indexer: &memdb.UintFieldIndex{Field: "Id1"},
			},
		},
	}

	var meetingTable = &memdb.TableSchema{
		Name: meetingTableName,
		Indexes: map[string]*memdb.IndexSchema {
			idIdx: {
				Name: idIdx,
				Unique: true,
				Indexer: &memdb.UintFieldIndex{Field: "Id1"},
			},
		},
	}

	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema {
			requestTableName: requestTable,
			meetingTableName: meetingTable,
		},
	}
}