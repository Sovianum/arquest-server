package dao

import (
	"database/sql"
	"fmt"
	"github.com/Sovianum/arquest-server/model"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

type MarkTestSuite struct {
	suite.Suite
	db      *sql.DB
	mock    sqlmock.Sqlmock
	markDAO MarkDAO
}

func (s *MarkTestSuite) SetupTest() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	s.Require().NoError(err)
	s.markDAO = NewMarkDAO(s.db)
}

func (s *MarkTestSuite) TestGetQuestsOk() {
	rows := sqlmock.NewRows([]string{"user_id", "quest_id", "mark"}).
		AddRow(1, 1, 1.).
		AddRow(1, 2, 2.)

	s.mock.
		ExpectQuery("SELECT user_id").
		WithArgs(1).
		WillReturnRows(rows)

	marks, err := s.markDAO.GetUserMarks(1)
	s.Require().NoError(err)
	s.Equal(
		[]model.Mark{
			{UserID: 1, QuestID: 1, Mark: 1},
			{UserID: 1, QuestID: 2, Mark: 2},
		},
		marks,
	)
}

func (s *MarkTestSuite) TestGetQuestsEmpty() {
	rows := sqlmock.NewRows([]string{"user_id", "quest_id", "mark"})

	s.mock.
		ExpectQuery("SELECT user_id").
		WithArgs(1).
		WillReturnRows(rows)

	marks, err := s.markDAO.GetUserMarks(1)
	s.Require().NoError(err)
	s.Equal(
		[]model.Mark{},
		marks,
	)
}

func (s *MarkTestSuite) TestGetQuestsError() {
	s.mock.
		ExpectQuery("SELECT user_id").
		WithArgs(1).
		WillReturnError(fmt.Errorf("fail"))

	_, err := s.markDAO.GetUserMarks(1)
	s.Require().Error(err)
	s.Equal("fail", err.Error())
}

func (s *MarkTestSuite) TestMarkQuestOk() {
	var mark float32 = 3.
	userID := 1
	questID := 2

	s.mock.
		ExpectExec("UPDATE quest_user_link").
		WithArgs(mark, userID, questID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.
		ExpectExec("UPDATE quest SET rating").
		WithArgs(mark, questID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.markDAO.MarkQuest(userID, questID, mark)
	s.NoError(err)
}

func (s *MarkTestSuite) TestMarkQuestFailMark() {
	var mark float32 = 3.
	userID := 1
	questID := 2

	s.mock.
		ExpectExec("UPDATE quest_user_link").
		WithArgs(mark, userID, questID).
		WillReturnError(fmt.Errorf("fail mark"))

	s.mock.
		ExpectExec("UPDATE quest SET rating").
		WithArgs(mark, questID).
		WillReturnError(fmt.Errorf("fail rating"))

	err := s.markDAO.MarkQuest(userID, questID, mark)
	s.Require().Error(err)
	s.Equal("fail mark", err.Error())
}

func (s *MarkTestSuite) TestMarkQuestFailRating() {
	var mark float32 = 3.
	userID := 1
	questID := 2

	s.mock.
		ExpectExec("UPDATE quest_user_link").
		WithArgs(mark, userID, questID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.
		ExpectExec("UPDATE quest SET rating").
		WithArgs(mark, questID).
		WillReturnError(fmt.Errorf("fail rating"))

	err := s.markDAO.MarkQuest(userID, questID, mark)
	s.Require().Error(err)
	s.Equal("fail rating", err.Error())
}

func (s *MarkTestSuite) TestFinishOk() {
	userID := 1
	questID := 2

	s.mock.
		ExpectExec("INSERT INTO quest_user_link").
		WithArgs(userID, questID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.markDAO.FinishQuest(userID, questID)
	s.NoError(err)
}

func (s *MarkTestSuite) TestFinishError() {
	userID := 1
	questID := 2

	s.mock.
		ExpectExec("INSERT INTO quest_user_link").
		WithArgs(userID, questID).
		WillReturnError(fmt.Errorf("fail"))

	err := s.markDAO.FinishQuest(userID, questID)
	s.Require().Error(err)
	s.Equal("fail", err.Error())
}

func TestMarkTestSuite(t *testing.T) {
	suite.Run(t, new(MarkTestSuite))
}
