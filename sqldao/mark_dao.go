package sqldao

import (
	"database/sql"
	"github.com/Sovianum/arquest-server/model"
	"github.com/Sovianum/arquest-server/dao"
	"net/http"
)

const (
	getUserVotes = `
		SELECT user_id, quest_id, mark FROM quest_user_link WHERE user_id = $1
	`
	markQuest = `
		UPDATE quest_user_link SET mark = $1, marked = TRUE WHERE user_id = $2 AND quest_id = $3
	`
	updateRating = `
		UPDATE quest SET (rating, mark_count) = (SELECT sum(mark) / count(*) AS rating, count(*) AS mark_count FROM quest_user_link) WHERE id = $1
	`
	finishQuest = `
		INSERT INTO quest_user_link (user_id, quest_id, completed) VALUES ($1, $2, TRUE) ON CONFLICT ON CONSTRAINT ux_user_id_quest_id DO UPDATE SET completed = TRUE 
	`
)

func NewMarkDAO(db *sql.DB) dao.MarkDAO {
	return &dbMarkDAO{db: db}
}

type dbMarkDAO struct {
	db *sql.DB
}

func (sqldao *dbMarkDAO) FinishQuest(userID, questID int) dao.DBError {
	result, err := sqldao.db.Exec(finishQuest, userID, questID)
	if err != nil {
		return dao.NewCrashDBErr(err)
	}
	return getResultErr(result)
}

func (sqldao *dbMarkDAO) MarkQuest(userID, questID int, mark float32) dao.DBError {
	if r, err := sqldao.db.Exec(markQuest, mark, userID, questID); err != nil {
		return dao.NewCrashDBErr(err)
	} else if rErr := getResultErr(r); rErr != nil {
		return rErr
	}

	if r, err := sqldao.db.Exec(updateRating, questID); err != nil {
		return dao.NewCrashDBErr(err)
	} else if rErr := getResultErr(r); rErr != nil {
		return rErr
	}
	return nil
}

func (sqldao *dbMarkDAO) GetUserMarks(userID int) ([]model.Mark, dao.DBError) {
	marks, err := sqldao.getMarks(getUserVotes, userID)
	return marks, dao.NewCrashDBErr(err)
}

func (sqldao *dbMarkDAO) getMarks(sql string, args ...interface{}) ([]model.Mark, error) {
	var rows, err = sqldao.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.Mark, 0)
	for rows.Next() {
		vote := model.Mark{}
		err = rows.Scan(
			&vote.UserID,
			&vote.QuestID,
			&vote.Mark,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, vote)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getResultErr(r sql.Result) dao.DBError {
	if affected, err := r.RowsAffected(); err != nil {
		return dao.NewCrashDBErr(err)
	} else if affected == 0 {
		return dao.NewDBErr(http.StatusNotFound, "quest not found")
	}
	return nil
}
