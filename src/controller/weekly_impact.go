package controller

import (
	"github.com/gin-gonic/gin"
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

type Period struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

func getPeriod(from, to time.Time) Period {
	return Period{
		DateFrom: from.Format("2006-01-02"),
		DateTo:   to.Format("2006-01-02"),
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
	ImpactPeriod    Period       `json:"period"`
	ImpactScore     ImpactMetric `json:"impactScore"`
	ActiveDays      ImpactMetric `json:"activeDays"`
	CommitsPerDay   ImpactMetric `json:"commitsPerDay"`
	MostChurnedFile ChurnMetric  `json:"mostChurnedFile"`
}

func getWeeklyImpact(c *gin.Context) {
	repoID := c.Param("id")
	if !validateRepoUser(c, repoID) {
		return
	}

	current := time.Now()
	lastWeek := getPeriod(getWeekRange(current.AddDate(0, 0, -7)))
	previousLastWeek := getPeriod(getWeekRange(current.AddDate(0, 0, -14)))

	weeklyData := WeeklyImpactData{
		ImpactPeriod:    lastWeek,
		ImpactScore:     getImpactScore(repoID, lastWeek, previousLastWeek),
		ActiveDays:      getActiveDays(repoID, lastWeek, previousLastWeek),
		CommitsPerDay:   getCommitsPerDay(repoID, lastWeek, previousLastWeek),
		MostChurnedFile: getMostChurnedFile(repoID, lastWeek),
	}

	c.JSON(http.StatusOK, weeklyData)
	return
}

func getImpactScore(repoID string, lastWeek, previousWeek Period) ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  184,
		PreviousPeriod: 10,
	}
}

func getActiveDays(repoID string, lastWeek, previousWeek Period) ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  10,
		PreviousPeriod: 7,
	}
}

func getCommitsPerDay(repoID string, lastWeek, previousWeek Period) ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  12.4,
		PreviousPeriod: 7.9,
	}
}

func getMostChurnedFile(repoID string, lastWeek Period) ChurnMetric {
	return ChurnMetric{
		FileName: "example-file",
		Value:    120,
	}
}
