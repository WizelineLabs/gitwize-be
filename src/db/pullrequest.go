package db

import (
	"github.com/jinzhu/gorm"
	"time"
)

func GetPullRequestInfo(idRepository string, epochFrom int64, epochTo int64) ([]PullRequestInfo, error) {
	dateFrom := time.Unix(epochFrom, 0)
	dateTo := time.Unix(epochTo, 0)
	yearFrom, monthFrom, dayFrom, hourFrom := dateFrom.Year(), int(dateFrom.Month()), dateFrom.Day(), dateFrom.Hour()
	from := ((yearFrom*100+monthFrom)*100+dayFrom)*100 + hourFrom
	yearTo, monthTo, dayTo, hourTo := dateTo.Year(), int(dateTo.Month()), dateTo.Day(), dateTo.Hour()
	to := ((yearTo*100+monthTo)*100+dayTo)*100 + hourTo

	pullRequestSizes := make([]PullRequestInfo, 0)
	if err := gormDB.Debug().
		Select("title,url,state,additions,deletions,review_duration,created_hour,closed_hour,created_by").
		Where("repository_id = ? AND ("+
			"(closed_hour > 0 AND ((closed_hour >= ? AND closed_hour <= ?) OR (closed_hour > ? AND created_hour <= ?)))"+
			"OR (closed_hour = 0 AND created_hour <= ?)"+
			")", idRepository, from, to, to, to, to).
		Find(&pullRequestSizes).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}

	return pullRequestSizes, nil
}
