package server

import (
	"net/http"
	"io/ioutil"
	"github.com/Sovianum/acquaintanceServer/model"
	"encoding/json"
	"errors"
	"github.com/Sovianum/acquaintanceServer/common"
)

const (
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
		w.WriteHeader(http.StatusConflict)
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

func (env *Env) UpdateRequest(w http.ResponseWriter, r *http.Request) {
	var update, parseCode, parseErr = parseRequestUpdate(r)
	if parseErr != nil {
		w.WriteHeader(parseCode)
		w.Write(common.GetErrorJson(parseErr))
		return
	}

	var userId, tokenCode, tokenErr = env.getIdFromRequest(r)
	if tokenErr != nil {
		w.WriteHeader(tokenCode)
		w.Write(common.GetErrorJson(tokenErr))
	}

	var rowsAffected, dbErr = env.meetRequestDAO.UpdateRequest(update.Id, userId, update.Status)
	if dbErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(common.GetErrorJson(dbErr))
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write(common.GetErrorJson(errors.New(requestNotFound)))
	}
}

func (env *Env) handleRequestAccept(update *model.MeetRequestUpdate) error {
	return nil	// TODO implement accept-specific logic
}

func (env *Env) handleRequestDecline(update *model.MeetRequestUpdate) error {
	return nil // TODO implement decline-specific logic
}

func parseRequestUpdate(r *http.Request) (*model.MeetRequestUpdate, int, error) {
	var body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := r.Body.Close(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var update = new(model.MeetRequestUpdate)
	if err := json.Unmarshal(body, &update); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return update, http.StatusOK, nil
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
