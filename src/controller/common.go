package controller

const (
	gwEndPoint          = "/api/v1/repositories/"
	gwEndPointGetPutDel = gwEndPoint + "/:id"
	gwEndPointPost      = gwEndPoint
	statsEndPoint       = gwEndPoint + "/:id/stats"
)

type RepoInfoPost struct {
	Name     string `json:"name" binding:"required"`
	Url      string `json:"url"  binding:"required"`
	Status   string `json:"status"  binding:"required"`
	User     string `json:"username" binding:"required"`
	Password string `json:"password"`
}

type RepoInfoGet struct {
	Name   string `json:"name" binding:"required"`
	Url    string `json:"url"  binding:"required"`
	Status string `json:"status"  binding:"required"`
}
