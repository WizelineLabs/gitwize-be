package controller

import (
	"github.com/gin-gonic/gin"
	"gitwize-be/src/db"
	"gitwize-be/src/utils"
	"net/http"
	"strconv"
	"time"
)

func getStatsQuarterlyTrends(c *gin.Context) {
	defer utils.TimeTrack(time.Now(), utils.GetFuncName())
	userId := extractUserInfo(c)
	if userId == "" {
		return
	}
	idRepository := c.Param("id")

	from, err := strconv.Atoi(c.Query("date_from"))
	if err != nil {
		c.JSON(http.StatusBadRequest, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}

	to, err := strconv.Atoi(c.Query("date_to"))
	if err != nil {
		c.JSON(http.StatusBadRequest, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}

	repo := db.Repository{}
	if err := db.GetOneRepoUser(userId, idRepository, &repo); err != nil {
		c.JSON(http.StatusInternalServerError, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}

	quarterlyTrends, err := db.GetQuarterlyTrends(idRepository, int64(from), int64(to))
	if err != nil {
		c.JSON(http.StatusInternalServerError, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, quarterlyTrends)
}
