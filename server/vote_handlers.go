package server

import (
	"github.com/Sovianum/arquest-server/common"
	"github.com/Sovianum/arquest-server/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (env *Env) FinishQuest(c *gin.Context) {
	env.updateLinkTable(c, func(vote model.Vote) error {
		return env.voteDAO.FinishQuest(vote.UserID, vote.QuestID)
	})
}

func (env *Env) MarkQuest(c *gin.Context) {
	env.updateLinkTable(c, func(vote model.Vote) error {
		return env.voteDAO.MarkQuest(vote.UserID, vote.QuestID, vote.Mark)
	})
}

func (env *Env) GetUserVotes(c *gin.Context) {
	id := c.GetInt(UserID)
	votes, err := env.voteDAO.GetUserVotes(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrResponse(err))
		return
	}
	c.JSON(http.StatusOK, votes)
}

func (env *Env) updateLinkTable(c *gin.Context, updateFunc func(vote model.Vote) error) {
	id := c.GetInt(UserID)
	var vote model.Vote
	if err := c.ShouldBindJSON(&vote); err != nil {
		c.JSON(http.StatusBadRequest, common.GetErrResponse(err))
		return
	}
	vote.UserID = id

	if err := updateFunc(vote); err != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrResponse(err))
		return
	}
	c.JSON(http.StatusOK, common.GetEmptyResponse())
}
