package dao

import (
	"github.com/Sovianum/acquaintanceServer/model"
)

type UserDAO interface {
	Save(user *model.User) (int, error)
	GetUserById(id int) (*model.User, error)
	GetUserByLogin(login string) (*model.User, error)
	GetNeighbourUsers(id int, distance float64) ([]*model.User, error)
	GetIdByLogin(login string) (int, error)
	ExistsById(id int) (bool, error)
	ExistsByLogin(login string) (bool, error)
}

type PositionDAO interface {
	Save(position *model.Position) error
	GetUserPosition(id int) (*model.Position, error)
}
