package server

import (
	"fmt"
	"github.com/Sovianum/arquest-server/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (env *Env) GetAllQuests(c *gin.Context) {
	quests, err := env.questDAO.GetAllQuests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrResponse(err))
		return
	}
	for i := range quests {
		quests[i].DataPath = getQuestDataUrl(env.conf.Logic.QuestDataTemplate, quests[i].ID)
	}
	c.JSON(http.StatusOK, common.GetDataResponse(quests))
}

func (env *Env) GetFinishedQuests(c *gin.Context) {
	id := c.GetInt(UserID)
	quests, err := env.questDAO.GetFinishedQuests(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.GetErrResponse(err))
		return
	}
	c.JSON(http.StatusOK, common.GetDataResponse(quests))
}

func getQuestDataUrl(template string, questID int) string {
	return fmt.Sprintf(template, questID)
}
