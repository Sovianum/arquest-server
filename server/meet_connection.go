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

type requestMapType map[int]*model.MeetRequest

func NewMeetConnection() MeetConnection {
	return &meetConnection{
		syncChan:     make(chan int, 1),
		accepted:     false,
		requestsLock: sync.RWMutex{},
		acceptedLock: sync.RWMutex{},
		requestMap:   make(requestMapType),
	}
}

type MeetConnection interface {
	AddAccept(request *model.MeetRequest) error
	AddDecline(request *model.MeetRequest)
	AddPending(request *model.MeetRequest)
	Remove(requestId int)
	GetAll(seconds int) []*model.MeetRequest
}

type meetConnection struct {
	syncChan     chan int
	accepted     bool
	acceptedLock sync.RWMutex
	requestMap   requestMapType
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
	requestCopy.Status = model.StatusAccepted
	*requestCopy = *request

	conn.requestMap[request.Id] = requestCopy
	conn.requestsLock.Unlock()

	select {
	case conn.syncChan <- 1:
	default:
	}

	return nil
}

func (conn *meetConnection) AddDecline(request *model.MeetRequest) {
	conn.addNonAccept(request, model.StatusDeclined)
}

func (conn *meetConnection) AddPending(request *model.MeetRequest) {
	conn.addNonAccept(request, model.StatusPending)
}

func (conn *meetConnection) Remove(requestId int) {
	delete(conn.requestMap, requestId)
}

func (conn *meetConnection) GetAll(seconds int) []*model.MeetRequest {
	var result = make([]*model.MeetRequest, 0)

	select {
	case <-conn.syncChan:
		conn.requestsLock.Lock()
		for _, request := range conn.requestMap {
			result= append(result, request)
		}
		conn.requestMap = make(requestMapType)
		conn.requestsLock.Unlock()
	case <-time.After(time.Second * time.Duration(seconds)):
	}

	conn.acceptedLock.Lock()
	conn.accepted = false
	conn.acceptedLock.Unlock()
	return result
}

func (conn *meetConnection) addNonAccept(request *model.MeetRequest, status string) {
	conn.requestsLock.Lock()
	var requestCopy = new(model.MeetRequest)
	*requestCopy = *request
	requestCopy.Status = status

	conn.requestMap[request.Id] = requestCopy
	conn.requestsLock.Unlock()

	select {
	case conn.syncChan <- 1:
	default:
	}
}
