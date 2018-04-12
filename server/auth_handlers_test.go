package server

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Sovianum/arquest-server/config"
	"github.com/Sovianum/arquest-server/dao"
	"github.com/Sovianum/arquest-server/model"
	"github.com/Sovianum/arquest-server/mylog"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	urlSample  = "URL"
	expireDays = 100
	tokenKey   = "token90"
)

type headerPair struct {
	key   string
	value string
}

func getEnv(db *sql.DB) *Env {
	return &Env{
		userDAO: dao.NewDBUserDAO(db),
		conf:    getAuthConf(),
		hashFunc: func(password []byte) ([]byte, error) {
			var h = sha256.New()
			h.Write(password)
			return h.Sum(nil), nil
		},
		hashValidator: func(password []byte, hash []byte) error {
			var h = sha256.New()
			h.Write(password)
			var passHash = h.Sum(nil)

			if string(passHash) != string(hash) {
				return fmt.Errorf("hashes %s, %s do not match", string(passHash), string(hash))
			}
			return nil
		},
		logger: mylog.NewLogger(ioutil.Discard),
	}
}

func getAuthConf() *config.Conf {
	return &config.Conf{
		Auth: config.AuthConfig{
			ExpireDays: expireDays,
			TokenKey:   tokenKey,
		},
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

type AuthTestSuite struct {
	suite.Suite
	user *model.User
	hash []byte
	db   *sql.DB
	env  *Env
	mock sqlmock.Sqlmock
}

func (s *AuthTestSuite) SetupTest() {
	s.user = &model.User{
		Login:    "login",
		Password: "password",
		About:    "about",
		Sex:      model.MALE,
		Age:      100,
	}
	var err error
	s.db, s.mock, err = sqlmock.New()
	s.Require().NoError(err)

	s.env = getEnv(s.db)
	s.hash, _ = s.env.hashFunc([]byte(s.user.Password))
	s.env.logger = mylog.NewLogger(ioutil.Discard)
}

func (s *AuthTestSuite) TestRegisterSuccess() {
	// mock exists
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(s.user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))

	// mock user insertion
	s.mock.
		ExpectExec("INSERT INTO users").
		WithArgs(s.user.Login, string(s.hash), s.user.Age, s.user.Sex, s.user.About).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// mock id selection
	s.mock.
		ExpectQuery("SELECT id FROM").
		WithArgs(s.user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	requestMsg, jsonErr := json.Marshal(s.user)
	s.Require().NoError(jsonErr)

	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.UserRegisterPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)
	s.Require().NoError(recErr)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *AuthTestSuite) TestRegisterParseFail() {
	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.UserRegisterPost,
		strings.NewReader("Invalid json"),
		headerPair{"Content-Type", "application/json"},
	)
	s.Require().NoError(recErr)
	s.Equal(http.StatusBadRequest, rec.Code)
}

func (s *AuthTestSuite) TestRegisterCheckFail() {
	// mock exists
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(s.user.Login).
		WillReturnError(fmt.Errorf("db fail"))

	requestMsg, jsonErr := json.Marshal(s.user)
	s.Require().Nil(jsonErr)

	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.UserRegisterPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)
	s.Require().NoError(recErr)
	s.Equal(http.StatusInternalServerError, rec.Code)
}

func (s *AuthTestSuite) TestRegisterConflict() {
	// mock exists
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(s.user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))

	requestMsg, jsonErr := json.Marshal(s.user)
	s.NoError(jsonErr)

	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.UserRegisterPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)
	s.Require().NoError(recErr)
	s.Equal(http.StatusConflict, rec.Code)
}

func (s *AuthTestSuite) TestRegisterSaveErr() {
	// mock exists
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(s.user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))

	// mock user insertion
	s.mock.
		ExpectExec("INSERT INTO Users").
		WithArgs(s.user.Login, s.hash, s.user.Age, s.user.Sex, s.user.About).
		WillReturnError(fmt.Errorf("db fail"))

	// mock id selection
	s.mock.
		ExpectQuery("SELECT id FROM").
		WithArgs(s.user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	requestMsg, jsonErr := json.Marshal(s.user)
	s.Require().NoError(jsonErr)

	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.UserRegisterPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)
	s.Require().NoError(recErr)
	s.Equal(http.StatusInternalServerError, rec.Code)
}

func (s *AuthTestSuite) TestRegisterIdExtractionErr() {
	// mock exists
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(s.user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))

	// mock user insertion
	s.mock.
		ExpectExec("INSERT INTO Users").
		WithArgs(s.user.Login, s.hash, s.user.Age, s.user.Sex, s.user.About).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// mock id selection
	s.mock.
		ExpectQuery("SELECT id FROM").
		WithArgs(s.user.Login).
		WillReturnError(fmt.Errorf("db fail"))

	requestMsg, jsonErr := json.Marshal(s.user)
	s.Require().NoError(jsonErr)

	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.UserRegisterPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)
	s.Require().NoError(recErr)
	s.Equal(http.StatusInternalServerError, rec.Code)
}

func (s *AuthTestSuite) TestRegisterNoLogin() {
	s.user = &model.User{
		Password: "password",
		About:    "about",
		Sex:      model.MALE,
		Age:      100,
	}

	requestMsg, jsonErr := json.Marshal(s.user)
	s.Require().NoError(jsonErr)

	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.UserRegisterPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)
	s.Require().NoError(recErr)
	s.Equal(http.StatusBadRequest, rec.Code)
}

func (s *AuthTestSuite) TestRegisterNoPassword() {
	s.user = &model.User{
		Login: "login",
		About: "about",
		Sex:   model.MALE,
		Age:   100,
	}

	requestMsg, jsonErr := json.Marshal(s.user)
	s.Require().NoError(jsonErr)

	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.UserRegisterPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)

	s.NoError(recErr)
	s.Equal(http.StatusBadRequest, rec.Code)
}

func (s *AuthTestSuite) TestSignInSuccess() {
	// mock exists
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(s.user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))

	// mock user extraction
	s.mock.
		ExpectQuery("SELECT id").
		WithArgs(s.user.Login).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{"id", "login", "password", "age", "sex", "about"}).
				AddRow(1, "login", s.hash, 100, model.MALE, "about"),
		)

	requestMsg, jsonErr := json.Marshal(s.user)
	s.NoError(jsonErr)

	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.UserSignInPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)
	s.Require().NoError(recErr)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *AuthTestSuite) TestUserSignWrongPassword() {
	// mock exists
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(s.user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))

	// mock user extraction
	s.mock.
		ExpectQuery("SELECT id").
		WithArgs(s.user.Login).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{"id", "login", "password", "age", "sex", "about"}).
				AddRow(1, "login", "pass", 100, model.MALE, "about"),
		)

	requestMsg, jsonErr := json.Marshal(s.user)
	s.Require().NoError(jsonErr)

	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.UserSignInPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)
	s.NoError(recErr)
	s.Equal(http.StatusNotFound, rec.Code)
}

func (s *AuthTestSuite) TestSignInParseError() {
	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.UserSignInPost,
		strings.NewReader("Invalid json"),
		headerPair{"Content-Type", "application/json"},
	)
	s.Require().NoError(recErr)
	s.Equal(http.StatusBadRequest, rec.Code)
}

func (s *AuthTestSuite) TestSignInDBFail() {
	// mock exists
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(s.user.Login).
		WillReturnError(fmt.Errorf("db fail"))

	requestMsg, jsonErr := json.Marshal(s.user)
	s.Require().NoError(jsonErr)

	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.UserSignInPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)

	s.Require().NoError(recErr)
	s.Equal(http.StatusInternalServerError, rec.Code)
}

func (s *AuthTestSuite) TestSignInNotFound() {
	// mock exists
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(s.user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))

	requestMsg, jsonErr := json.Marshal(s.user)
	s.Require().NoError(jsonErr)

	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.UserSignInPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)

	s.Require().NoError(recErr)
	s.Equal(http.StatusNotFound, rec.Code)
}

func (s *AuthTestSuite) TestSignInIdExtractionFail() {
	// mock exists
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(s.user.Login).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))

	// mock id selection
	s.mock.
		ExpectQuery("SELECT id FROM").
		WithArgs(s.user.Login).
		WillReturnError(fmt.Errorf("db fail"))

	requestMsg, jsonErr := json.Marshal(s.user)
	s.Require().NoError(jsonErr)

	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.UserSignInPost,
		strings.NewReader(string(requestMsg)),
		headerPair{"Content-Type", "application/json"},
	)

	s.Require().NoError(recErr)
	s.Equal(http.StatusInternalServerError, rec.Code)
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
