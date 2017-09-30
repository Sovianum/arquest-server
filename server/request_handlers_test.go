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
