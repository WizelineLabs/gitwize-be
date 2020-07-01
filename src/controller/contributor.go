package controller

import (
	"github.com/gin-gonic/gin"
	"gitwize-be/src/db"
	"net/http"
	"strconv"
	"time"
)

type ContributorData struct {
	Data         []db.ContributorStats
	Contributors []db.ContributorName
}

func getContributorStats(c *gin.Context) {
	userId := extractUserInfo(c)
	repoID := c.Param("id")

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
	if err := db.GetOneRepoUser(userId, repoID, &repo); err != nil {
		c.JSON(http.StatusInternalServerError, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}

	if repo.ID == 0 {
		c.JSON(ErrCodeEntityNotFound, RestErr{
			ErrorKey:     ErrKeyEntityNotFound,
			ErrorMessage: ErrMsgEntityNotFound})
		return
	}

	data, err := db.GetChartContributorStats(repoID, getStartDateFromEpoch(from), getEndDateFromEpoch(to))
	if err != nil {
		c.JSON(http.StatusInternalServerError, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}
	contributorList, err := db.GetListContributors(repoID, getStartDateFromEpoch(from), getEndDateFromEpoch(to))
	if err != nil {
		c.JSON(http.StatusInternalServerError, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}

	contributorData := ContributorData{
		Data:         data,
		Contributors: contributorList,
	}

	c.JSON(http.StatusOK, contributorData)
	return
}

func getStartDateFromEpoch(epoch int) string {
	dateFrom := time.Unix(int64(epoch), 0)
	yearFrom, monthFrom, dayFrom := dateFrom.Year(), int(dateFrom.Month()), dateFrom.Day()
	return strconv.Itoa(yearFrom) + "-" + strconv.Itoa(monthFrom) + "-" + strconv.Itoa(dayFrom)
}

func getEndDateFromEpoch(epoch int) string {
	oneDay := 60 * 60 * 24
	return getStartDateFromEpoch(epoch + oneDay)
}
