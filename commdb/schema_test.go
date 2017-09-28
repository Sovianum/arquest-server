package interdb

import (
	"testing"
	"github.com/hashicorp/go-memdb"
	"time"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestGetSchema_Smoke(t *testing.T) {
	var db, err = memdb.NewMemDB(GetSchema())
	if err != nil {
		t.Fatal(err)
	}

	var txn = db.Txn(true)
	var link = &UserLink{Id1: 1, Id2: 2, Ts: time.Now()}

	var insertErr = txn.Insert(requestTableName, link)
	assert.Nil(t, insertErr)
	txn.Commit()

	txn = db.Txn(false)
	defer txn.Abort()

	var raw, getErr = txn.First(requestTableName, idIdx, uint(1))
	assert.Nil(t, getErr)

	fmt.Println(raw.(*UserLink))
}
