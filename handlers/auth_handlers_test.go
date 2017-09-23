package handlers

import (
	"encoding/json"
	"github.com/Sovianum/acquaintanceServer/config"
	"github.com/Sovianum/acquaintanceServer/dao"
	"github.com/Sovianum/acquaintanceServer/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"errors"
)

const (
	urlSample = "URL"
)

type headerPair struct {
	key   string
	value string
}

func TestEnv_UserRegisterPost_Success(t *testing.T) {
	var db, mock, dbErr = sqlmock.New()

	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer db.Close()

	var user = &model.User{
		Login:    "login",
		Password: "password",
		About:    "about",
		Sex:      model.MALE,
		Age:      100,
	}

	// mock exists
	mock.
		ExpectQuery("SELECT count").
		WithArgs(user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))

	// mock user insertion
	mock.
		ExpectExec("INSERT INTO Users").
		WithArgs(user.Login, user.Password, user.Age, user.Sex, user.About).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// mock id selection
	mock.
		ExpectQuery("SELECT id FROM").
		WithArgs(user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var env = &Env{
		userDAO:  dao.NewDBUserDAO(db),
		authConf: getAuthConf(),
	}

	var requestMsg, jsonErr = json.Marshal(user)
	assert.Nil(t, jsonErr)

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.UserRegisterPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestEnv_UserRegisterPost_CheckFail(t *testing.T) {
	var db, mock, dbErr = sqlmock.New()

	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer db.Close()

	var user = &model.User{
		Login:    "login",
		Password: "password",
		About:    "about",
		Sex:      model.MALE,
		Age:      100,
	}

	// mock exists
	mock.
	ExpectQuery("SELECT count").
		WithArgs(user.Login).
		WillReturnError(errors.New("db fail"))

	var env = &Env{
		userDAO:  dao.NewDBUserDAO(db),
		authConf: getAuthConf(),
	}

	var requestMsg, jsonErr = json.Marshal(user)
	assert.Nil(t, jsonErr)

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.UserRegisterPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestEnv_UserRegisterPost_Conflict(t *testing.T) {
	var db, mock, dbErr = sqlmock.New()

	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer db.Close()

	var user = &model.User{
		Login:    "login",
		Password: "password",
		About:    "about",
		Sex:      model.MALE,
		Age:      100,
	}

	// mock exists
	mock.
	ExpectQuery("SELECT count").
		WithArgs(user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))

	var env = &Env{
		userDAO:  dao.NewDBUserDAO(db),
		authConf: getAuthConf(),
	}

	var requestMsg, jsonErr = json.Marshal(user)
	assert.Nil(t, jsonErr)

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.UserRegisterPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestEnv_UserRegisterPost_SaveErr(t *testing.T) {
	var db, mock, dbErr = sqlmock.New()

	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer db.Close()

	var user = &model.User{
		Login:    "login",
		Password: "password",
		About:    "about",
		Sex:      model.MALE,
		Age:      100,
	}

	// mock exists
	mock.
	ExpectQuery("SELECT count").
		WithArgs(user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))

	// mock user insertion
	mock.
	ExpectExec("INSERT INTO Users").
		WithArgs(user.Login, user.Password, user.Age, user.Sex, user.About).
		WillReturnError(errors.New("db fail"))

	// mock id selection
	mock.
	ExpectQuery("SELECT id FROM").
		WithArgs(user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var env = &Env{
		userDAO:  dao.NewDBUserDAO(db),
		authConf: getAuthConf(),
	}

	var requestMsg, jsonErr = json.Marshal(user)
	assert.Nil(t, jsonErr)

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.UserRegisterPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestEnv_UserRegisterPost_IdExtraction(t *testing.T) {
	var db, mock, dbErr = sqlmock.New()

	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer db.Close()

	var user = &model.User{
		Login:    "login",
		Password: "password",
		About:    "about",
		Sex:      model.MALE,
		Age:      100,
	}

	// mock exists
	mock.
	ExpectQuery("SELECT count").
		WithArgs(user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))

	// mock user insertion
	mock.
	ExpectExec("INSERT INTO Users").
		WithArgs(user.Login, user.Password, user.Age, user.Sex, user.About).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// mock id selection
	mock.
	ExpectQuery("SELECT id FROM").
		WithArgs(user.Login).
		WillReturnError(errors.New("db fail"))

	var env = &Env{
		userDAO:  dao.NewDBUserDAO(db),
		authConf: getAuthConf(),
	}

	var requestMsg, jsonErr = json.Marshal(user)
	assert.Nil(t, jsonErr)

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.UserRegisterPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestEnv_UserSignInPost_Success(t *testing.T) {
	var db, mock, dbErr = sqlmock.New()

	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer db.Close()

	var user = &model.User{
		Login:    "login",
		Password: "password",
		About:    "about",
		Sex:      model.MALE,
		Age:      100,
	}

	// mock exists
	mock.
	ExpectQuery("SELECT count").
		WithArgs(user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))

	// mock id selection
	mock.
	ExpectQuery("SELECT id FROM").
		WithArgs(user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var env = &Env{
		userDAO:  dao.NewDBUserDAO(db),
		authConf: getAuthConf(),
	}

	var requestMsg, jsonErr = json.Marshal(user)
	assert.Nil(t, jsonErr)

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.UserSignInPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestEnv_UserSignInPost_CheckFail(t *testing.T) {
	var db, mock, dbErr = sqlmock.New()

	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer db.Close()

	var user = &model.User{
		Login:    "login",
		Password: "password",
		About:    "about",
		Sex:      model.MALE,
		Age:      100,
	}

	// mock exists
	mock.
	ExpectQuery("SELECT count").
		WithArgs(user.Login).
		WillReturnError(errors.New("db fail"))

	var env = &Env{
		userDAO:  dao.NewDBUserDAO(db),
		authConf: getAuthConf(),
	}

	var requestMsg, jsonErr = json.Marshal(user)
	assert.Nil(t, jsonErr)

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.UserSignInPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestEnv_UserSignInPost_IdExtractionFail(t *testing.T) {
	var db, mock, dbErr = sqlmock.New()

	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer db.Close()

	var user = &model.User{
		Login:    "login",
		Password: "password",
		About:    "about",
		Sex:      model.MALE,
		Age:      100,
	}

	// mock exists
	mock.
	ExpectQuery("SELECT count").
		WithArgs(user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))

	// mock id selection
	mock.
	ExpectQuery("SELECT id FROM").
		WithArgs(user.Login).
		WillReturnError(errors.New("db fail"))

	var env = &Env{
		userDAO:  dao.NewDBUserDAO(db),
		authConf: getAuthConf(),
	}

	var requestMsg, jsonErr = json.Marshal(user)
	assert.Nil(t, jsonErr)

	var rec, recErr = getRecorder(
		urlSample,
		http.MethodPost,
		env.UserSignInPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)

	assert.Nil(t, recErr)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func getAuthConf() config.AuthConfig {
	return config.AuthConfig{
		ExpireDays: 100,
		TokenKey:   "token90",
	}
}

func getRecorder(
	url string,
	method string,
	handlerFunc func(http.ResponseWriter, *http.Request),
	body io.Reader,
	headers ...headerPair,
) (*httptest.ResponseRecorder, error) {
	var req, err = http.NewRequest(
		method,
		url,
		body,
	)

	for _, hp := range headers {
		req.Header.Set(hp.key, hp.value)
	}

	if err != nil {
		return nil, err
	}

	var handler = http.HandlerFunc(handlerFunc)
	var rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	return rec, nil
}
