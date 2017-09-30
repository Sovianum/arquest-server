package server

import (
	"github.com/Sovianum/acquaintanceServer/config"
	"testing"
	"github.com/stretchr/testify/assert"
	"strings"
	"encoding/json"
	"net/http"
	"github.com/Sovianum/acquaintanceServer/server/mocks"
	"github.com/Sovianum/acquaintanceServer/model"
	"fmt"
	"time"
	"github.com/dgrijalva/jwt-go"
)

func TestEnv_CreateRequest_Success(t *testing.T) {
	var meetRequest = model.MeetRequest{RequesterId:1, RequestedId:2}
	var requestMsg, jsonErr = json.Marshal(meetRequest)
	assert.Nil(t, jsonErr)

	var env = &Env{conf:getTotalConf(), meetRequestDAO:&mocks.MeetRequestDAOMockSuccess{}}
	var tokenStr, _ = env.generateTokenString(1, "login")

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.CreateRequest,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, fmt.Sprintf("Bearer %s", tokenStr)},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestEnv_CreateRequest_NoIdInToken(t *testing.T) {
	var env = &Env{conf:getTotalConf(), meetRequestDAO:&mocks.MeetRequestDAOMockSuccess{}}
	var tokenStr, _ = getIncompleteToken(env)

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.CreateRequest,
		strings.NewReader(""),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, fmt.Sprintf("Bearer %s", tokenStr)},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestEnv_CreateRequest_BadToken(t *testing.T) {
	var env = &Env{conf:getTotalConf(), meetRequestDAO:&mocks.MeetRequestDAOMockSuccess{}}
	var tokenStr = "Bad token"

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.CreateRequest,
		strings.NewReader(""),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, fmt.Sprintf("Bearer %s", tokenStr)},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestEnv_CreateRequest_WrongRequesterId(t *testing.T) {
	var meetRequest = model.MeetRequest{RequesterId:3, RequestedId:2}
	var requestMsg, jsonErr = json.Marshal(meetRequest)
	assert.Nil(t, jsonErr)

	var env = &Env{conf:getTotalConf(), meetRequestDAO:&mocks.MeetRequestDAOMockSuccess{}}
	var tokenStr, _ = env.generateTokenString(1, "login")

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.CreateRequest,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, fmt.Sprintf("Bearer %s", tokenStr)},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestEnv_CreateRequest_Conflict(t *testing.T) {
	var meetRequest = model.MeetRequest{RequesterId:1, RequestedId:2}
	var requestMsg, jsonErr = json.Marshal(meetRequest)
	assert.Nil(t, jsonErr)

	var env = &Env{conf:getTotalConf(), meetRequestDAO:&mocks.MeetRequestDAOMockCreateConflict{}}
	var tokenStr, _ = env.generateTokenString(1, "login")

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.CreateRequest,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, fmt.Sprintf("Bearer %s", tokenStr)},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestEnv_CreateRequest_Error(t *testing.T) {
	var meetRequest = model.MeetRequest{RequesterId:1, RequestedId:2}
	var requestMsg, jsonErr = json.Marshal(meetRequest)
	assert.Nil(t, jsonErr)

	var env = &Env{conf:getTotalConf(), meetRequestDAO:&mocks.MeetRequestDAOMockCreateError{}}
	var tokenStr, _ = env.generateTokenString(1, "login")

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.CreateRequest,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, fmt.Sprintf("Bearer %s", tokenStr)},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestEnv_GetRequests_Success(t *testing.T) {
	var env = &Env{conf:getTotalConf(), meetRequestDAO:&mocks.MeetRequestDAOMockSuccess{}}
	var tokenStr, _ = env.generateTokenString(1, "login")

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodGet,
		env.GetRequests,
		strings.NewReader(""),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, fmt.Sprintf("Bearer %s", tokenStr)},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusOK, rec.Code)

	var requests, _ = env.meetRequestDAO.GetRequests(1)
	var gotRequests = make([]*model.MeetRequest, 0)
	var jsonErr = json.Unmarshal(rec.Body.Bytes(), &gotRequests)

	assert.Nil(t, jsonErr)
	assert.Equal(t, len(requests), len(gotRequests))
	for i := range requests {
		assert.Equal(t, *requests[i], *gotRequests[i])
	}
}

func TestEnv_GetRequests_Empty(t *testing.T) {
	var env = &Env{conf:getTotalConf(), meetRequestDAO:&mocks.MeetRequestDAOMockGetRequestsEmpty{}}
	var tokenStr, _ = env.generateTokenString(1, "login")

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodGet,
		env.GetRequests,
		strings.NewReader(""),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, fmt.Sprintf("Bearer %s", tokenStr)},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusOK, rec.Code)

	var requests, _ = env.meetRequestDAO.GetRequests(1)
	var gotRequests = make([]*model.MeetRequest, 0)
	var jsonErr = json.Unmarshal(rec.Body.Bytes(), &gotRequests)

	assert.Nil(t, jsonErr)
	assert.Equal(t, len(requests), len(gotRequests))
	for i := range requests {
		assert.Equal(t, *requests[i], *gotRequests[i])
	}
}

func TestEnv_GetRequests_NoIdInToken(t *testing.T) {
	var env = &Env{conf:getTotalConf(), meetRequestDAO:&mocks.MeetRequestDAOMockSuccess{}}
	var tokenStr, _ = getIncompleteToken(env)

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodGet,
		env.GetRequests,
		strings.NewReader(""),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, fmt.Sprintf("Bearer %s", tokenStr)},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestEnv_GetRequests_BadToken(t *testing.T) {
	var env = &Env{conf:getTotalConf(), meetRequestDAO:&mocks.MeetRequestDAOMockSuccess{}}
	var tokenStr = "Bad token"

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodGet,
		env.GetRequests,
		strings.NewReader(""),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, fmt.Sprintf("Bearer %s", tokenStr)},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestEnv_GetRequests_Error(t *testing.T) {
	var env = &Env{conf:getTotalConf(), meetRequestDAO:&mocks.MeetRequestDAOMockGetRequestsError{}}
	var tokenStr, _ = env.generateTokenString(1, "login")

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodGet,
		env.GetRequests,
		strings.NewReader(""),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, fmt.Sprintf("Bearer %s", tokenStr)},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestEnv_UpdateRequest_Success(t *testing.T) {
	var env = &Env{conf:getTotalConf(), meetRequestDAO:&mocks.MeetRequestDAOMockSuccess{}}
	var tokenStr, _ = env.generateTokenString(1, "login")

	var update = model.MeetRequestUpdate{Id:1, Status:model.StatusAccepted}
	var requestMsg, err = json.Marshal(update)
	assert.Nil(t, err)

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.UpdateRequest,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, fmt.Sprintf("Bearer %s", tokenStr)},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func getIncompleteToken(env *Env) (string, error) {
	var token = jwt.New(jwt.SigningMethodHS256)
	var claims = token.Claims.(jwt.MapClaims)

	claims[loginStr] = "login"
	claims[expStr] = time.Now().Add(time.Hour * 24 * time.Duration(env.conf.Auth.ExpireDays)).Unix()

	var tokenKey = env.conf.Auth.GetTokenKey()
	return token.SignedString(tokenKey)
}

func getTotalConf() config.Conf {
	return config.Conf{
		Auth: config.AuthConfig{
			ExpireDays: expireDays,
			TokenKey:   tokenKey,
		},
		Logic: config.LogicConfig{
			OnlineTimeout: onlineTimeout,
			Distance:      distance,
		},
	}
}
