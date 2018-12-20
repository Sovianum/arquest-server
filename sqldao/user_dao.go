package sqldao

import (
	"database/sql"
	"github.com/Sovianum/arquest-server/model"
	"github.com/Sovianum/arquest-server/dao"
)

const (
	saveUser         = `INSERT INTO users (login, password, age, sex, about) VALUES ($1, $2, $3, $4, $5)`
	getUserById      = `SELECT id, login, password, age, sex, about FROM users WHERE id = $1`
	getUserByLogin   = `SELECT id, login, password, age, sex, about FROM users WHERE login = $1`
	getIdByLogin     = `SELECT id FROM users WHERE login = $1`
	checkUserById    = `SELECT count(*) cnt FROM users u WHERE u.id = $1`
	checkUserByLogin = `SELECT count(*) cnt FROM users u WHERE u.login = $1`
)

type dbUserDAO struct {
	db *sql.DB
}

func NewDBUserDAO(db *sql.DB) dao.UserDAO {
	result := new(dbUserDAO)
	result.db = db
	return result
}

func (sqldao *dbUserDAO) Save(user model.User) (int, dao.DBError) {
	_, saveErr := sqldao.db.Exec(saveUser, user.Login, user.Password, user.Age, user.Sex, user.About)
	if saveErr != nil {
		return 0, dao.NewCrashDBErr(saveErr)
	}

	id, getErr := sqldao.getIdByLogin(user.Login)
	if getErr != nil {
		return 0, dao.NewCrashDBErr(getErr)
	}
	// TODO add handling of the case when user saved but not extracted

	return id, nil
}

func (sqldao *dbUserDAO) GetIdByLogin(login string) (int, dao.DBError) {
	return sqldao.getIdByLogin(login)
}

func (sqldao *dbUserDAO) GetUserById(id int) (model.User, dao.DBError) {
	u := model.User{}
	err := sqldao.db.QueryRow(getUserById, id).Scan(&u.Id, &u.Login, &u.Password, &u.Age, &u.Sex, &u.About)
	if err != nil {
		return u, dao.NewCrashDBErr(err)
	}
	return u, nil
}

func (sqldao *dbUserDAO) GetUserByLogin(login string) (*model.User, dao.DBError) {
	u := new(model.User)
	err := sqldao.db.QueryRow(getUserByLogin, login).Scan(&u.Id, &u.Login, &u.Password, &u.Age, &u.Sex, &u.About)
	if err != nil {
		return nil, dao.NewCrashDBErr(err)
	}
	return u, nil
}

func (sqldao *dbUserDAO) ExistsById(id int) (bool, dao.DBError) {
	cnt := 0
	if err := sqldao.db.QueryRow(checkUserById, id).Scan(&cnt); err != nil {
		return false, dao.NewCrashDBErr(err)
	}
	return cnt > 0, nil
}

func (sqldao *dbUserDAO) ExistsByLogin(login string) (bool, dao.DBError) {
	cnt := 0
	if err := sqldao.db.QueryRow(checkUserByLogin, login).Scan(&cnt); err != nil {
		return false, dao.NewCrashDBErr(err)
	}
	return cnt > 0, nil
}

func (sqldao *dbUserDAO) getIdByLogin(login string) (int, dao.DBError) {
	id := 0
	getErr := sqldao.db.QueryRow(getIdByLogin, login).Scan(&id)
	return id, dao.NewCrashDBErr(getErr)
}
