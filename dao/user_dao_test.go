package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Sovianum/arquest-server/model"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

const (
	password = "password"
)

type UserTestSuite struct {
	suite.Suite
	db       *sql.DB
	mock     sqlmock.Sqlmock
	userDAO  UserDAO
	passHash []byte
}

func (s *UserTestSuite) SetupTest() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	s.Require().NoError(err)
	s.userDAO = NewDBUserDAO(s.db)
	s.passHash, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (s *UserTestSuite) TestExistsByIDFound() {
	rows := sqlmock.NewRows([]string{"cnt"}).
		AddRow(1)

	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(10).
		WillReturnRows(rows)

	exists, dbErr := s.userDAO.ExistsById(10)

	s.NoError(dbErr)
	s.True(exists)
}

func (s *UserTestSuite) TestExistsByIDNoFound() {
	rows := sqlmock.NewRows([]string{"cnt"}).
		AddRow(0)

	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(10).
		WillReturnRows(rows)

	exists, dbErr := s.userDAO.ExistsById(10)

	s.NoError(dbErr)
	s.False(exists)
}

func (s *UserTestSuite) TestExistsByIDError() {
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs(10).
		WillReturnError(fmt.Errorf("failed to check"))

	_, dbErr := s.userDAO.ExistsById(10)

	s.Error(dbErr)
	s.Equal("failed to check", dbErr.Error())
}

func (s *UserTestSuite) TestExistsByLoginFound() {
	rows := sqlmock.NewRows([]string{"cnt"}).
		AddRow(1)

	s.mock.
		ExpectQuery("SELECT count").
		WithArgs("login").
		WillReturnRows(rows)

	exists, dbErr := s.userDAO.ExistsByLogin("login")

	s.NoError(dbErr)
	s.True(exists)
}

func (s *UserTestSuite) TestExistsByLoginNotFound() {
	rows := sqlmock.NewRows([]string{"cnt"}).
		AddRow(0)

	s.mock.
		ExpectQuery("SELECT count").
		WithArgs("login").
		WillReturnRows(rows)

	exists, dbErr := s.userDAO.ExistsByLogin("login")

	s.NoError(dbErr)
	s.False(exists)
}

func (s *UserTestSuite) TestExistsByLoginError() {
	s.mock.
		ExpectQuery("SELECT count").
		WithArgs("login").
		WillReturnError(fmt.Errorf("fail"))

	_, dbErr := s.userDAO.ExistsByLogin("login")

	s.Error(dbErr)
	s.Equal("fail", dbErr.Error())
}

func (s *UserTestSuite) TestSaveSuccess() {
	s.mock.
		ExpectExec("INSERT INTO").
		WithArgs("login", string(s.passHash), 100, model.FEMALE, "").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.
		ExpectQuery("SELECT").
		WithArgs("login").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	user := &model.User{Login: "login", Password: string(s.passHash), Sex: model.FEMALE, Age: 100}

	userDAO := NewDBUserDAO(s.db)
	id, saveErr := userDAO.Save(user)

	s.NoError(saveErr)
	s.Equal(1, id)
}

func (s *UserTestSuite) TestSaveDuplicateLogin() {
	s.mock.
		ExpectExec("INSERT INTO").
		WithArgs("login", string(s.passHash), 100, model.FEMALE, "").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.
		ExpectQuery("SELECT").
		WithArgs("login").
		WillReturnError(fmt.Errorf("duplicate id"))

	user := &model.User{Login: "login", Password: string(s.passHash), Sex: model.FEMALE, Age: 100}

	userDAO := NewDBUserDAO(s.db)
	_, saveErr := userDAO.Save(user)

	s.Error(saveErr)
	s.Equal("duplicate id", saveErr.Error())
}

func (s *UserTestSuite) TestGetUserByIdSuccess() {
	rows := sqlmock.NewRows([]string{"id", "login", "password", "age", "sex", "about"}).
		AddRow(1, "login", "pass", 100, model.MALE, "about")

	s.mock.
		ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnRows(rows)

	user := &model.User{Id: 1, Login: "login", Password: "pass", Sex: model.MALE, Age: 100, About: "about"}
	dbUser, userErr := s.userDAO.GetUserById(1)

	s.NoError(userErr)
	s.Equal(user, dbUser)
}

func (s *UserTestSuite) TestGetUserByIdNotFound() {
	s.mock.
		ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnError(errors.New("user not found"))

	_, userErr := s.userDAO.GetUserById(1)

	s.Error(userErr)
	s.Equal("user not found", userErr.Error())
}

func (s *UserTestSuite) TestGetUserByLoginSuccess() {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	s.mock.
		ExpectQuery("SELECT").
		WithArgs("login").
		WillReturnRows(rows)

	id, err := s.userDAO.GetIdByLogin("login")

	s.NoError(err)
	s.Equal(1, id)
}

func (s *UserTestSuite) TestGetUserByLoginNotFound() {
	s.mock.
		ExpectQuery("SELECT").
		WithArgs("login").
		WillReturnError(errors.New("user not found"))

	_, err := s.userDAO.GetIdByLogin("login")

	s.Error(err)
	s.Equal("user not found", err.Error())
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
