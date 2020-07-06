package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

//GetCodeChangeVelocity get netChanges
func GetNetChanges(id string, start time.Time, end time.Time) (map[string]string, error) {
	startMonth, startYear := int(start.Month()), start.Year()
	endMonth, endYear := int(end.Month()), end.Year()
	netChanges := make([]NetChange, 0)

	err := gormDB.
		Select("month, SUM(addition_loc)-SUM(deletion_loc) as value").
		Where("repository_id = ? AND year >= ? AND year <= ? AND month >= ? AND month <= ?", id, startYear, endYear, startMonth, endMonth).
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
