package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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

type duration struct {
	from time.Time
	to   time.Time
}

func getWeeklyImpact(c *gin.Context) {
	repoID := c.Param("id")
	if !validateRepoUser(c, repoID) {
		return
	}

	current := time.Now()
	lastMonday, lastSunday := getWeekRange(current.AddDate(0, 0, -7))
	previousMonday, previousSunday := getWeekRange(current.AddDate(0, 0, -14))

	weeklyData := WeeklyImpactData{
		ImpactScore:     getImpactScore(repoID, duration{lastMonday, lastSunday}, duration{previousMonday, previousSunday}),
		ActiveDays:      getActiveDays(repoID, duration{lastMonday, lastSunday}, duration{previousMonday, previousSunday}),
		CommitsPerDay:   getCommitsPerDay(repoID, duration{lastMonday, lastSunday}, duration{previousMonday, previousSunday}),
		MostChurnedFile: getMostChurnedFile(repoID, duration{lastMonday, lastSunday}),
	}

	c.JSON(http.StatusOK, weeklyData)
	return
}

func getImpactScore(repoID string, lastWeek, previousWeek duration) ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  184,
		PreviousPeriod: 10,
	}
}

func getActiveDays(repoID string, lastWeek, previousWeek duration) ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  10,
		PreviousPeriod: 7,
	}
}

func getCommitsPerDay(repoID string, lastWeek, prevWeek duration) ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  12.4,
		PreviousPeriod: 7.9,
	}
}

func getMostChurnedFile(repoID string, lastWeek duration) ChurnMetric {
	return ChurnMetric{
		FileName: "example-file",
		Value:    120,
	}
}
