package lambda

import (
	"log"
	"time"
)

const (
	commitTable   = "commit_data"
	fileStatTable = "file_stat_data"
	gitDateFormat = "2006-01-02"
	batchSize     = 1000
	directory     = "/tmp/"
)

type DateRange struct {
	Since *time.Time
	Until *time.Time
}

func GetFullGitDateRange() DateRange {
	since := time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2040, 1, 1, 0, 0, 0, 0, time.UTC)
	return DateRange{Since: &since, Until: &until}
}

func GetLastNDayDateRange(n int) DateRange {
	nDayBefore := time.Now().AddDate(0, 0, -n)
	tomorrow := time.Now().AddDate(0, 0, +1)
	since := time.Date(nDayBefore.Year(), nDayBefore.Month(), nDayBefore.Day(), 0, 0, 0, 0, time.UTC)
	until := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.UTC)
	return DateRange{Since: &since, Until: &until}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("\n%s took %s", name, elapsed)
}
