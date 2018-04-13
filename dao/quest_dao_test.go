package dao

import (
	"database/sql"
	"fmt"
	"github.com/Sovianum/arquest-server/model"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

type QuestTestSuite struct {
	suite.Suite
	db       *sql.DB
	mock     sqlmock.Sqlmock
	questDAO QuestDAO
}

func (s *QuestTestSuite) SetupTest() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	s.Require().NoError(err)
	s.questDAO = NewQuestDAO(s.db)
}

func (s *QuestTestSuite) TestAllOk() {
	rows := sqlmock.NewRows([]string{"id", "name", "description", "rating"}).
		AddRow(1, "n1", "d1", 1).
		AddRow(2, "n2", "d2", 2)

	s.mock.
		ExpectQuery("SELECT id").
		WillReturnRows(rows)

	quests, err := s.questDAO.GetAllQuests()
	s.Require().NoError(err)
	s.Equal(
		[]model.Quest{
			{ID: 1, Name: "n1", Description: "d1", Rating: 1},
			{ID: 2, Name: "n2", Description: "d2", Rating: 2},
		},
		quests,
	)
}

func (s *QuestTestSuite) TestAllEmpty() {
	rows := sqlmock.NewRows([]string{"id", "name", "description", "rating"})

	s.mock.
		ExpectQuery("SELECT id").
		WillReturnRows(rows)

	quests, err := s.questDAO.GetAllQuests()
	s.Require().NoError(err)
	s.Equal(
		[]model.Quest{},
		quests,
	)
}

func (s *QuestTestSuite) TestAllError() {
	s.mock.
		ExpectQuery("SELECT id").
		WillReturnError(fmt.Errorf("fail"))

	_, err := s.questDAO.GetAllQuests()
	s.Require().Error(err)
	s.Equal("fail", err.Error())
}

func (s *QuestTestSuite) TestFinishedOk() {
	rows := sqlmock.NewRows([]string{"id", "name", "description", "rating"}).
		AddRow(1, "n1", "d1", 1).
		AddRow(2, "n2", "d2", 2)

	s.mock.
		ExpectQuery("SELECT q.id .+ quest_user_link").
		WithArgs(10).
		WillReturnRows(rows)

	quests, err := s.questDAO.GetFinishedQuests(10)
	s.Require().NoError(err)
	s.Equal(
		[]model.Quest{
			{ID: 1, Name: "n1", Description: "d1", Rating: 1},
			{ID: 2, Name: "n2", Description: "d2", Rating: 2},
		},
		quests,
	)
}

func (s *QuestTestSuite) TestFinishedEmpty() {
	rows := sqlmock.NewRows([]string{"id", "name", "description", "rating"})

	s.mock.
		ExpectQuery("SELECT q.id .+ quest_user_link").
		WithArgs(10).
		WillReturnRows(rows)

	quests, err := s.questDAO.GetFinishedQuests(10)
	s.Require().NoError(err)
	s.Equal(
		[]model.Quest{},
		quests,
	)
}

func (s *QuestTestSuite) TestFinishedError() {
	s.mock.
		ExpectQuery("SELECT q.id .+ quest_user_link").
		WithArgs(10).
		WillReturnError(fmt.Errorf("fail"))

	_, err := s.questDAO.GetFinishedQuests(10)
	s.Require().Error(err)
	s.Equal("fail", err.Error())
}

func TestQuestTestSuite(t *testing.T) {
	suite.Run(t, new(QuestTestSuite))
}
