package server

import (
	"database/sql"
	"fmt"
	"github.com/Sovianum/arquest-server/model"
	"github.com/Sovianum/arquest-server/mylog"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

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
	gin.SetMode(gin.ReleaseMode)
}

func (s *AuthTestSuite) TestAuthOk() {
	// mock exists
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(s.user.Id).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))

	token, tokenErr := s.env.generateTokenString(s.user.Id, s.user.Login)
	s.Require().NoError(tokenErr)

	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.CheckAuthorization,
		strings.NewReader(""),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, token},
	)
	s.Require().NoError(recErr)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *AuthTestSuite) TestAuthUserNotFound() {
	// mock exists
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(s.user.Id).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))

	token, tokenErr := s.env.generateTokenString(s.user.Id, s.user.Login)
	s.Require().NoError(tokenErr)

	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.CheckAuthorization,
		strings.NewReader(""),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, token},
	)
	s.Require().NoError(recErr)
	s.Equal(http.StatusNotFound, rec.Code)
}

func (s *AuthTestSuite) TestAuthUserNoToken() {
	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.CheckAuthorization,
		strings.NewReader(""),
		headerPair{"Content-Type", "application/json"},
	)
	s.Require().NoError(recErr)
	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *AuthTestSuite) TestAuthUserDBErr() {
	// mock exists
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(s.user.Id).
		WillReturnError(fmt.Errorf("fail"))

	token, tokenErr := s.env.generateTokenString(s.user.Id, s.user.Login)
	s.Require().NoError(tokenErr)

	rec, recErr := getRecorder(
		urlSample,
		http.MethodPost,
		s.env.CheckAuthorization,
		strings.NewReader(""),
		headerPair{"Content-Type", "application/json"},
		headerPair{authorizationStr, token},
	)
	s.Require().NoError(recErr)
	s.Equal(http.StatusInternalServerError, rec.Code)
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
