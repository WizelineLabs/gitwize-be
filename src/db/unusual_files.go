package db

import (
	"github.com/jinzhu/gorm"
	"time"
)

// GetUnusualFiles get file with addition > 1000 in a commit
func GetUnusualFiles(repID string, since, until time.Time) ([]FileDetail, error) {
	files := make([]FileDetail, 0)
	err := gormDB.
		Where("repository_id = ? AND commit_time_stamp >= ? AND commit_time_stamp <= ? AND addition_loc > 1000", repID, since, until).
		Select("file_name, sum(addition_loc) as addition_loc, sum(deletion_loc) as deletion_loc").
		Group("file_name").Find(&files).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}
	return files, nil
}
