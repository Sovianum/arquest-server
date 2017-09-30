package server

import (
	"github.com/Sovianum/acquaintanceServer/model"
	"github.com/go-errors/errors"
	"sync"
	"time"
)

const (
	userHasAlreadyAcceptedRequest = "user has already accepted request"
)

func NewMeetConnection() MeetConnection {
	return &meetConnection{
		syncChan:     make(chan int, 1),
		accepted:     false,
		requestsLock: sync.RWMutex{},
		acceptedLock: sync.RWMutex{},
		requests:     make([]*model.MeetRequest, 0),
	}
}

type MeetConnection interface {
	AddAccept(request *model.MeetRequest) error
	AddDecline(request *model.MeetRequest)
	GetAll(seconds int) []*model.MeetRequest
}

type meetConnection struct {
	syncChan     chan int
	accepted     bool
	acceptedLock sync.RWMutex
	requests     []*model.MeetRequest
	requestsLock sync.RWMutex
}

func (conn *meetConnection) AddAccept(request *model.MeetRequest) error {
	conn.acceptedLock.Lock()
	if conn.accepted {
		return errors.New(userHasAlreadyAcceptedRequest)
	} else {
		conn.accepted = true
	}
	conn.acceptedLock.Unlock()

	conn.requestsLock.Lock()
	var requestCopy = new(model.MeetRequest)
	*requestCopy = *request
	conn.requests = append(conn.requests, requestCopy)
	conn.requestsLock.Unlock()

	select {
	case conn.syncChan <- 1:
	default:
	}

	return nil
}

func (conn *meetConnection) AddDecline(request *model.MeetRequest) {
	conn.requestsLock.Lock()
	var requestCopy = new(model.MeetRequest)
	*requestCopy = *request
	conn.requests = append(conn.requests, requestCopy)
	conn.requestsLock.Unlock()

	select {
	case conn.syncChan <- 1:
	default:
	}
}

func (conn *meetConnection) GetAll(seconds int) []*model.MeetRequest {
	var result = make([]*model.MeetRequest, 0)

	select {
	case <-conn.syncChan:
		conn.requestsLock.Lock()
		result = conn.requests
		conn.requests = make([]*model.MeetRequest, 0)
		conn.requestsLock.Unlock()
	case <-time.After(time.Second * time.Duration(seconds)):
	}

	conn.acceptedLock.Lock()
	conn.accepted = false
	conn.acceptedLock.Unlock()
	return result
}
