package commdb

import (
	"github.com/hashicorp/go-memdb"
	"time"
	"github.com/go-errors/errors"
)

type CommunicationDAO interface {
	CreateRequest(from uint, to uint) error
	GetRequests(to uint) ([]UserLink, error)
	//CreateLink(from uint, to uint) error
	//DeleteRequest(from uint, to uint)
}

func NewCommunicationDAO(db *memdb.MemDB) CommunicationDAO {
	return &communicationDAO{
		db: db,
	}
}

type communicationDAO struct {
	db *memdb.MemDB
}

func (dao *communicationDAO) CreateRequest(from uint, to uint) error {
	var txn = dao.db.Txn(true)
	var link = UserLink{
		Id1: to,
		Id2: from,
		Ts: time.Now(),
	}

	var insertErr = txn.Insert(requestTableName, link)
	if insertErr != nil {
		txn.Abort()
		return insertErr
	}

	txn.Commit()
	return nil
}

func (dao *communicationDAO) GetRequests(to uint) ([]UserLink, error) {
	var txn = dao.db.Txn(false)
	defer txn.Abort()

	var iter, iterErr = txn.Get(requestTableName, idIdx, to)
	if iterErr != nil {
		return nil, iterErr
	}

	var result = make([]UserLink, 0)
	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		var link, ok = raw.(UserLink)
		if !ok {
			return nil, errors.Errorf("Failed to convert value %v", raw)
		}

		result = append(result, link)
	}

	return result, nil
}
