package controller

import (
	"net/http"
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
