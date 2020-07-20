package controller

import (
	"github.com/gin-gonic/gin"
	"gitwize-be/src/db"
	"net/http"
	"strconv"
	"time"
)

const (
	gwEndPointAdmin      = "/api/v1/admin/"
	gwAdminOp            = ":op_id"
	gwEndPointRepository = "/api/v1/repositories/"
	gwRepoGetPutDel      = ":id"
	gwRepoPost           = gwEndPointRepository
	gwRepoStats          = ":id/stats"
	gwContributorStats   = ":id/contributor"
	gwWeeklyImpact       = ":id/impact/weekly"
	gwCodeVelocity       = ":id/code-velocity"
	gwQuarterlyTrend     = ":id/trends"
)

type AdminOperation int

const (
	UPDATE_METRIC_TABLE AdminOperation = iota + 1
)

const (
	statusDataLoading   = "LOADING"
	statusDataAvailable = "AVAILABLE"
)

const (
	ErrCodeNotAuthenticatedUser = http.StatusBadRequest
	ErrKeyNotAuthenticatedUser  = "common.NotAuthenticatedUser"
	ErrMsgNotAuthenticatedUser  = "User's email does not exist."

	ErrCodeUnauthorized = http.StatusUnauthorized
	ErrKeyUnauthorized  = "common.unauthorized"
	ErrMsgUnauthorized  = "Unauthorized."

	ErrCodeEntityNotFound = http.StatusNotFound
	ErrKeyEntityNotFound  = "common.entityNotFound"
	ErrMsgEntityNotFound  = "Entity not found."

	ErrCodeBadJsonFormat = http.StatusBadRequest
	ErrKeyBadJsonFormat  = "common.badJsonFormat"
	ErrMsgBadJsonFormat  = "Not able to parse json format."

	ErrCodeRepoExisted = http.StatusConflict
	ErrKeyRepoExisted  = "repository.existed"
	ErrMsgRepoExisted  = "Repository existed."

	ErrCodeRepoNotFound = http.StatusNotFound
	ErrKeyRepoNotFound  = "repository.notFound"
	ErrMsgRepoNotFound  = "Repository not found."

	ErrCodeRepoBadCredential = http.StatusForbidden
	ErrKeyRepoBadCredential  = "repository.badCredentials"
	ErrMsgRepoBadCredential  = "Provided repository credentials is invalid."

	ErrCodeRepoInvalidUrl = http.StatusBadRequest
	ErrKeyRepoInvalidUrl  = "repository.invalidURL"
	ErrMsgRepoInvalidUrl  = "Repo URL is invalid."

	ErrKeyUnknownIssue = "Unknown"
)

type RestErr struct {
	ErrorKey     string `json:"errorKey"`
	ErrorMessage string `json:"errorMessage"`
}

type RepoInfoPost struct {
	Name        string   `json:"name"`
	Url         string   `json:"url"  binding:"required"`
	Branches    []string `json:"branches"`
	AccessToken string   `json:"password"`
}

type RepoInfoGet struct {
	ID          int       `json:"id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Url         string    `json:"url"  binding:"required"`
	Status      string    `json:"status"  binding:"required"`
	Branches    []string  `json:"branches" binding:"required"`
	LastUpdated time.Time `json:"last_updated" binding:"required"`
}

func getIntParams(c *gin.Context, params ...string) ([]int, error) {
	values := make([]int, len(params))
	for i, param := range params {
		value, err := strconv.Atoi(c.Query(param))
		if err != nil {
			c.JSON(http.StatusBadRequest, RestErr{
				ErrKeyUnknownIssue,
				err.Error(),
			})
			return nil, err
		}
		values[i] = value
	}
	return values, nil
}

func validateRepoUser(c *gin.Context, repoID string) bool {
	userID := extractUserInfo(c)
	repo := db.Repository{}
	if err := db.GetOneRepoUser(userID, repoID, &repo); err != nil {
		c.JSON(http.StatusInternalServerError, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return false
	}
	if repo.ID == 0 {
		c.JSON(ErrCodeEntityNotFound, RestErr{
			ErrorKey:     ErrKeyEntityNotFound,
			ErrorMessage: ErrMsgEntityNotFound})
		return false
	}
	return true
}

func hasErrUnknown(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(http.StatusInternalServerError, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return true
	}
	return false
}

type TimeRange struct {
	from time.Time
	to   time.Time
}

func getWeekRange(t time.Time) TimeRange {
	monday := t.AddDate(0, 0, -int(t.Weekday())+1)
	begin := monday.Truncate(24 * time.Hour)
	end := begin.Add(7*24*time.Hour - 1*time.Microsecond)
	return TimeRange{begin, end}
}
