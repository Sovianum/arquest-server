package dao

import (
	"database/sql"
	"github.com/Sovianum/acquaintanceServer/model"
)

const (
	saveUser  = `INSERT INTO Users (login, password, age, sex, about) SELECT $1, $2, $3, $4, $5`
	getUser = `SELECT FROM Users (login, password, age, sex, about) WHERE id = $1`
	getNeighbourUsers = `SELECT u2.id, u2.login, u2.password, u2.age, u2.sex, u2.about
						 FROM Users u1
						 	JOIN Users u2 ON u2.id != u1.id
						 	JOIN Position p1 ON u1.id = p1.id
						 	JOIN Position p2 ON u2.id = p2.id
						 WHERE u1.id = $1 AND ST_DistanceSphere(p1.geom, p2.geom) <= $2`
	checkUser = `SELECT count(*) cnt FROM Users u WHERE u.id = $1`
)

type dbUserDAO struct {
	db *sql.DB
}

func NewDBUserDAO(db *sql.DB) UserDAO {
	var result = new(dbUserDAO)
	result.db = db
	return result
}

func (dao *dbUserDAO) Save(user *model.User) error {
	_, err := dao.db.Exec(saveUser, user.Login, user.Password, user.Age, user.Sex, user.About)
	return err
}

func (dao *dbUserDAO) Get(id int) (*model.User, error) {
	var user  = new(model.User)
	var err = dao.db.QueryRow(getUser, id).Scan(&user.Id, &user.Login, &user.Password, &user.Age, &user.Sex, &user.About)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (dao *dbUserDAO) GetNeighbour(id int, distance float64) ([]*model.User, error) {
	var rows, err = dao.db.Query(getNeighbourUsers, id, distance)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result = make([]*model.User, 0)
	for rows.Next() {
		var user = new(model.User)
		err = rows.Scan(&user.Id, &user.Login, &user.Password, &user.Age, &user.Sex, &user.About)
		if err != nil {
			return nil, err
		}

		result = append(result, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (dao *dbUserDAO) Exists(id int) (bool, error) {
	var cnt int
	var err = dao.db.QueryRow(checkUser, id).Scan(&cnt)
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}
