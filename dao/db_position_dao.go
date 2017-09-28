package dao

import (
	"database/sql"
	"github.com/Sovianum/acquaintanceServer/model"
	"time"
)

const (
	savePosition = "INSERT INTO Position (userId, point) SELECT $1, ST_GeomFromText('POINT($2 $3)')"
	getLastPosition = "SELECT id, userId, ST_X(p.point) x, ST_Y(p.point) y, time" +
		              "FROM Position p WHERE p.userId = $1 ORDER BY time DESC LIMIT 1"
)

type dbPositionDAO struct {
	db *sql.DB
}

func NewDBPositionDAO(db *sql.DB) PositionDAO {
	var result = new(dbPositionDAO)
	result.db = db
	return result
}

func (dao *dbPositionDAO) Save(position *model.Position) error {
	_, err := dao.db.Exec(savePosition, position.UserId, position.Point.X, position.Point.Y)
	return err
}

func (dao *dbPositionDAO) GetUserPosition(id int) (*model.Position, error) {
	var position = new(model.Position)
	var posTime time.Time
	var err = dao.db.QueryRow(getLastPosition, id).Scan(
		&position.Id, &position.UserId, &position.Point.X, &position.Point.Y, &posTime,
	)
	if err != nil {
		return nil, err
	}

	position.Time = model.QuotedTime(posTime)
	return position, nil
}
