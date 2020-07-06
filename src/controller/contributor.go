package controller

import (
	"github.com/gin-gonic/gin"
	"gitwize-be/src/db"
	"net/http"
	"strconv"
	"time"
)

type ContributorAPIData struct {
	Table        []ContributorTableItem           `json:"table"`
	Chart        map[string][]db.ContributorStats `json:"chart"`
	Contributors []db.Contributor                 `json:"contributors"`
}

const AVERAGE = "average"

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

	dataByPerson, err := db.GetContributorStatsByPerson(repoID, getStartDateFromEpoch(from), getEndDateFromEpoch(to))
	if err != nil {
		c.JSON(http.StatusInternalServerError, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}

	dataTotal, err := db.GetTotalContributorStats(repoID, getStartDateFromEpoch(from), getEndDateFromEpoch(to))
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

	chartData := getChartData(dataByPerson, dataTotal, len(contributorList))
	contributorData := ContributorAPIData{
		Table:        getTableData(chartData),
		Chart:        chartData,
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

func getChartData(dataPerson []db.ContributorStats, dataTotal []db.ContributorStats, numbOfContributor int) map[string][]db.ContributorStats {
	chartData := map[string][]db.ContributorStats{}
	dataByDayMap := getStatByDayMap(dataTotal)

	for _, stat := range dataPerson {
		dataByDay := dataByDayMap[stat.Date]
		if dataByDay.AdditionLoc+dataByDay.DeletionLoc == 0 {
			stat.LOCPercent = 0
		} else {
			stat.LOCPercent = float32(stat.AdditionLoc+stat.DeletionLoc) / float32(dataByDay.AdditionLoc+dataByDay.DeletionLoc) * 100
		}
		items := chartData[stat.Email]
		items = append(items, stat)
		chartData[stat.Email] = items
	}
	chartData[AVERAGE] = getAverageStatByDay(dataTotal, numbOfContributor)
	return chartData
}

func getStatByDayMap(data []db.ContributorStats) map[string]db.ContributorStats {
	result := make(map[string]db.ContributorStats, len(data))
	for _, stat := range data {
		result[stat.Date] = stat
	}
	return result
}

func getAverageStatByDay(data []db.ContributorStats, numbOfContributor int) []db.ContributorStats {
	result := make([]db.ContributorStats, len(data))
	if numbOfContributor != 0 {
		for i, item := range data {
			result[i] = db.ContributorStats{
				RepositoryID: item.RepositoryID,
				Commits:      item.Commits / numbOfContributor,
				AdditionLoc:  item.AdditionLoc / numbOfContributor,
				DeletionLoc:  item.DeletionLoc / numbOfContributor,
				NumFiles:     item.NumFiles / numbOfContributor,
				LOCPercent:   100 / float32(numbOfContributor),
				Date:         item.Date,
			}
		}
	}
	return result
}

type ContributorTableItem struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Commits     int    `json:"commits"`
	Additions   int    `json:"additions"`
	Deletions   int    `json:"deletions"`
	ActiveDays  int    `json:"activeDays"`
	FilesChange int    `json:"filesChange"`
}

func getTableData(dataMap map[string][]db.ContributorStats) []ContributorTableItem {
	result := []ContributorTableItem{}
	for email, contributorData := range dataMap {
		if email == AVERAGE {
			continue
		} else {
			tableItem := buildTableItem(contributorData)
			result = append(result, tableItem)
		}
	}
	return result
}

func buildTableItem(data []db.ContributorStats) ContributorTableItem {
	item := ContributorTableItem{}
	if len(data) == 0 {
		return item
	}
	item.Name = data[0].Name
	item.Email = data[0].Email
	commits, additions, deletions, activeDays, fileChanges := 0, 0, 0, 0, 0
	for _, dataItem := range data {
		activeDays++
		commits += dataItem.Commits
		additions += dataItem.AdditionLoc
		deletions += dataItem.DeletionLoc
		fileChanges += dataItem.NumFiles
	}
	item.Commits = commits
	item.Additions = additions
	item.Deletions = deletions
	item.ActiveDays = activeDays
	item.FilesChange = fileChanges
	return item
}
