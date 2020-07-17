package db

import (
	"github.com/jinzhu/gorm"
	"time"
)

const numberFileLimit = 1

// GetFileChurn get most edited file, increase numberFile to get more files
func GetFileChurn(repID string, since, until time.Time) ([]FileChurn, error) {
	fileChurns := make([]FileChurn, 0)
	err := gormDB.
		Where("repository_id = ? AND commit_time_stamp >= ? AND commit_time_stamp <= ?", repID, since, until).
		Select("COUNT(*) as count, file_name").
		Group("file_name").Order("count DESC, file_name").Limit(numberFileLimit).Find(&fileChurns).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}
	return fileChurns, nil
}
