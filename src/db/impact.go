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
SUM(deletion_loc) as deletions,
SUM(insertion_point) as insertion_point,
FROM commit_data
WHERE repository_id = 2 AND num_parents = 1;
*/

// GetCommitDurationStat get statistics summary in a period
func GetCommitDurationStat(repoID string, from, to time.Time) (DurationStat, error) {
	durationStat := DurationStat{}
	err := gormDB.
		Where("repository_id = ? AND commit_time_stamp >= ? AND commit_time_stamp <= ? AND num_parents = 1", repoID, from, to).
		Select("COUNT(DISTINCT(Date(commit_time_stamp))) AS active_days, COUNT(*) AS total_commits, " +
			"SUM(num_files) AS num_files, SUM(addition_loc) AS additions, SUM(deletion_loc) AS deletions," +
			"SUM(insertion_point) AS insertion_point").
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
		Select("SUM(IFNULL(modification_loc,0)) AS modifications, SUM(addition_loc) AS additions, SUM(deletion_loc) AS deletions").
		// Select("SUM(modification_loc) AS modifications, SUM(addition_loc) AS additions, SUM(deletion_loc) AS deletions").
		Find(&modificationStat).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return ModificationStat{}, err
		}
	}
	return modificationStat, nil
}
