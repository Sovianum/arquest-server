package dao

import (
	"github.com/Sovianum/acquaintanceServer/model"
	"database/sql"
)

const (
	countPendingRequests = `
		SELECT count(*) FROM MeetRequest
		WHERE requesterId = $1 AND requestedId = $2 AND status = 'PENDING'
	`
	getPendingRequests = `
		SELECT id, requesterId, requestedId, status, time FROM MeetRequest mr
			JOIN Users u ON mr.requesterId = u.id
			JOIN Position p ON p.userId = u.id
		WHERE mr.requestedId = $1 AND status = 'PENDING'
	`
	checkAccessibility = `
		SELECT ST_DistanceSphere(p1.geom, p2.geom) < $1 FROM
			(
				SELECT * FROM Position
				WHERE userId = $2 AND age(now(), time) < '$4 minutes'
				ORDER BY time DESC
				LIMIT 1
			) p1,
			(
				SELECT * FROM Position
				WHERE userId = $3 AND age(now(), time) < '$4 minutes'
				ORDER BY time DESC
				LIMIT 1
			) p2
	`
	createRequest = `
		INSERT INTO MeetRequest (requesterId, requestedId) SELECT $1, $2
	`
	updateRequestStatus = `
		UPDATE MeetRequest SET status = $1 WHERE id = $2 AND requestedId = $3
	`
)

type MeetRequestDAO interface {
	CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error)
	GetRequests(requestedId int) ([]*model.MeetRequest, error)
	AcceptRequest(id int, requestedId int) (int, error)
	DeclineRequest(id int, requestedId int) (int, error)
}

type meetRequestDAO struct {
	db *sql.DB
}

func NewMeetDAO(db *sql.DB) MeetRequestDAO {
	return &meetRequestDAO{
		db: db,
	}
}

func (dao *meetRequestDAO) DeclineRequest(id int, requestedId int) (int, error) {
	return dao.updateRequestStatus(id, requestedId, model.StatusDeclined)
}

func (dao *meetRequestDAO) AcceptRequest(id int, requestedId int) (int, error) {
	return dao.updateRequestStatus(id, requestedId, model.StatusAccepted)
}

func (dao *meetRequestDAO) GetRequests(requestedId int) ([]*model.MeetRequest, error) {
	return dao.getPendingRequests(requestedId)
}

func (dao *meetRequestDAO) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (int, error) {
	var requestCnt, countErr = dao.countPendingRequests(requesterId, requestedId)
	if countErr != nil {
		return 0, countErr
	}

	if requestCnt > 0 {
		return 0, nil
	}

	var accessible, accessErr = dao.isAccessible(requesterId, requestedId, maxDistance, requestTimeoutMin)
	if accessErr != nil {
		return 0, accessErr
	}

	if !accessible {
		return 0, nil
	}

	var result, createErr = dao.db.Exec(createRequest, requesterId, requestedId)
	if createErr != nil {
		return 0, createErr
	}

	var rowsAffected, rowsErr = result.RowsAffected()
	return int(rowsAffected), rowsErr
}

func (dao *meetRequestDAO) updateRequestStatus(id int, requestedId int, status string) (int, error) {
	var result, err = dao.db.Exec(updateRequestStatus, status, id, requestedId)
	if err != nil {
		return 0, err
	}

	var rowsAffected, rowsErr = result.RowsAffected()
	if rowsErr != nil {
		return 0, rowsErr
	}

	return int(rowsAffected), nil
}

func (dao *meetRequestDAO) isAccessible(id1 int, id2 int, maxDistance float64, timeoutMin int) (bool, error) {
	var rows, err = dao.db.Query(checkAccessibility, maxDistance, id1, id2, timeoutMin)
	if err != nil {
		return false, err
	}

	var accessible = false
	for rows.Next() {
		err = rows.Scan(&accessible)
		if err != nil {
			return false, err
		}
	}
	err = rows.Err()
	if err != nil {
		return false, err
	}

	return accessible, nil
}

func (dao *meetRequestDAO) createRequest(requesterId int, requestedId int) error {
	var _, err = dao.db.Exec(createRequest, requesterId, requestedId)
	return err
}

func (dao *meetRequestDAO) getPendingRequests(requestedId int) ([]*model.MeetRequest, error) {
	var rows, err = dao.db.Query(getPendingRequests, requestedId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result = make([]*model.MeetRequest, 0)
	for rows.Next() {
		var request = new(model.MeetRequest)
		err = rows.Scan(&request.Id, &request.RequesterId, &request.RequestedId, &request.Status, &request.Time)
		if err != nil {
			return nil, err
		}
		result = append(result, request)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (dao *meetRequestDAO) countPendingRequests(requesterId int, requestedId int) (int, error) {
	var cnt int
	var err = dao.db.QueryRow(countPendingRequests, requesterId, requestedId).Scan(&cnt)
	return cnt, err
}
