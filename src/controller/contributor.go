package controller

import (
	"github.com/gin-gonic/gin"
	"gitwize-be/src/db"
	"net/http"
)

//ContributorAPIData data struct for contributor api response
type ContributorAPIData struct {
	Table        []ContributorTableItem           `json:"table"`
	Chart        map[string][]db.ContributorStats `json:"chart"`
	Contributors []db.Contributor                 `json:"contributors"`
}

const average = "average"

func getContributorStats(c *gin.Context) {
	repoID := c.Param("id")

	if !validateRepoUser(c, repoID) {
		return
	}

	values, err := getIntParams(c, "date_from", "date_to")
	if err != nil {
		return
	}
	from, to := values[0], values[1]

	dataByPerson, err := db.GetContributorStatsByPerson(repoID, getStartDateFromEpoch(from), getEndDateFromEpoch(to))
	if hasErrUnknown(c, err) {
		return
	}

	dataTotal, err := db.GetTotalContributorStats(repoID, getStartDateFromEpoch(from), getEndDateFromEpoch(to))
	if hasErrUnknown(c, err) {
		return
	}

	contributorList, err := db.GetListContributors(repoID, getStartDateFromEpoch(from), getEndDateFromEpoch(to))
	if hasErrUnknown(c, err) {
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
	chartData[average] = getAverageStatByDay(dataTotal, numbOfContributor)
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
	NetChanges  int    `json:"netChanges"`
	ActiveDays  int    `json:"activeDays"`
	FilesChange int    `json:"filesChange"`
}

func getTableData(dataMap map[string][]db.ContributorStats) []ContributorTableItem {
	result := []ContributorTableItem{}
	for email, contributorData := range dataMap {
		if email == average {
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
	item.NetChanges = additions - deletions
	item.ActiveDays = activeDays
	item.FilesChange = fileChanges
	return item
}
