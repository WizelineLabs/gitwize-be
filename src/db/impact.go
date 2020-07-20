package db

import (
	"github.com/jinzhu/gorm"
	"time"
)

/* Example
SELECT
COUNT(DISTINCT(Date(commit_time_stamp))) AS active_days,
COUNT(*) AS total_commits,
SUM(num_files) as num_files,
SUM(addition_loc) as additions,
SUM(deletion_loc) as deletions
FROM commit_data
WHERE repository_id = 2 AND num_parents = 1;
*/

// GetCommitDurationStat get statistics summary in a period
func GetCommitDurationStat(repoID string, from, to time.Time) (DurationStat, error) {
	durationStat := DurationStat{}
	err := gormDB.
		Where("repository_id = ? AND commit_time_stamp >= ? AND commit_time_stamp <= ? AND num_parents = 1", repoID, from, to).
		Select("COUNT(DISTINCT(Date(commit_time_stamp))) AS active_days, COUNT(*) AS total_commits, " +
			"SUM(num_files) as num_files, SUM(addition_loc) as additions, SUM(deletion_loc) as deletions").
		Find(&durationStat).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return DurationStat{}, err
		}
	}
	return durationStat, nil
}

//GetModificationStat get modification stats
func GetModificationStat(repoID string, from, to time.Time) (ModificationStat, error) {
	modificationStat := ModificationStat{}
	err := gormDB.
		Where("repository_id = ? AND commit_time_stamp >= ? AND commit_time_stamp <= ?", repoID, from, to).
		Select("SUM(modification_loc) as modifications").
		Find(&modificationStat).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return ModificationStat{}, err
		}
	}
	return modificationStat, nil
}
