package server

import (
	"net/http"
	"io/ioutil"
	"github.com/Sovianum/acquaintanceServer/model"
	"encoding/json"
	"errors"
	"github.com/Sovianum/acquaintanceServer/common"
	"github.com/gorilla/mux"
	"strconv"
)

const (
	RequestId = "requestId"
	requestIdNotSet = "request id not set"
	badRequestIdValue = "bad request id value"
	requestNotFound = "request not found"
)

func (env *Env) CreateRequest(w http.ResponseWriter, r *http.Request) {
	var meetRequest, parseCode, parseErr = parseRequest(r)
	if parseErr != nil {
		w.WriteHeader(parseCode)
		w.Write(common.GetErrorJson(parseErr))
		return
	}
	var userId, tokenCode, tokenErr = env.getIdFromRequest(r)
	if tokenErr != nil {
		w.WriteHeader(tokenCode)
		w.Write(common.GetErrorJson(tokenErr))
		return
	}
	if userId != meetRequest.RequesterId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var rowsAffected, dbErr = env.meetRequestDAO.CreateRequest(
		meetRequest.RequesterId, meetRequest.RequestedId, env.conf.Logic.RequestExpiration, env.conf.Logic.Distance,
	)
	if dbErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(common.GetErrorJson(dbErr))
		return
	}
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}
}

func (env *Env) GetRequests(w http.ResponseWriter, r *http.Request) {
	var userId, tokenCode, tokenErr = env.getIdFromRequest(r)
	if tokenErr != nil {
		w.WriteHeader(tokenCode)
		w.Write(common.GetErrorJson(tokenErr))
		return
	}
	var requests, requestsErr = env.meetRequestDAO.GetRequests(userId)
	if requestsErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(common.GetErrorJson(requestsErr))
	}

	var msg, msgErr = json.Marshal(requests)
	if msgErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(common.GetErrorJson(msgErr))
	}
	w.Write(msg)
}

func (env *Env) AcceptRequest(w http.ResponseWriter, r *http.Request) {
	// TODO добавить проверку на занятость отправителя
	var code, err = env.updateRequestStatus(env.meetRequestDAO.AcceptRequest, r)
	if err != nil {
		w.WriteHeader(code)
		w.Write(common.GetErrorJson(err))
	}
	// TODO добавить добавление принятого канала в кэш
}

func (env *Env) DeclineRequest(w http.ResponseWriter, r *http.Request) {
	var code, err = env.updateRequestStatus(env.meetRequestDAO.DeclineRequest, r)
	if err != nil {
		w.WriteHeader(code)
		w.Write(common.GetErrorJson(err))
	}
	// TODO добавить добавление отклоненного канала в кэш
}

func (env *Env) updateRequestStatus(f func(int, int) (int, error), r *http.Request) (int, error) {
	var userId, tokenCode, tokenErr = env.getIdFromRequest(r)
	if tokenErr != nil {
		return tokenCode, tokenErr
	}

	var vars = mux.Vars(r)
	var requestId, ok = vars[RequestId]
	if !ok {
		return http.StatusBadRequest, errors.New(requestIdNotSet)
	}

	var casted, castErr = strconv.Atoi(requestId)
	if castErr != nil {
		return http.StatusBadRequest, errors.New(badRequestIdValue)
	}

	var rowsAffected, dbErr = f(casted, userId)
	if dbErr != nil {
		return http.StatusInternalServerError, dbErr
	}

	if rowsAffected == 0 {
		return http.StatusNotFound, errors.New(requestNotFound)
	}

	return http.StatusOK, nil
}

func parseRequest(r *http.Request) (*model.MeetRequest, int, error) {
	var body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := r.Body.Close(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var request = new(model.MeetRequest)
	request.RequestedId = -1
	request.RequesterId = -1

	if err := json.Unmarshal(body, &request); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if request.RequesterId == -1 || request.RequestedId == -1 {
		return nil, http.StatusBadRequest, errors.New("Empty request")
	}

	return request, http.StatusOK, nil
}
