package dao

import (
	"github.com/Sovianum/myTgtTest/model"
	"time"
)

type UserDAO interface {
	Save(r model.Registration) error
	Exists(id uint) (bool, error)
}

type StatsDAO interface {
	Save(s model.Stats) error
	GetStatsSlice(dates []time.Time, action string, limit int) (model.StatsSlice, error)
}
