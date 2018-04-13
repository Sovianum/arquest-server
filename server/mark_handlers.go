package server

import (
	"fmt"
	"github.com/Sovianum/arquest-server/common"
	"github.com/Sovianum/arquest-server/dao"
	"github.com/Sovianum/arquest-server/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (env *Env) FinishQuest(c *gin.Context) {
	env.updateLinkTable(c, func(vote model.Mark) dao.DBError {
		return env.markDAO.FinishQuest(vote.UserID, vote.QuestID)
	})
}

func (env *Env) MarkQuest(c *gin.Context) {
	env.updateLinkTable(c, func(mark model.Mark) dao.DBError {
		return env.markDAO.MarkQuest(mark.UserID, mark.QuestID, mark.Mark)
	})
}

func (env *Env) GetUserVotes(c *gin.Context) {
	id := c.GetInt(UserID)
	votes, err := env.markDAO.GetUserMarks(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrResponse(err))
		return
	}
	c.JSON(http.StatusOK, votes)
}

func (env *Env) updateLinkTable(c *gin.Context, updateFunc func(vote model.Mark) dao.DBError) {
	id := c.GetInt(UserID)
	var vote model.Mark
	if err := c.ShouldBindJSON(&vote); err != nil {
		c.JSON(http.StatusBadRequest, common.GetErrResponse(err))
		return
	}
	if vote.UserID == 0 { // default value
		vote.UserID = id
	}

	if exists, err := env.questDAO.ExistsByID(vote.QuestID); err != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrResponse(err))
		return
	} else if !exists {
		c.JSON(http.StatusNotFound, common.GetErrResponse(fmt.Errorf("quest not found")))
		return
	}

	if id != vote.UserID {
		c.JSON(http.StatusForbidden, common.GetErrResponse(fmt.Errorf("you can not vote as another person")))
		return
	}

	if err := updateFunc(vote); err != nil {
		c.JSON(err.Code(), common.GetErrResponse(err))
		return
	}
	c.JSON(http.StatusOK, common.GetEmptyResponse())
}
