package dao

import (
	"errors"
	"github.com/Sovianum/myTgtTest/model"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestDbUserDAO_Exists_ClientFound(t *testing.T) {
	var db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var rows = sqlmock.NewRows([]string{"cnt"}).
		AddRow(1)

	mock.
		ExpectQuery("SELECT count").
		WithArgs(10).
		WillReturnRows(rows)

	var userDAO = NewDBUserDAO(db)
	var exists, dbErr = userDAO.Exists(10)

	if dbErr != nil {
		t.Error(dbErr)
	}

	if !exists {
		t.Error("Failed to find existing user")
	}
}

func TestDbUserDAO_Exists_ClientNotFound(t *testing.T) {
	var db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var rows = sqlmock.NewRows([]string{"cnt"}).
		AddRow(0)

	mock.
		ExpectQuery("SELECT count").
		WithArgs(10).
		WillReturnRows(rows)

	var userDAO = NewDBUserDAO(db)
	var exists, dbErr = userDAO.Exists(10)

	if dbErr != nil {
		t.Error(dbErr)
	}

	if exists {
		t.Error("Succeeded to find non existing user")
	}
}

func TestDbUserDAO_Exists_DBFailed(t *testing.T) {
	var db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.
		ExpectQuery("SELECT count").
		WithArgs(10).
		WillReturnError(errors.New("Failed to check"))

	var userDAO = NewDBUserDAO(db)
	var _, dbErr = userDAO.Exists(10)

	if dbErr == nil {
		t.Error("Had to crash")
	}

	if dbErr.Error() != "Failed to check" {
		t.Errorf("Wrong error expected %v got %v", "\"Failed to check\"", dbErr.Error())
	}
}

func TestDbUserDAO_Save_Success(t *testing.T) {
	var db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.
		ExpectExec("INSERT INTO").
		WithArgs(1, 1, "F").
		WillReturnResult(sqlmock.NewResult(1, 1))

	var r = model.Registration{Id: 1, Age: 1, Sex: model.FEMALE}

	var userDAO = NewDBUserDAO(db)
	var saveErr = userDAO.Save(r)

	if saveErr != nil {
		t.Error(saveErr.Error())
	}
}

func TestDbUserDAO_Save_DuplicateId(t *testing.T) {
	var db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.
		ExpectExec("INSERT INTO").
		WithArgs(1, 1, model.FEMALE).
		WillReturnError(errors.New("Duplicate id"))

	var r = model.Registration{Id: 1, Age: 1, Sex: model.FEMALE}

	var userDAO = NewDBUserDAO(db)
	var saveErr = userDAO.Save(r)

	if saveErr == nil {
		t.Error("Had to crash")
	}

	if saveErr.Error() != "Duplicate id" {
		t.Errorf("Wrong error expected %v got %v", "\"Duplicate id\"", saveErr.Error())
	}
}
