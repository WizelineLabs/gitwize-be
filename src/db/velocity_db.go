package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

//GetCodeChangeVelocity get netChanges
func GetCodeVelocity(id string, start time.Time, end time.Time) (map[string]string, map[string]string, error) {
	start = beginningOfMonth(start)
	end = endOfMonth(end)
	cvdbEntities := make([]CodeVelocityDBEntity, 0)

	err := gormDB.Debug().
		Select("month, SUM(addition_loc)-SUM(deletion_loc) as net_changes, COUNT(*) as no_commits").
		Where("repository_id = ? AND commit_time_stamp >= ? AND commit_time_stamp <= ?", id, start, end).
		Group("year, month").
		Order("year, month").
		Find(&cvdbEntities).Error

	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, nil, err
		}
	}

	netChangesMap := make(map[string]string)
	noCommitsMap := make(map[string]string)
	for _, v := range cvdbEntities {
		netChangesMap[time.Month(v.Month).String()] = v.NetChanges
		noCommitsMap[time.Month(v.Month).String()] = v.NoCommits
	}

	return noCommitsMap, netChangesMap, nil
}

func beginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

func endOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}
