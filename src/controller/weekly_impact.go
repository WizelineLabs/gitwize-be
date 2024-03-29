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
    },],
    "newCodePercentage": {
      "currentPeriod": 5,
      "previousPeriod": 10
    },
    "churnPercentage": {
      "currentPeriod": 17,
      "previousPeriod": 23
	},
	"legacyPercentage": {
		"currentPeriod": 17.33,
		"previousPeriod": 23.44
	},
	"fileChanged": {
		"currentPeriod": 20,
		"previousPeriod": 30
	},
	"insertionPoints": {
		"currentPeriod": 10,
		"previousPeriod": 15
	},
	"additions": {
		"currentPeriod": 10,
		"previousPeriod": 15
	},
	"deletions": {
		"currentPeriod": 10,
		"previousPeriod": 15
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
	ImpactPeriod     DatePeriod      `json:"period"`
	ImpactScore      ImpactMetric    `json:"impactScore"`
	ActiveDays       ImpactMetric    `json:"activeDays"`
	CommitsPerDay    ImpactMetric    `json:"commitsPerDay"`
	MostChurnedFiles []db.FileChurn  `json:"mostChurnedFiles"`
	NewCodePercent   ImpactMetric    `json:"newCodePercentage"`
	ChurnPercent     ImpactMetric    `json:"churnPercentage"`
	LegacyPercent    ImpactMetric    `json:"legacyPercentage"`
	FileChanged      ImpactMetric    `json:"fileChanged"`
	InsertionPoints  ImpactMetric    `json:"insertionPoints"`
	Additions        ImpactMetric    `json:"additions"`
	Deletions        ImpactMetric    `json:"deletions"`
	UnusualFiles     []db.FileDetail `json:"unusualFiles"`
}

func getWeeklyImpact(c *gin.Context) {
	repoID := c.Param("id")
	if !validateRepoUser(c, repoID) {
		return
	}

	values, err := getIntParams(c, "date_from")
	if err != nil {
		return
	}
	from := values[0]
	t := time.Unix(int64(from), 0)

	currentDuration := getWeekRange(t.UTC())
	prevDuration := getWeekRange(t.UTC().AddDate(0, 0, -7))

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

	unusualfiles, err := db.GetUnusualFiles(repoID, currentDuration.from, currentDuration.to)
	if hasErrUnknown(c, err) {
		return
	}

	currentNewCodePercent, currentChurnPercent := getNewCodeAndChurnPercentage(currentModification)
	prevNewCodePercent, prevChurnPercent := getNewCodeAndChurnPercentage(prevModification)

	weeklyData := WeeklyImpactData{
		ImpactPeriod:     getDatePeriod(currentDuration),
		ImpactScore:      getImpactScore(currentStat, prevStat, currentModification, prevModification),
		ActiveDays:       getActiveDays(currentStat, prevStat),
		CommitsPerDay:    getCommitsPerDay(currentStat, prevStat),
		MostChurnedFiles: mostChurnedFiles,
		NewCodePercent:   ImpactMetric{currentNewCodePercent, prevNewCodePercent},
		ChurnPercent:     ImpactMetric{currentChurnPercent, prevChurnPercent},
		LegacyPercent:    getDumbLegacyPercent(),
		FileChanged:      getFileChanged(currentStat, prevStat),
		InsertionPoints:  getInsertionPoints(currentStat, prevStat),
		Additions:        getAdditions(currentStat, prevStat),
		Deletions:        getDeletions(currentStat, prevStat),
		UnusualFiles:     unusualfiles,
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
	numEditLocation := durationStat.Insertions
	numPercentageNewcode := 0.0
	if durationStat.Addtions != 0 {
		numPercentageNewcode = float64(modificationStat.Additions) * 100 / float64(durationStat.Addtions)
	}
	impact := 5*float64(durationStat.NumFiles) + 5*float64(numEditLocation) + numPercentageNewcode/10 + float64(durationStat.Addtions-durationStat.Deletions)/10
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

func getNewCodeAndChurnPercentage(stat db.ModificationStat) (newCodePercent, churnPercent float64) {
	totalAddition := stat.Additions + stat.Modifications
	if totalAddition > 0 {
		newCodePercent = float64(stat.Additions) / float64(totalAddition) * 100
		churnPercent = float64(stat.Modifications) / float64(totalAddition) * 100
	}
	return newCodePercent, churnPercent
}

// TODO: update when legacy code is ready in db
func getDumbLegacyPercent() ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  17.33,
		PreviousPeriod: 23.44,
	}
}

func getFileChanged(currentStat, prevStat db.DurationStat) ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  float64(currentStat.NumFiles),
		PreviousPeriod: float64(prevStat.NumFiles),
	}
}

func getAdditions(currentStat, prevStat db.DurationStat) ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  float64(currentStat.Addtions),
		PreviousPeriod: float64(prevStat.Addtions),
	}
}

func getDeletions(currentStat, prevStat db.DurationStat) ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  float64(currentStat.Deletions),
		PreviousPeriod: float64(prevStat.Deletions),
	}
}

func getInsertionPoints(currentStat, prevStat db.DurationStat) ImpactMetric {
	return ImpactMetric{
		CurrentPeriod:  float64(currentStat.Insertions),
		PreviousPeriod: float64(prevStat.Insertions),
	}
}
