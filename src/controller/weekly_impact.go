package controller

import (
	"github.com/gin-gonic/gin"
	"gitwize-be/src/db"
	"net/http"
	"time"
)

/*
spec https://wizeline.atlassian.net/wiki/spaces/GWZ/pages/1424393330/API+spec+-+Weekly+Impact
{
    "period": {
      date_from: "2020-06-29",
      date_to: "2020-07-05"
    }
    "impactScore": {
      "currentPeriod": 184,
      "previousPeriod": 10
    },
    "activeDays": {
      "currentPeriod": 15,
      "previousPeriod": 12
    },
    "commitsPerDay": {
      "currentPeriod": 13.5,
      "previousPeriod": 11.2
    },
    "mostChurnedFile": {
      "fileName": "abc.js",
      "value": 30
    }
}

*/

type DatePeriod struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

func getDatePeriod(r TimeRange) DatePeriod {
	return DatePeriod{
		DateFrom: r.from.Format("2006-01-02"),
		DateTo:   r.to.Format("2006-01-02"),
	}
}

type ImpactMetric struct {
	CurrentPeriod  float64 `json:"currentPeriod"`
	PreviousPeriod float64 `json:"previousPeriod"`
}

type ChurnMetric struct {
	FileName string `json:"fileName"`
	Value    int    `json:"value"`
}

type WeeklyImpactData struct {
	ImpactPeriod     DatePeriod     `json:"period"`
	ImpactScore      ImpactMetric   `json:"impactScore"`
	ActiveDays       ImpactMetric   `json:"activeDays"`
	CommitsPerDay    ImpactMetric   `json:"commitsPerDay"`
	MostChurnedFiles []db.FileChurn `json:"mostChurnedFiles"`
}

func getWeeklyImpact(c *gin.Context) {
	repoID := c.Param("id")
	if !validateRepoUser(c, repoID) {
		return
	}

	current := time.Now().UTC()
	lastWeekRange := getWeekRange(current.AddDate(0, 0, -7))
	previousLastWeek := getWeekRange(current.AddDate(0, 0, -14))

	mostChurnedFiles, err := db.GetFileChurn(repoID, lastWeekRange.from, lastWeekRange.to)
	if hasErrUnknown(c, err) {
		return
	}

	weeklyData := WeeklyImpactData{
		ImpactPeriod:     getDatePeriod(lastWeekRange),
		ImpactScore:      getImpactScore(repoID, lastWeekRange, previousLastWeek),
		ActiveDays:       getActiveDays(repoID, lastWeekRange, previousLastWeek),
		CommitsPerDay:    getCommitsPerDay(repoID, lastWeekRange, previousLastWeek),
		MostChurnedFiles: mostChurnedFiles,
	}

	c.JSON(http.StatusOK, weeklyData)
	return
}

func getImpactScore(repoID string, lastWeek, previousWeek TimeRange) ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  184,
		PreviousPeriod: 10,
	}
}

func getActiveDays(repoID string, lastWeek, previousWeek TimeRange) ImpactMetric {

	return ImpactMetric{
		CurrentPeriod:  10,
		PreviousPeriod: 7,
	}
}

func getCommitsPerDay(repoID string, lastWeek, previousWeek TimeRange) ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  12.4,
		PreviousPeriod: 7.9,
	}
}
