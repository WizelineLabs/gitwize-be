package db

import (
	"github.com/jinzhu/gorm"
	"time"
)

// GetCommitDurationStat get statistics summary in a period
func GetCommitDurationStat(repoID string, from, to time.Time) (DurationStat, error) {
	durationStat := DurationStat{}
	err := gormDB.
		Where("repository_id = ? AND commit_time_stamp >= ? AND commit_time_stamp <= ? AND num_parents = 1", repoID, from, to).
		Select("COUNT(DISTINCT(Date(commit_time_stamp))) AS active_days, Count(*) AS total_commits").
		Find(&durationStat).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return DurationStat{}, err
		}
	}
	return durationStat, nil
}
