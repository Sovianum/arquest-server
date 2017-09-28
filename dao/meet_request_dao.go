package dao

import (
	"github.com/Sovianum/acquaintanceServer/model"
	"database/sql"
	"net/http"
)

const (
	countPendingRequests = `
		SELECT count(*) FROM MeetRequest
		WHERE requesterId = $1 AND requestedId = $2 AND status = 'PENDING' AND age(now(), time) < '$3 minutes'
	`
	getPendingRequests = `
		SELECT id, requesterId, requestedId, status, time FROM MeetRequest mr
			JOIN Users u ON mr.requesterId = u.id
			JOIN Position p ON p.userId = u.id
		WHERE mr.requestedId = $1 AND age(now(), p.time) < $2 AND status = 'PENDING'
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
)

type MeetRequestDAO interface {
	CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error)
	GetRequests(requestedId int, onlineTimeoutMin int) ([]*model.MeetRequest, error)
}

type meetRequestDAO struct {
	db *sql.DB
}

func NewMeetDAO(db *sql.DB) MeetRequestDAO {
	return &meetRequestDAO{
		db: db,
	}
}

func (dao *meetRequestDAO) GetRequests(requestedId int, onlineTimeoutMin int) ([]*model.MeetRequest, error) {
	return dao.getPendingRequests(requestedId, onlineTimeoutMin)
}

func (dao *meetRequestDAO) CreateRequest(requesterId int, requestedId int, requestTimeoutMin int, maxDistance float64) (code int, dbErr error) {
	var requestCnt, countErr = dao.countPendingRequests(requesterId, requestedId, requestTimeoutMin)
	if countErr != nil {
		return http.StatusInternalServerError, countErr
	}

	if requestCnt > 0 {
		return http.StatusForbidden, nil
	}

	var accessible, accessErr = dao.isAccessible(requesterId, requestedId, maxDistance, requestTimeoutMin)
	if accessErr != nil {
		return http.StatusInternalServerError, accessErr
	}

	if !accessible {
		return http.StatusForbidden, nil
	}

	countErr = dao.createRequest(requesterId, requestedId)
	if countErr != nil {
		return http.StatusInternalServerError, countErr
	}

	return http.StatusOK, nil
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

func (dao *meetRequestDAO) getPendingRequests(requestedId int, onlineTimeoutMin int) ([]*model.MeetRequest, error) {
	var rows, err = dao.db.Query(getPendingRequests, requestedId, onlineTimeoutMin)
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

func (dao *meetRequestDAO) countPendingRequests(requesterId int, requestedId int, requestTimeoutMin int) (int, error) {
	var cnt int
	var err = dao.db.QueryRow(countPendingRequests, requesterId, requestedId, requestTimeoutMin).Scan(&cnt)
	return cnt, err
}
