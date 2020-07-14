package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
spec https://wizeline.atlassian.net/wiki/spaces/GWZ/pages/1424393330/API+spec+-+Weekly+Impact
{
    "impactScore": {
      "currentPeriod": 184,
      "previousPeriod": 10,
    },
    "activeDays": {
      "currentPeriod": 15,
      "previousPeriod": 12,
    },
    "commitsPerDay": {
      "currentPeriod": 13.5,
      "previousPeriod": 11.2,
    },
    "mostChurnedFile": {
      "fileName": "abc.js"
      "value": 30
    }
}
*/

type ImpactMetric struct {
	CurrentPeriod  float32 `json:"currentPeriod"`
	PreviousPeriod float32 `json:"previousPeriod"`
}

type ChurnMetric struct {
	FileName string `json:"fileName"`
	Value    int    `json:"value"`
}

type WeeklyImpactData struct {
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

	values, err := getIntParams(c, "date_from", "date_to")
	if err != nil {
		return
	}
	from, to := values[0], values[1]

	weeklyData := WeeklyImpactData{
		ImpactScore:     getImpactScore(repoID, from, to),
		ActiveDays:      getActiveDays(repoID, from, to),
		CommitsPerDay:   getCommitsPerDay(repoID, from, to),
		MostChurnedFile: getMostChurnedFile(repoID, from, to),
	}

	c.JSON(http.StatusOK, weeklyData)
	return
}

func getImpactScore(repoID string, from, to int) ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  184,
		PreviousPeriod: 10,
	}
}

func getActiveDays(repoID string, from, to int) ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  10,
		PreviousPeriod: 7,
	}
}

func getCommitsPerDay(repoID string, from, to int) ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  12.4,
		PreviousPeriod: 7.9,
	}
}

func getMostChurnedFile(repoID string, from, to int) ChurnMetric {
	return ChurnMetric{
		FileName: "example-file",
		Value:    120,
	}
}
