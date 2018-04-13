package dao

import (
	"database/sql"
	"github.com/Sovianum/arquest-server/model"
)

const (
	saveUser         = `INSERT INTO users (login, password, age, sex, about) VALUES ($1, $2, $3, $4, $5)`
	getUserById      = `SELECT id, login, password, age, sex, about FROM users WHERE id = $1`
	getUserByLogin   = `SELECT id, login, password, age, sex, about FROM users WHERE login = $1`
	getIdByLogin     = `SELECT id FROM users WHERE login = $1`
	checkUserById    = `SELECT count(*) cnt FROM users u WHERE u.id = $1`
	checkUserByLogin = `SELECT count(*) cnt FROM users u WHERE u.login = $1`
)

type UserDAO interface {
	Save(user model.User) (int, DBError)
	GetUserById(id int) (model.User, DBError)
	GetUserByLogin(login string) (*model.User, DBError)
	GetIdByLogin(login string) (int, DBError)
	ExistsById(id int) (bool, DBError)
	ExistsByLogin(login string) (bool, DBError)
}

type dbUserDAO struct {
	db *sql.DB
}

func NewDBUserDAO(db *sql.DB) UserDAO {
	result := new(dbUserDAO)
	result.db = db
	return result
}

func (dao *dbUserDAO) Save(user model.User) (int, DBError) {
	_, saveErr := dao.db.Exec(saveUser, user.Login, user.Password, user.Age, user.Sex, user.About)
	if saveErr != nil {
		return 0, NewCrashDBErr(saveErr)
	}

	id, getErr := dao.getIdByLogin(user.Login)
	if getErr != nil {
		return 0, NewCrashDBErr(getErr)
	}
	// TODO add handling of the case when user saved but not extracted

	return id, nil
}

func (dao *dbUserDAO) GetIdByLogin(login string) (int, DBError) {
	return dao.getIdByLogin(login)
}

func (dao *dbUserDAO) GetUserById(id int) (model.User, DBError) {
	u := model.User{}
	err := dao.db.QueryRow(getUserById, id).Scan(&u.Id, &u.Login, &u.Password, &u.Age, &u.Sex, &u.About)
	if err != nil {
		return u, NewCrashDBErr(err)
	}
	return u, nil
}

func (dao *dbUserDAO) GetUserByLogin(login string) (*model.User, DBError) {
	u := new(model.User)
	err := dao.db.QueryRow(getUserByLogin, login).Scan(&u.Id, &u.Login, &u.Password, &u.Age, &u.Sex, &u.About)
	if err != nil {
		return nil, NewCrashDBErr(err)
	}
	return u, nil
}

func (dao *dbUserDAO) ExistsById(id int) (bool, DBError) {
	cnt := 0
	if err := dao.db.QueryRow(checkUserById, id).Scan(&cnt); err != nil {
		return false, NewCrashDBErr(err)
	}
	return cnt > 0, nil
}

func (dao *dbUserDAO) ExistsByLogin(login string) (bool, DBError) {
	cnt := 0
	if err := dao.db.QueryRow(checkUserByLogin, login).Scan(&cnt); err != nil {
		return false, NewCrashDBErr(err)
	}
	return cnt > 0, nil
}

func (dao *dbUserDAO) getIdByLogin(login string) (int, DBError) {
	id := 0
	getErr := dao.db.QueryRow(getIdByLogin, login).Scan(&id)
	return id, NewCrashDBErr(getErr)
}
