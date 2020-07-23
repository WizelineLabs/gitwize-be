package controller

import (
	"github.com/gin-gonic/gin"
	"gitwize-be/src/db"
	"math"
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
    "mostChurnedFile": [{
      "fileName": "abc.js",
      "value": 30
    }]
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

	now := time.Now().UTC()
	currentDuration := getWeekRange(now.AddDate(0, 0, -7))
	prevDuration := getWeekRange(now.AddDate(0, 0, -14))

	mostChurnedFiles, err := db.GetFileChurn(repoID, currentDuration.from, currentDuration.to)
	if hasErrUnknown(c, err) {
		return
	}

	currentStat, err := db.GetCommitDurationStat(repoID, currentDuration.from, currentDuration.to)
	if hasErrUnknown(c, err) {
		return
	}
	prevStat, err := db.GetCommitDurationStat(repoID, prevDuration.from, prevDuration.to)
	if hasErrUnknown(c, err) {
		return
	}

	currentModification, err := db.GetModificationStat(repoID, currentDuration.from, currentDuration.to)
	if hasErrUnknown(c, err) {
		return
	}

	prevModification, err := db.GetModificationStat(repoID, prevDuration.from, prevDuration.to)
	if hasErrUnknown(c, err) {
		return
	}

	weeklyData := WeeklyImpactData{
		ImpactPeriod:     getDatePeriod(currentDuration),
		ImpactScore:      getImpactScore(currentStat, prevStat, currentModification, prevModification),
		ActiveDays:       getActiveDays(currentStat, prevStat),
		CommitsPerDay:    getCommitsPerDay(currentStat, prevStat),
		MostChurnedFiles: mostChurnedFiles,
	}

	c.JSON(http.StatusOK, weeklyData)
	return
}

func getImpactScore(currentStat, prevStat db.DurationStat, currentModification, prevModification db.ModificationStat) ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  getImpactScoreForPeriod(currentStat, currentModification),
		PreviousPeriod: getImpactScoreForPeriod(prevStat, prevModification),
	}
}

// Impact = (5 * numFilesChanged) + (5 * numeditLocation) + (numPercentageNewcode/10) + (netChange/10)
func getImpactScoreForPeriod(durationStat db.DurationStat, modificationStat db.ModificationStat) float64 {
	numeditLocation := durationStat.Insertions
	numPercentageNewcode := 0.0
	if durationStat.Addtions != 0 {
		numPercentageNewcode = float64(modificationStat.Modifications) * 100 / float64(durationStat.Addtions)
	}
	impact := 5*float64(durationStat.NumFiles) + 5*float64(numeditLocation) + numPercentageNewcode/10 + float64(durationStat.Addtions-durationStat.Deletions)/10
	return math.Round(impact)
}

func getActiveDays(currentStat, prevStat db.DurationStat) ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  float64(currentStat.ActiveDays),
		PreviousPeriod: float64(prevStat.ActiveDays),
	}
}

func getCommitsPerDay(currentStat, prevStat db.DurationStat) ImpactMetric {
	cur, prev := 0.0, 0.0
	if currentStat.ActiveDays != 0 {
		cur = float64(currentStat.TotalCommits) / float64(currentStat.ActiveDays)
	}
	if prevStat.ActiveDays != 0 {
		prev = float64(prevStat.TotalCommits) / float64(prevStat.ActiveDays)
	}
	return ImpactMetric{
		CurrentPeriod:  cur,
		PreviousPeriod: prev,
	}
}
