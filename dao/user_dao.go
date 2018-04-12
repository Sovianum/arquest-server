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
	Save(user model.User) (int, error)
	GetUserById(id int) (model.User, error)
	GetUserByLogin(login string) (*model.User, error)
	GetIdByLogin(login string) (int, error)
	ExistsById(id int) (bool, error)
	ExistsByLogin(login string) (bool, error)
}

type dbUserDAO struct {
	db *sql.DB
}

func NewDBUserDAO(db *sql.DB) UserDAO {
	result := new(dbUserDAO)
	result.db = db
	return result
}

func (dao *dbUserDAO) Save(user model.User) (int, error) {
	_, saveErr := dao.db.Exec(saveUser, user.Login, user.Password, user.Age, user.Sex, user.About)
	if saveErr != nil {
		return 0, saveErr
	}

	id, getErr := dao.getIdByLogin(user.Login)
	if getErr != nil {
		return 0, getErr
	}
	// TODO add handling of the case when user saved but not extracted

	return id, nil
}

func (dao *dbUserDAO) GetIdByLogin(login string) (int, error) {
	return dao.getIdByLogin(login)
}

func (dao *dbUserDAO) GetUserById(id int) (model.User, error) {
	u := model.User{}
	err := dao.db.QueryRow(getUserById, id).Scan(&u.Id, &u.Login, &u.Password, &u.Age, &u.Sex, &u.About)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (dao *dbUserDAO) GetUserByLogin(login string) (*model.User, error) {
	u := new(model.User)
	err := dao.db.QueryRow(getUserByLogin, login).Scan(&u.Id, &u.Login, &u.Password, &u.Age, &u.Sex, &u.About)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (dao *dbUserDAO) ExistsById(id int) (bool, error) {
	cnt := 0
	if err := dao.db.QueryRow(checkUserById, id).Scan(&cnt); err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (dao *dbUserDAO) ExistsByLogin(login string) (bool, error) {
	cnt := 0
	if err := dao.db.QueryRow(checkUserByLogin, login).Scan(&cnt); err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (dao *dbUserDAO) getIdByLogin(login string) (int, error) {
	id := 0
	getErr := dao.db.QueryRow(getIdByLogin, login).Scan(&id)
	return id, getErr
}
