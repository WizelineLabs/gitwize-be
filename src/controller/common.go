package controller

import "time"

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

type RepoInfoPost struct {
	Name     string   `json:"name"`
	Url      string   `json:"url"  binding:"required"`
	Status   string   `json:"status"`
	Branches []string `json:"branches"`
	User     string   `json:"username"`
	Password string   `json:"password"`
}

type RepoInfoGet struct {
	ID          int       `json:"id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Url         string    `json:"url"  binding:"required"`
	Status      string    `json:"status"  binding:"required"`
	Branches    []string  `json:"branches" binding:"required"`
	LastUpdated time.Time `json:"last_updated" binding:"required"`
}
