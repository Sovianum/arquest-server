package mocks

import (
	"errors"
	"github.com/Sovianum/acquaintanceServer/model"
)

const (
	createErr   = "create error"
	getError    = "get error"
	acceptError = "accept error"
)

type createRequestFuncType func(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (int, error)
type getRequestsFuncType func(requestedId int) ([]*model.MeetRequest, error)
type updateRequestFuncType func(id int, requestedId int) (int, error)

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

var acceptRequestSuccess updateRequestFuncType = func(int, int) (int, error) {
	return 1, nil
}
var acceptRequestNoRequest updateRequestFuncType = func(int, int) (int, error) {
	return 0, nil
}
var acceptRequestError updateRequestFuncType = func(int, int) (int, error) {
	return 0, errors.New(acceptError)
}

var declineRequestSuccess updateRequestFuncType = func(int, int) (int, error) {
	return 1, nil
}
var declineRequestNoRequest updateRequestFuncType = func(int, int) (int, error) {
	return 0, nil
}
var declineRequestError updateRequestFuncType = func(int, int) (int, error) {
	return 0, errors.New(acceptError)
}

type MeetRequestDAOMockSuccess struct{}

func (*MeetRequestDAOMockSuccess) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	return createRequestSuccess(requesterId, requestedId, requestTimeoutMin, maxDistance)
}

func (*MeetRequestDAOMockSuccess) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return getRequestsSuccess(requestedId)
}

func (*MeetRequestDAOMockSuccess) AcceptRequest(id int, requestedId int) (int, error) {
	return acceptRequestSuccess(id, requestedId)
}

func (*MeetRequestDAOMockSuccess) DeclineRequest(id int, requestedId int) (int, error) {
	return declineRequestSuccess(id, requestedId)
}

type MeetRequestDAOMockCreateConflict struct{}

func (*MeetRequestDAOMockCreateConflict) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	return createRequestConflict(requesterId, requestedId, requestTimeoutMin, maxDistance)
}

func (*MeetRequestDAOMockCreateConflict) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return getRequestsSuccess(requestedId)
}

func (*MeetRequestDAOMockCreateConflict) AcceptRequest(id int, requestedId int) (int, error) {
	return acceptRequestSuccess(id, requestedId)
}

func (*MeetRequestDAOMockCreateConflict) DeclineRequest(id int, requestedId int) (int, error) {
	return declineRequestSuccess(id, requestedId)
}

type MeetRequestDAOMockCreateError struct{}

func (*MeetRequestDAOMockCreateError) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	return createRequestError(requesterId, requestedId, requestTimeoutMin, maxDistance)
}

func (*MeetRequestDAOMockCreateError) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return getRequestsSuccess(requestedId)
}

func (*MeetRequestDAOMockCreateError) AcceptRequest(id int, requestedId int) (int, error) {
	return acceptRequestSuccess(id, requestedId)
}

func (*MeetRequestDAOMockCreateError) DeclineRequest(id int, requestedId int) (int, error) {
	return declineRequestSuccess(id, requestedId)
}

type MeetRequestDAOMockGetRequestsEmpty struct{}

func (*MeetRequestDAOMockGetRequestsEmpty) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	return createRequestSuccess(requesterId, requestedId, requestTimeoutMin, maxDistance)
}

func (*MeetRequestDAOMockGetRequestsEmpty) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return getRequestsEmpty(requestedId)
}

func (*MeetRequestDAOMockGetRequestsEmpty) AcceptRequest(id int, requestedId int) (int, error) {
	return acceptRequestSuccess(id, requestedId)
}

func (*MeetRequestDAOMockGetRequestsEmpty) DeclineRequest(id int, requestedId int) (int, error) {
	return declineRequestSuccess(id, requestedId)
}

type MeetRequestDAOMockGetRequestsError struct{}

func (*MeetRequestDAOMockGetRequestsError) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	return createRequestSuccess(requesterId, requestedId, requestTimeoutMin, maxDistance)
}

func (*MeetRequestDAOMockGetRequestsError) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return getRequestsError(requestedId)
}

func (*MeetRequestDAOMockGetRequestsError) AcceptRequest(id int, requestedId int) (int, error) {
	return acceptRequestSuccess(id, requestedId)
}

func (*MeetRequestDAOMockGetRequestsError) DeclineRequest(id int, requestedId int) (int, error) {
	return declineRequestSuccess(id, requestedId)
}

type MeetRequestDAOMockAcceptNoRequest struct{}

func (*MeetRequestDAOMockAcceptNoRequest) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	return createRequestSuccess(requesterId, requestedId, requestTimeoutMin, maxDistance)
}

func (*MeetRequestDAOMockAcceptNoRequest) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return getRequestsSuccess(requestedId)
}

func (*MeetRequestDAOMockAcceptNoRequest) AcceptRequest(id int, requestedId int) (int, error) {
	return acceptRequestNoRequest(id, requestedId)
}

func (*MeetRequestDAOMockAcceptNoRequest) DeclineRequest(id int, requestedId int) (int, error) {
	return declineRequestSuccess(id, requestedId)
}

type MeetRequestDAOMockAcceptError struct{}

func (*MeetRequestDAOMockAcceptError) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	return createRequestSuccess(requesterId, requestedId, requestTimeoutMin, maxDistance)
}

func (*MeetRequestDAOMockAcceptError) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return getRequestsSuccess(requestedId)
}

func (*MeetRequestDAOMockAcceptError) AcceptRequest(id int, requestedId int) (int, error) {
	return acceptRequestError(id, requestedId)
}

func (*MeetRequestDAOMockAcceptError) DeclineRequest(id int, requestedId int) (int, error) {
	return declineRequestSuccess(id, requestedId)
}

type MeetRequestDAOMockDeclineNoRequest struct{}

func (*MeetRequestDAOMockDeclineNoRequest) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	return createRequestSuccess(requesterId, requestedId, requestTimeoutMin, maxDistance)
}

func (*MeetRequestDAOMockDeclineNoRequest) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return getRequestsSuccess(requestedId)
}

func (*MeetRequestDAOMockDeclineNoRequest) AcceptRequest(id int, requestedId int) (int, error) {
	return acceptRequestSuccess(id, requestedId)
}

func (*MeetRequestDAOMockDeclineNoRequest) DeclineRequest(id int, requestedId int) (int, error) {
	return declineRequestNoRequest(id, requestedId)
}

type MeetRequestDAOMockDeclineError struct{}

func (*MeetRequestDAOMockDeclineError) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	return createRequestSuccess(requesterId, requestedId, requestTimeoutMin, maxDistance)
}

func (*MeetRequestDAOMockDeclineError) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return getRequestsSuccess(requestedId)
}

func (*MeetRequestDAOMockDeclineError) AcceptRequest(id int, requestedId int) (int, error) {
	return acceptRequestSuccess(id, requestedId)
}

func (*MeetRequestDAOMockDeclineError) DeclineRequest(id int, requestedId int) (int, error) {
	return declineRequestError(id, requestedId)
}
