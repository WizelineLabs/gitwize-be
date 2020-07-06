package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

//GetCodeChangeVelocity get netChanges
func GetNetChanges(id string, start time.Time, end time.Time) (map[string]string, error) {
	start = beginningOfMonth(start)
	end = endOfMonth(end)
	netChanges := make([]NetChange, 0)

	err := gormDB.Debug().
		Select("month, SUM(addition_loc)-SUM(deletion_loc) as value").
		Where("repository_id = ? AND commit_time_stamp >= ? AND commit_time_stamp <= ?", id, start, end).
		Group("year, month").
		Order("year, month").
		Find(&netChanges).Error

	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}

	result := make(map[string]string)
	for _, v := range netChanges {
		result[time.Month(v.Month).String()] = v.Value
	}
	return result, nil
}

func beginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

func endOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}
