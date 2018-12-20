package dao

import (
	"github.com/Sovianum/arquest-server/model"
)

type UserDAO interface {
	Save(user model.User) (int, DBError)
	GetUserById(id int) (model.User, DBError)
	GetUserByLogin(login string) (*model.User, DBError)
	GetIdByLogin(login string) (int, DBError)
	ExistsById(id int) (bool, DBError)
	ExistsByLogin(login string) (bool, DBError)
}
