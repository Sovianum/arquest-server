package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Sovianum/arquest-server/common"
	"github.com/Sovianum/arquest-server/sqldao"
	"github.com/Sovianum/arquest-server/model"
	"github.com/Sovianum/arquest-server/mylog"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MarkTestSuite struct {
	suite.Suite
	user *model.User
	db   *sql.DB
	env  *Env
	mock sqlmock.Sqlmock
	c    *gin.Context
	rw   *httptest.ResponseRecorder
}

func (s *MarkTestSuite) SetupTest() {
	s.user = &model.User{
		Id:       20,
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
	s.env.markDAO = sqldao.NewMarkDAO(s.db)
	s.env.questDAO = sqldao.NewQuestDAO(s.db)
	s.env.logger = mylog.NewLogger(ioutil.Discard)
	gin.SetMode(gin.ReleaseMode)

	s.rw = httptest.NewRecorder()
	s.c, _ = gin.CreateTestContext(s.rw)
	s.c.Set(UserID, s.user.Id)
}

func (s *MarkTestSuite) TestFinishQuestSuccess() {
	mark := model.Mark{UserID: 20, QuestID: 2}
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(mark.QuestID).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))
	s.mock.
		ExpectExec("INSERT INTO quest_user_link").
		WithArgs(mark.UserID, mark.QuestID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	msg, err := json.Marshal(mark)
	s.Require().NoError(err)
	s.c.Request, err = getRequest(urlSample, http.MethodPost, strings.NewReader(string(msg)))

	s.env.FinishQuest(s.c)

	resp := common.ResponseMsg{}
	data := s.rw.Body.Bytes()
	json.Unmarshal(data, &resp)

	s.Require().Nil(resp.ErrMsg)
	s.Equal(http.StatusOK, s.rw.Code)
}

func (s *MarkTestSuite) TestFinishQuestForbidden() {
	mark := model.Mark{UserID: 1, QuestID: 2}
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(mark.QuestID).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))
	s.mock.
		ExpectExec("UPDATE quest_user_link SET finished = TRUE").
		WithArgs(mark.UserID, mark.QuestID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	msg, err := json.Marshal(mark)
	s.Require().NoError(err)
	s.c.Request, err = getRequest(urlSample, http.MethodPost, strings.NewReader(string(msg)))

	s.env.FinishQuest(s.c)

	resp := common.ResponseMsg{}
	data := s.rw.Body.Bytes()
	json.Unmarshal(data, &resp)

	s.Require().NotNil(resp.ErrMsg)
	s.Equal(http.StatusForbidden, s.rw.Code)
}

func (s *MarkTestSuite) TestFinishQuestError() {
	mark := model.Mark{UserID: 20, QuestID: 2}
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(mark.QuestID).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))
	s.mock.
		ExpectExec("UPDATE quest_user_link SET finished = TRUE").
		WithArgs(mark.UserID, mark.QuestID).
		WillReturnError(fmt.Errorf("fail"))

	msg, err := json.Marshal(mark)
	s.Require().NoError(err)
	s.c.Request, err = getRequest(urlSample, http.MethodPost, strings.NewReader(string(msg)))

	s.env.FinishQuest(s.c)
	s.Equal(http.StatusInternalServerError, s.rw.Code)
}

func (s *MarkTestSuite) TestMarkQuestSuccess() {
	mark := model.Mark{UserID: 20, QuestID: 2, Mark: 4}
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(mark.QuestID).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))
	s.mock.
		ExpectExec("UPDATE quest_user_link SET mark").
		WithArgs(mark.Mark, mark.UserID, mark.QuestID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.
		ExpectExec("UPDATE quest SET").
		WithArgs(mark.QuestID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	msg, err := json.Marshal(mark)
	s.Require().NoError(err)
	s.c.Request, err = getRequest(urlSample, http.MethodPost, strings.NewReader(string(msg)))

	s.env.MarkQuest(s.c)

	resp := common.ResponseMsg{}
	data := s.rw.Body.Bytes()
	json.Unmarshal(data, &resp)

	s.Require().Nil(resp.ErrMsg)
	s.Equal(http.StatusOK, s.rw.Code)
}

func (s *MarkTestSuite) TestMarkQuestForbidden() {
	mark := model.Mark{UserID: 1, QuestID: 2, Mark: 4}
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(mark.QuestID).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))
	s.mock.
		ExpectExec("UPDATE quest_user_link SET mark").
		WithArgs(mark.Mark, mark.UserID, mark.QuestID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.
		ExpectExec("UPDATE quest SET rating = mark_count").
		WithArgs(mark.Mark, mark.QuestID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	msg, err := json.Marshal(mark)
	s.Require().NoError(err)
	s.c.Request, err = getRequest(urlSample, http.MethodPost, strings.NewReader(string(msg)))

	s.env.MarkQuest(s.c)
	s.Equal(http.StatusForbidden, s.rw.Code)
}

func (s *MarkTestSuite) TestMarkQuestMarkErr() {
	mark := model.Mark{UserID: 20, QuestID: 2, Mark: 4}
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(mark.QuestID).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))
	s.mock.
		ExpectExec("UPDATE quest_user_link SET mark").
		WithArgs(mark.Mark, mark.UserID, mark.QuestID).
		WillReturnError(fmt.Errorf("mark err"))

	s.mock.
		ExpectExec("UPDATE quest SET rating = mark_count").
		WithArgs(mark.Mark, mark.QuestID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	msg, err := json.Marshal(mark)
	s.Require().NoError(err)
	s.c.Request, err = getRequest(urlSample, http.MethodPost, strings.NewReader(string(msg)))

	s.env.MarkQuest(s.c)

	resp := common.ResponseMsg{}
	data := s.rw.Body.Bytes()
	json.Unmarshal(data, &resp)

	s.Require().NotNil(resp.ErrMsg)
	s.Equal("mark err", resp.ErrMsg)
	s.Equal(http.StatusInternalServerError, s.rw.Code)
}

func (s *MarkTestSuite) TestMarkQuestRatingErr() {
	mark := model.Mark{UserID: 20, QuestID: 2, Mark: 4}
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(mark.QuestID).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))
	s.mock.
		ExpectExec("UPDATE quest_user_link SET mark").
		WithArgs(mark.Mark, mark.UserID, mark.QuestID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.
		ExpectExec("UPDATE quest SET").
		WithArgs(mark.QuestID).
		WillReturnError(fmt.Errorf("rating err"))

	msg, err := json.Marshal(mark)
	s.Require().NoError(err)
	s.c.Request, err = getRequest(urlSample, http.MethodPost, strings.NewReader(string(msg)))

	s.env.MarkQuest(s.c)

	resp := common.ResponseMsg{}
	data := s.rw.Body.Bytes()
	json.Unmarshal(data, &resp)

	s.Require().NotNil(resp.ErrMsg)
	s.Equal("rating err", resp.ErrMsg)
	s.Equal(http.StatusInternalServerError, s.rw.Code)
}

func TestMarkTestSuite(t *testing.T) {
	suite.Run(t, new(MarkTestSuite))
}
