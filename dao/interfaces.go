package dao

import (
	"github.com/Sovianum/acquaintanceServer/model"
)

type UserDAO interface {
	Save(user *model.User) error
	Get(id int) (*model.User, error)
	GetNeighbour(id int, distance float64) ([]*model.User, error)
	Exists(id int) (bool, error)
}

type PositionDAO interface {
	Save(position *model.Position) error
	GetUserPosition(id int) (*model.Position, error)
}
