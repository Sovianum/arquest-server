package commdb

import (
	"testing"
	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/assert"
)

func TestCommunicationDAO_CreateRequest_Success(t *testing.T) {
	var db, err = memdb.NewMemDB(GetSchema())
	if err != nil {
		t.Fatal(err)
	}

	var dao = NewCommunicationDAO(db)
	err = dao.CreateRequest(uint(1), uint(2))

	assert.Nil(t, err)

	var links, getErr = dao.GetRequests(2)
	assert.Nil(t, getErr)
	assert.Equal(t, 1, len(links))

	var link = links[0]
	assert.Equal(t, uint(2), link.Id1)
	assert.Equal(t, uint(1), link.Id2)
}


