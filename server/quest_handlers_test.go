package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Sovianum/arquest-server/common"
	"github.com/Sovianum/arquest-server/dao"
	"github.com/Sovianum/arquest-server/model"
	"github.com/Sovianum/arquest-server/mylog"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type QuestTestSuite struct {
	suite.Suite
	user *model.User
	db   *sql.DB
	env  *Env
	mock sqlmock.Sqlmock
	c    *gin.Context
	rw   *httptest.ResponseRecorder
}

func (s *QuestTestSuite) SetupTest() {
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
	s.env.questDAO = dao.NewQuestDAO(s.db)
	s.env.logger = mylog.NewLogger(ioutil.Discard)
	gin.SetMode(gin.ReleaseMode)

	s.rw = httptest.NewRecorder()
	s.c, _ = gin.CreateTestContext(s.rw)
}

func (s *QuestTestSuite) TestAllQuestsSuccess() {
	s.mock.
		ExpectQuery("SELECT id, name").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "description", "rating"}).
				AddRow(1, "n1", "d1", 1.),
		)
	s.env.GetAllQuests(s.c)

	resp := common.ResponseMsg{}
	data := s.rw.Body.Bytes()
	json.Unmarshal(data, &resp)

	s.Require().Nil(resp.ErrMsg)

	quests, err := getResponseQuests(resp)
	s.Require().NoError(err)

	s.Equal(
		[]model.Quest{
			{ID: 1, Name: "n1", Description: "d1", Rating: 1.},
		},
		quests,
	)
	s.Equal(http.StatusOK, s.rw.Code)
}

func (s *QuestTestSuite) TestAllQuestsEmpty() {
	s.mock.
		ExpectQuery("SELECT id, name").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "description", "rating"}),
		)
	s.env.GetAllQuests(s.c)

	resp := common.ResponseMsg{}
	data := s.rw.Body.Bytes()
	json.Unmarshal(data, &resp)

	s.Require().Nil(resp.ErrMsg)

	quests, err := getResponseQuests(resp)
	s.Require().NoError(err)

	s.Equal(
		[]model.Quest{},
		quests,
	)
	s.Equal(http.StatusOK, s.rw.Code)
}

func (s *QuestTestSuite) TestAllQuestsError() {
	s.mock.
		ExpectQuery("SELECT id, name").
		WillReturnError(fmt.Errorf("fail"))
	s.env.GetAllQuests(s.c)

	resp := common.ResponseMsg{}
	data := s.rw.Body.Bytes()
	json.Unmarshal(data, &resp)

	s.NotNil(resp.ErrMsg)
	s.Nil(resp.Data)
	s.Equal(http.StatusInternalServerError, s.rw.Code)
}

func (s *QuestTestSuite) TestFinishedQuestsSuccess() {
	s.c.Set(UserID, s.user.Id)
	s.mock.
		ExpectQuery("SELECT q.id AS id, q.name AS name, q.description").
		WithArgs(s.user.Id).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "description", "rating"}).
				AddRow(1, "n1", "d1", 1.),
		)
	s.env.GetFinishedQuests(s.c)

	resp := common.ResponseMsg{}
	data := s.rw.Body.Bytes()
	json.Unmarshal(data, &resp)

	s.Require().Nil(resp.ErrMsg)

	quests, err := getResponseQuests(resp)
	s.Require().NoError(err)

	s.Equal(
		[]model.Quest{
			{ID: 1, Name: "n1", Description: "d1", Rating: 1.},
		},
		quests,
	)
	s.Equal(http.StatusOK, s.rw.Code)
}

func (s *QuestTestSuite) TestFinishedQuestsEmpty() {
	s.c.Set(UserID, s.user.Id)
	s.mock.
		ExpectQuery("SELECT q.id AS id, q.name AS name, q.description").
		WithArgs(s.user.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "rating"}))
	s.env.GetFinishedQuests(s.c)

	resp := common.ResponseMsg{}
	data := s.rw.Body.Bytes()
	json.Unmarshal(data, &resp)

	s.Require().Nil(resp.ErrMsg)

	quests, err := getResponseQuests(resp)
	s.Require().NoError(err)

	s.Equal(
		[]model.Quest{},
		quests,
	)
	s.Equal(http.StatusOK, s.rw.Code)
}

func (s *QuestTestSuite) TestFinishedQuestsError() {
	s.c.Set(UserID, s.user.Id)
	s.mock.
		ExpectQuery("SELECT q.id id, q.name name, q.description").
		WithArgs(s.user.Id).
		WillReturnError(fmt.Errorf("fail"))
	s.env.GetFinishedQuests(s.c)

	resp := common.ResponseMsg{}
	data := s.rw.Body.Bytes()
	json.Unmarshal(data, &resp)

	s.Require().NotNil(resp.ErrMsg)
	s.Equal(http.StatusInternalServerError, s.rw.Code)
}

func TestQuestTestSuite(t *testing.T) {
	suite.Run(t, new(QuestTestSuite))
}

func getResponseQuests(r common.ResponseMsg) ([]model.Quest, error) {
	items := r.Data.([]interface{})
	result := make([]model.Quest, len(items))

	for i, item := range items {
		q, e := questFromMap(item.(map[string]interface{}))
		if e != nil {
			return nil, e
		}
		result[i] = q
	}
	return result, nil
}

func questFromMap(m map[string]interface{}) (model.Quest, error) {
	q := model.Quest{}
	if id, ok := m["id"]; !ok {
		return q, fmt.Errorf("no id")
	} else {
		q.ID = int(id.(float64))
	}

	if name, ok := m["name"]; !ok {
		return q, fmt.Errorf("no name")
	} else {
		q.Name = name.(string)
	}

	if description, ok := m["description"]; !ok {
		return q, fmt.Errorf("no description")
	} else {
		q.Description = description.(string)
	}

	if rating, ok := m["rating"]; !ok {
		return q, fmt.Errorf("no rating")
	} else {
		q.Rating = float32(rating.(float64))
	}
	return q, nil
}
