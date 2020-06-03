package controller

import (
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"gitwize-be/src/auth"
	"gitwize-be/src/configuration"
	"gitwize-be/src/cypher"
	"gitwize-be/src/db"
	"net/http"
	"strconv"
	"time"
)

func getRepos(c *gin.Context) {
	id := c.Param("id")
	var repo db.Repository
	if err := db.FindRepository(&repo, id); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	if repo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository " + id + " doesn't exist"})
	} else {
		c.JSON(http.StatusOK, repo)
	}
}

func getListRepos(c *gin.Context) {
	var repos []db.Repository
	if err := db.GetListRepository(&repos); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	var repoInfos []RepoInfoGet
	for _, repo := range repos {
		repoInfos = append(repoInfos, RepoInfoGet{
			ID:          repo.ID,
			Name:        repo.Name,
			Url:         repo.Url,
			Status:      repo.Status,
			LastUpdated: repo.CtlModifiedDate,
		})
	}
	c.JSON(http.StatusOK, repoInfos)
}

func postRepos(c *gin.Context) {
	var reqInfo RepoInfoPost
	if err := c.BindJSON(&reqInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdRepos := db.Repository{
		Name:            reqInfo.Name,
		Url:             reqInfo.Url,
		Status:          reqInfo.Status,
		UserName:        reqInfo.User,
		Password:        cypher.EncryptString(reqInfo.Password, configuration.CurConfiguration.Cypher.PassPhase),
		CtlCreatedBy:    reqInfo.User,
		CtlCreatedDate:  time.Now(),
		CtlModifiedBy:   reqInfo.User,
		CtlModifiedDate: time.Now(),
	}
	if err := db.CreateRepository(&createdRepos); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	repoInfo := RepoInfoGet{
		ID:          createdRepos.ID,
		Name:        createdRepos.Name,
		Url:         createdRepos.Url,
		Status:      createdRepos.Status,
		LastUpdated: createdRepos.CtlModifiedDate,
	}

	c.JSON(http.StatusCreated, repoInfo)
}

func putRepos(c *gin.Context) {
	var reqInfo RepoInfoPost
	if err := c.BindJSON(&reqInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	var repo db.Repository
	if err := db.FindRepository(&repo, id); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	if repo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository " + id + " doesn't exist"})
	} else {
		repo.Name = reqInfo.Name
		repo.UserName = reqInfo.User
		repo.Url = reqInfo.Url
		repo.Status = reqInfo.Status
		repo.CtlModifiedBy = reqInfo.User
		repo.CtlModifiedDate = time.Now()
		if err := db.UpdateRepository(&repo); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		c.JSON(http.StatusOK, repo)
	}
}

func delRepos(c *gin.Context) {
	id := c.Param("id")
	var repo db.Repository
	if err := db.FindRepository(&repo, id); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	if repo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository " + id + " doesn't exist"})
	} else {
		if err := db.DeleteRepository(&repo); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		c.JSON(http.StatusNoContent, nil)
	}
}

func getStats(c *gin.Context) {
	idRepository := c.Param("id")
	metricTypeName := c.DefaultQuery("metric_type", "ALL")
	metricTypeVal, ok := db.MapNameToTypeMetric[metricTypeName]
	if !ok {
		metricTypeVal = db.ALL
	}

	from, err := strconv.Atoi(c.Query("date_from"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	to, err := strconv.Atoi(c.Query("date_to"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	var repo db.Repository
	if err := db.FindRepository(&repo, idRepository); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	if repo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository " + idRepository + " doesn't exist"})
	} else {
		result, err := db.GetMetricBaseOnType(idRepository, metricTypeVal, int64(from), int64(to))
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else {
			repositoryDTO := db.RepositoryDTO{
				ID:      repo.ID,
				Name:    repo.Name,
				Status:  repo.Status,
				Url:     repo.Url,
				Metrics: result,
			}
			c.JSON(http.StatusOK, repositoryDTO)
		}
	}
}

// AuthMiddleware checks for valid access token
func AuthMiddleware(c *gin.Context) {
	authDisabled := configuration.CurConfiguration.Auth.AuthDisable == "true"
	if !authDisabled && !auth.IsAuthorized(nil, c.Request) {
		c.AbortWithStatusJSON(401, gin.H{
			"message.key": "system.unauthorized",
			"message":     "Unauthorized!",
		})
	}

	c.Next()
}

func corsHandler() gin.HandlerFunc {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{configuration.CurConfiguration.Endpoint.Frontend},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "OPTIONS", "HEAD", "DELETE"},
		AllowedHeaders:   []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
	})
}

func Initialize() *gin.Engine {
	db.Initialize()

	ginCont := gin.Default()
	ginCont.Use(corsHandler())
	ginCont.Use(AuthMiddleware)
	ginCont.GET(gwEndPoint, getListRepos)
	ginCont.GET(gwEndPointGetPutDel, getRepos)
	ginCont.POST(gwEndPointPost, postRepos)
	ginCont.PUT(gwEndPointGetPutDel, putRepos)
	ginCont.DELETE(gwEndPointGetPutDel, delRepos)
	ginCont.GET(statsEndPoint, getStats)
	return ginCont
}
