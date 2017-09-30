package server

import (
	"testing"
	"github.com/Sovianum/acquaintanceServer/model"
	"github.com/stretchr/testify/assert"
)

func TestMeetConnection_AddAccept(t *testing.T) {
	var conn = NewMeetConnection()
	var request = new(model.MeetRequest)

	var err1 = conn.AddAccept(request)
	assert.Nil(t, err1)

	var err2 = conn.AddAccept(request)
	assert.NotNil(t, err2)
}

func TestMeetConnection_AddDecline(t *testing.T) {
	var conn = NewMeetConnection()
	var request = new(model.MeetRequest)

	conn.AddDecline(request)

	var err2 = conn.AddAccept(request)
	assert.Nil(t, err2)
}

func TestMeetConnection_GetAll(t *testing.T) {
	var conn = NewMeetConnection()
	var request = new(model.MeetRequest)

	conn.AddDecline(request)

	var requests = conn.GetAll(1)
	assert.Equal(t, 1, len(requests))

	requests = conn.GetAll(1)
	assert.Equal(t, 0, len(requests))
}
