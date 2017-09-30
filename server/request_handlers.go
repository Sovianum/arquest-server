package server

import (
	"net/http"
	"io/ioutil"
	"github.com/Sovianum/acquaintanceServer/model"
	"encoding/json"
	"errors"
	"github.com/Sovianum/acquaintanceServer/common"
	"strconv"
	"github.com/patrickmn/go-cache"
	"fmt"
	"github.com/Sovianum/acquaintanceServer/dao"
)

const (
	requestNotFound = "request not found"
	alreadyAccepted = "user has already accepted another request"
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

	var requestId, dbErr = env.meetRequestDAO.CreateRequest(
		meetRequest.RequesterId, meetRequest.RequestedId, env.conf.Logic.RequestExpiration, env.conf.Logic.Distance,
	)
	if dbErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(common.GetErrorJson(dbErr))
		return
	}
	if requestId == dao.ImpossibleID {
		w.WriteHeader(http.StatusConflict)
		return
	}
	var code, err = env.handleRequestPending(requestId, userId)
	if err != nil {
		w.WriteHeader(code)
		w.Write(common.GetErrorJson(err))
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
		return
	}

	var msg, msgErr = json.Marshal(requests)
	if msgErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(common.GetErrorJson(msgErr))
		return
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
		return
	}

	// here cache is is used before accessing database
	// cos when checking incoming requests, only pending requests are taken into account
	switch update.Status {
	case model.StatusAccepted:
		var code, err = env.handleRequestAccept(update.Id, userId)
		if err != nil {
			w.WriteHeader(code)
			w.Write(common.GetErrorJson(err))
			return
		}
	case model.StatusDeclined:
		var code, err = env.handleRequestDecline(update.Id, userId)
		if err != nil {
			w.WriteHeader(code)
			w.Write(common.GetErrorJson(err))
			return
		}
	}

	var rowsAffected, dbErr = env.meetRequestDAO.UpdateRequest(update.Id, userId, update.Status)
	if dbErr != nil {
		env.revertCache(update.Id, userId)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(common.GetErrorJson(dbErr))
		return
	}

	if rowsAffected == 0 {
		env.revertCache(update.Id, userId)
		w.WriteHeader(http.StatusNotFound)
		w.Write(common.GetErrorJson(errors.New(requestNotFound)))
		return
	}
}

func (env *Env) GetNewRequests(w http.ResponseWriter, r *http.Request) {
	var userId, tokenCode, tokenErr = env.getIdFromRequest(r)
	if tokenErr != nil {
		w.WriteHeader(tokenCode)
		w.Write(common.GetErrorJson(tokenErr))
		return
	}

	var box, boxErr = env.getMailBox(userId)
	if boxErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(common.GetErrorJson(boxErr))
		return
	}

	var newRequestData = box.GetAll(env.conf.Logic.PollSeconds)
	var msg, jsonErr = json.Marshal(newRequestData)

	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(common.GetErrorJson(jsonErr))
		return
	}
	w.Write(msg)
}

func (env *Env) revertCache(requestId int, userId int) {
	var box, found = env.meetRequestCache.Get(strconv.Itoa(userId))
	if !found {
		return
	}

	box.(MailBox).Remove(requestId)
}

func (env *Env) handleRequestAccept(requestId int, userId int) (int, error) {
	var boxFunc = func(box MailBox, request *model.MeetRequest) (int, error) {
		if err := box.AddAccept(request); err != nil {
			return http.StatusUnavailableForLegalReasons, errors.New(alreadyAccepted)
		}
		return http.StatusOK, nil
	}
	var rightsCheckFunc = func(request *model.MeetRequest, userId int) bool {return request.RequestedId == userId}
	return env.dispatchRequest(boxFunc, rightsCheckFunc, requestId, userId)
}

func (env *Env) handleRequestDecline(requestId int, userId int) (int, error) {
	var boxFunc = func(box MailBox, request *model.MeetRequest) (int, error) {
		box.AddDecline(request)
		return http.StatusOK, nil
	}
	var rightsCheckFunc = func(request *model.MeetRequest, userId int) bool {return request.RequestedId == userId}
	return env.dispatchRequest(boxFunc, rightsCheckFunc, requestId, userId)
}

func (env *Env) handleRequestPending(requestId int, userId int) (int, error) {
	var boxFunc = func(box MailBox, request *model.MeetRequest) (int, error) {
		box.AddPending(request)
		return http.StatusOK, nil
	}
	var rightsCheckFunc = func(request *model.MeetRequest, userId int) bool {
		return request.RequesterId == userId
	}
	return env.dispatchRequest(boxFunc, rightsCheckFunc, requestId, userId)
}

func (env *Env) dispatchRequest(
	boxFunc func(MailBox, *model.MeetRequest) (int, error),
	rightsCheckFunc func(request *model.MeetRequest, userId int) bool,
	requestId int,
	userId int,
) (int, error) {
	var request, requestErr = env.meetRequestDAO.GetPendingRequestById(requestId)
	if requestErr != nil {
		return http.StatusNotFound, requestErr
	}

	if !rightsCheckFunc(request, userId) {
		return http.StatusNotFound, errors.New(requestNotFound)
	}

	var box, boxErr = env.getMailBox(userId)
	if boxErr != nil {
		return http.StatusInternalServerError, boxErr
	}

	return boxFunc(box, request)
}

func (env *Env) getMailBox(id int) (MailBox, error) {
	var box, found = env.meetRequestCache.Get(strconv.Itoa(id))
	if !found {
		box = NewMailBox()
		env.meetRequestCache.Set(strconv.Itoa(id), box, cache.DefaultExpiration)
	}

	var casted, ok = box.(MailBox)
	if !ok {
		return nil, fmt.Errorf("failed to cast box. Type: %T", box)
	}
	return casted, nil
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
