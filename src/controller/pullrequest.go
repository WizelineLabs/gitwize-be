package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitwize-be/src/db"
	"net/http"
	"time"
)

type PullRequestSize struct {
	Title       string `json:"title"`
	Size        int    `json:"size"`
	Status      string `json:"status"`
	ReviewTime  int    `json:"review_time"`
	Url         string `json:"url"`
	CreatedDate string `json:"created_date"`
	CreatedBy   string `json:"created_by"`
}

func parseDateFromHour(c *gin.Context, pullRequestHour int) (time.Time, int64) {
	day := (pullRequestHour / 100) % 100
	month := (pullRequestHour / 10000) % 100
	year := pullRequestHour / 1000000
	dateFormat := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	date, err := time.Parse(commonDateFormat, dateFormat)
	if hasErrUnknown(c, err) {
		return time.Time{}, 0
	}

	return date, date.Unix()
}
func getPullRequestSize(c *gin.Context) {
	repoID := c.Param("id")

	values, err := getIntParams(c, "date_from", "date_to")
	if err != nil {
		return
	}
	from, to := int64(values[0]), int64(values[1])

	result := make(map[string][]PullRequestSize)
	pullRequestInfos, err := db.GetPullRequestInfo(repoID, from, to)
	if hasErrUnknown(c, err) {
		return
	}

	for runner := from; runner <= to; runner += 24 * 3600 {
		result[time.Unix(runner, 0).Format(commonDateFormat)] = make([]PullRequestSize, 0)
	}
	for _, pullRequest := range pullRequestInfos {

		createdDate, createdInUnixTs := parseDateFromHour(c, pullRequest.CreatedHour)
		if createdDate.IsZero() {
			return
		}

		runner := createdInUnixTs
		if runner < from {
			runner = from
		}
		if pullRequest.ClosedHour == 0 { // open status
			for ; runner <= to; runner += 24 * 3600 {
				curDate := time.Unix(runner, 0).Format(commonDateFormat)
				result[curDate] = append(result[curDate], PullRequestSize{
					Title:       pullRequest.Title,
					Size:        pullRequest.Addition + pullRequest.Deletion,
					Status:      "opened",
					ReviewTime:  0,
					Url:         pullRequest.Url,
					CreatedDate: createdDate.Format(commonDateFormat),
					CreatedBy:   pullRequest.CreatedBy,
				})
			}
		} else { // merged or rejected
			closedDate, closedInUnixTs := parseDateFromHour(c, pullRequest.ClosedHour)
			if closedDate.IsZero() {
				return
			}
			end := closedInUnixTs
			if end <= to {
				result[closedDate.Format(commonDateFormat)] = append(result[closedDate.Format(commonDateFormat)], PullRequestSize{
					Title:       pullRequest.Title,
					Size:        pullRequest.Addition + pullRequest.Deletion,
					Status:      pullRequest.Status,
					ReviewTime:  pullRequest.ReviewDuration / 60,
					Url:         pullRequest.Url,
					CreatedDate: createdDate.Format(commonDateFormat),
					CreatedBy:   pullRequest.CreatedBy,
				})
				end = closedInUnixTs - 24*3600
			} else {
				end = to
			}

			for ; runner <= end; runner += 24 * 3600 {
				curDate := time.Unix(runner, 0).Format(commonDateFormat)
				result[curDate] = append(result[curDate], PullRequestSize{
					Title:       pullRequest.Title,
					Size:        pullRequest.Addition + pullRequest.Deletion,
					Status:      "opened",
					ReviewTime:  pullRequest.ReviewDuration / 60,
					Url:         pullRequest.Url,
					CreatedDate: createdDate.Format(commonDateFormat),
					CreatedBy:   pullRequest.CreatedBy,
				})
			}
		}
	}
	c.JSON(http.StatusOK, result)
}
