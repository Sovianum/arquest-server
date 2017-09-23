package server

import (
	"testing"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"github.com/stretchr/testify/assert"
	"strings"
	"github.com/Sovianum/acquaintanceServer/model"
	"github.com/Sovianum/acquaintanceServer/dao"
	"encoding/json"
	"net/http"
	"time"
	"fmt"
	"github.com/go-errors/errors"
)

func TestEnv_UserSavePositionPost_Success(t *testing.T) {
	var db, mock, dbErr = sqlmock.New()

	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer db.Close()

	var date = time.Date(2003, 10, 17, 0, 0, 0, 0, time.UTC)
	var pos = &model.Position{
		UserId:1,
		Point:model.Point{X:100, Y:200},
		Time:model.QuotedTime(date),
	}

	// mock position insertion
	mock.
	ExpectExec("INSERT INTO Position").
		WithArgs(pos.UserId, pos.Point.X, pos.Point.Y, pos.Time.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var env = &Env{
		positionDAO:  dao.NewDBPositionDAO(db),
		authConf: getAuthConf(),
	}

	var requestMsg, jsonErr = json.Marshal(pos)
	assert.Nil(t, jsonErr)

	var tokenStr, _ = env.generateTokenString(1, "login")
	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.UserSavePositionPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, fmt.Sprintf("Bearer %s", tokenStr)},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestEnv_UserSavePositionPost_Unauthorized(t *testing.T) {
	var date = time.Date(2003, 10, 17, 0, 0, 0, 0, time.UTC)
	var pos = &model.Position{
		UserId:1,
		Point:model.Point{X:100, Y:200},
		Time:model.QuotedTime(date),
	}

	var env = &Env{
		authConf: getAuthConf(),
	}

	var requestMsg, jsonErr = json.Marshal(pos)
	assert.Nil(t, jsonErr)

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.UserSavePositionPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestEnv_UserSavePositionPost_BadToken(t *testing.T) {
	var date = time.Date(2003, 10, 17, 0, 0, 0, 0, time.UTC)
	var pos = &model.Position{
		UserId:1,
		Point:model.Point{X:100, Y:200},
		Time:model.QuotedTime(date),
	}

	var env = &Env{
		authConf: getAuthConf(),
	}

	var requestMsg, jsonErr = json.Marshal(pos)
	assert.Nil(t, jsonErr)

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.UserSavePositionPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, fmt.Sprintf("Bearer some_strange_token", )},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestEnv_UserSavePositionPost_WrongId(t *testing.T) {
	var date = time.Date(2003, 10, 17, 0, 0, 0, 0, time.UTC)
	var pos = &model.Position{
		UserId:1,
		Point:model.Point{X:100, Y:200},
		Time:model.QuotedTime(date),
	}

	var env = &Env{
		authConf: getAuthConf(),
	}

	var requestMsg, jsonErr = json.Marshal(pos)
	assert.Nil(t, jsonErr)

	var tokenStr, _ = env.generateTokenString(10, "login")
	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.UserSavePositionPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, fmt.Sprintf("Bearer %s", tokenStr)},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestEnv_UserSavePositionPost_SaveErr(t *testing.T) {
	var db, mock, dbErr = sqlmock.New()

	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer db.Close()

	var date = time.Date(2003, 10, 17, 0, 0, 0, 0, time.UTC)
	var pos = &model.Position{
		UserId:1,
		Point:model.Point{X:100, Y:200},
		Time:model.QuotedTime(date),
	}

	// mock position insertion
	mock.
	ExpectExec("INSERT INTO Position").
		WithArgs(pos.UserId, pos.Point.X, pos.Point.Y, pos.Time.String()).
		WillReturnError(errors.New("Save error"))

	var env = &Env{
		positionDAO:  dao.NewDBPositionDAO(db),
		authConf: getAuthConf(),
	}

	var requestMsg, jsonErr = json.Marshal(pos)
	assert.Nil(t, jsonErr)

	var tokenStr, _ = env.generateTokenString(1, "login")
	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.UserSavePositionPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, fmt.Sprintf("Bearer %s", tokenStr)},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestEnv_getIdFromTokenString_Success(t *testing.T) {
	var env = &Env{
		authConf: getAuthConf(),
	}
	var tokenStr, _ = env.generateTokenString(1, "login")
	var token, _ = env.parseTokenString(tokenStr)
	var id, err = env.getIdFromTokenString(token)

	assert.Nil(t, err)
	assert.Equal(t, 1, id)
}

func TestEnv_parseTokenString_Success(t *testing.T) {
	var env = &Env{
		authConf: getAuthConf(),
	}
	var tokenStr, _ = env.generateTokenString(1, "login")
	var _, err = env.parseTokenString(tokenStr)

	assert.Nil(t, err)
}

func TestEnv_parseTokenString_Fail(t *testing.T) {
	var env = &Env{
		authConf: getAuthConf(),
	}
	var _, err = env.parseTokenString("Some_strange_str")

	assert.NotNil(t, err)
}

func TestRound(t *testing.T) {
	var testData = []struct{
		input float64
		expected int
	}{
		{0.999999, 1},
		{1, 1},
		{1.000000001, 1},
		{0.4999999, 0},
	}

	for i, item := range testData {
		assert.Equal(t, item.expected, round(item.input), i)
	}
}
