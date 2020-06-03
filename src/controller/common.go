package controller

import "time"

const (
	gwEndPointAdmin     = "/api/v1/admin/:op_id"
	gwEndPoint          = "/api/v1/repositories/"
	gwEndPointGetPutDel = gwEndPoint + "/:id"
	gwEndPointPost      = gwEndPoint
	statsEndPoint       = gwEndPoint + "/:id/stats"
)

type AdminOperation int

const (
	UPDATE_METRIC_TABLE AdminOperation = iota + 1
)

type RepoInfoPost struct {
	Name     string `json:"name"`
	Url      string `json:"url"  binding:"required"`
	Status   string `json:"status"`
	User     string `json:"username"`
	Password string `json:"password"`
}

type RepoInfoGet struct {
	ID          uint      `json:"id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Url         string    `json:"url"  binding:"required"`
	Status      string    `json:"status"  binding:"required"`
	LastUpdated time.Time `json:"last_updated" binding:"required"`
}
