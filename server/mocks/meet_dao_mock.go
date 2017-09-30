package mocks

import (
	"errors"
	"github.com/Sovianum/acquaintanceServer/model"
	"time"
)

const (
	createErr   = "create error"
	getError    = "get error"
	acceptError = "accept error"

	RequesterId = 2
	RequestedId = 3
)

type createRequestFuncType func(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (int, error)
type getRequestsFuncType func(requestedId int) ([]*model.MeetRequest, error)
type updateRequestFuncType func(id int, requestedId int, status string) (int, error)
type getRequestByIdFuncType func(id int) (*model.MeetRequest, error)

var createRequestSuccess createRequestFuncType = func(int, int, int, float64) (int, error) {
	return 1, nil
}
var createRequestConflict createRequestFuncType = func(int, int, int, float64) (int, error) {
	return 0, nil
}
var createRequestError createRequestFuncType = func(int, int, int, float64) (int, error) {
	return 0, errors.New(createErr)
}

var getRequestsSuccess getRequestsFuncType = func(requestedId int) ([]*model.MeetRequest, error) {
	return []*model.MeetRequest{
		{Id: 0, RequestedId: requestedId, RequesterId: 2, Status: model.StatusPending},
		{Id: 1, RequestedId: requestedId, RequesterId: 3, Status: model.StatusPending},
	}, nil
}
var getRequestsEmpty getRequestsFuncType = func(int) ([]*model.MeetRequest, error) {
	return []*model.MeetRequest{}, nil
}
var getRequestsError getRequestsFuncType = func(int) ([]*model.MeetRequest, error) {
	return nil, errors.New(getError)
}

var updateRequestSuccess updateRequestFuncType = func(int, int, string) (int, error) {
	return 1, nil
}
var updateRequestNoRequest updateRequestFuncType = func(int, int, string) (int, error) {
	return 0, nil
}
var updateRequestError updateRequestFuncType = func(int, int, string) (int, error) {
	return 0, errors.New(acceptError)
}

var getRequestByIdSuccess getRequestByIdFuncType = func(id int) (*model.MeetRequest, error) {
	return &model.MeetRequest{
		Time:        model.QuotedTime(time.Now()),
		Status:      model.StatusPending,
		RequestedId: RequestedId,
		RequesterId: RequesterId,
		Id:          id,
	}, nil
}
var getRequestByIdNotFound getRequestByIdFuncType = func(id int) (*model.MeetRequest, error) {
	return nil, errors.New("not found")
}

type MeetRequestDAOMockSuccess struct{}

func (*MeetRequestDAOMockSuccess) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	return createRequestSuccess(requesterId, requestedId, requestTimeoutMin, maxDistance)
}

func (*MeetRequestDAOMockSuccess) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return getRequestsSuccess(requestedId)
}

func (*MeetRequestDAOMockSuccess) UpdateRequest(id int, requestedId int, status string) (int, error) {
	return updateRequestSuccess(id, requestedId, status)
}

func (*MeetRequestDAOMockSuccess) GetRequestById(id int) (*model.MeetRequest, error) {
	return getRequestByIdSuccess(id)
}

type MeetRequestDAOMockCreateConflict struct{}

func (*MeetRequestDAOMockCreateConflict) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	return createRequestConflict(requesterId, requestedId, requestTimeoutMin, maxDistance)
}

func (*MeetRequestDAOMockCreateConflict) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return getRequestsSuccess(requestedId)
}

func (*MeetRequestDAOMockCreateConflict) UpdateRequest(id int, requestedId int, status string) (int, error) {
	return updateRequestSuccess(id, requestedId, status)
}

func (*MeetRequestDAOMockCreateConflict) GetRequestById(id int) (*model.MeetRequest, error) {
	return getRequestByIdSuccess(id)
}

type MeetRequestDAOMockCreateError struct{}

func (*MeetRequestDAOMockCreateError) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	return createRequestError(requesterId, requestedId, requestTimeoutMin, maxDistance)
}

func (*MeetRequestDAOMockCreateError) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return getRequestsSuccess(requestedId)
}

func (*MeetRequestDAOMockCreateError) UpdateRequest(id int, requestedId int, status string) (int, error) {
	return updateRequestSuccess(id, requestedId, status)
}

func (*MeetRequestDAOMockCreateError) GetRequestById(id int) (*model.MeetRequest, error) {
	return getRequestByIdSuccess(id)
}

type MeetRequestDAOMockGetRequestsEmpty struct{}

func (*MeetRequestDAOMockGetRequestsEmpty) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	return createRequestSuccess(requesterId, requestedId, requestTimeoutMin, maxDistance)
}

func (*MeetRequestDAOMockGetRequestsEmpty) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return getRequestsEmpty(requestedId)
}

func (*MeetRequestDAOMockGetRequestsEmpty) UpdateRequest(id int, requestedId int, status string) (int, error) {
	return updateRequestSuccess(id, requestedId, status)
}

func (*MeetRequestDAOMockGetRequestsEmpty) GetRequestById(id int) (*model.MeetRequest, error) {
	return getRequestByIdSuccess(id)
}

type MeetRequestDAOMockGetRequestsError struct{}

func (*MeetRequestDAOMockGetRequestsError) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	return createRequestSuccess(requesterId, requestedId, requestTimeoutMin, maxDistance)
}

func (*MeetRequestDAOMockGetRequestsError) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return getRequestsError(requestedId)
}

func (*MeetRequestDAOMockGetRequestsError) UpdateRequest(id int, requestedId int, status string) (int, error) {
	return updateRequestSuccess(id, requestedId, status)
}

func (*MeetRequestDAOMockGetRequestsError) GetRequestById(id int) (*model.MeetRequest, error) {
	return getRequestByIdSuccess(id)
}

type MeetRequestDAOMockUpdateNoRequest struct{}

func (*MeetRequestDAOMockUpdateNoRequest) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	return createRequestSuccess(requesterId, requestedId, requestTimeoutMin, maxDistance)
}

func (*MeetRequestDAOMockUpdateNoRequest) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return getRequestsSuccess(requestedId)
}

func (*MeetRequestDAOMockUpdateNoRequest) UpdateRequest(id int, requestedId int, status string) (int, error) {
	return updateRequestNoRequest(id, requestedId, status)
}

func (*MeetRequestDAOMockUpdateNoRequest) GetRequestById(id int) (*model.MeetRequest, error) {
	return getRequestByIdSuccess(id)
}

type MeetRequestDAOMockUpdateError struct{}

func (*MeetRequestDAOMockUpdateError) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	return createRequestSuccess(requesterId, requestedId, requestTimeoutMin, maxDistance)
}

func (*MeetRequestDAOMockUpdateError) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return getRequestsSuccess(requestedId)
}

func (*MeetRequestDAOMockUpdateError) UpdateRequest(id int, requestedId int, status string) (int, error) {
	return updateRequestError(id, requestedId, status)
}

func (*MeetRequestDAOMockUpdateError) GetRequestById(id int) (*model.MeetRequest, error) {
	return getRequestByIdSuccess(id)
}

type MeetRequestDAOMockGetRequestByIdNotFound struct{}

func (*MeetRequestDAOMockGetRequestByIdNotFound) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	return createRequestSuccess(requesterId, requestedId, requestTimeoutMin, maxDistance)
}

func (*MeetRequestDAOMockGetRequestByIdNotFound) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return getRequestsSuccess(requestedId)
}

func (*MeetRequestDAOMockGetRequestByIdNotFound) UpdateRequest(id int, requestedId int, status string) (int, error) {
	return updateRequestSuccess(id, requestedId, status)
}

func (*MeetRequestDAOMockGetRequestByIdNotFound) GetRequestById(id int) (*model.MeetRequest, error) {
	return getRequestByIdNotFound(id)
}
