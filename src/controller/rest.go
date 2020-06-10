package controller

import (
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"gitwize-be/src/auth"
	"gitwize-be/src/configuration"
	"gitwize-be/src/cypher"
	"gitwize-be/src/db"
	"gitwize-be/src/lambda"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func posAdminOperation(c *gin.Context) {
	opId, err := strconv.Atoi(c.Param("op_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Admin key does not exist"})
		return
	}

	adminKey := strings.Split(authHeader, "Bearer ")[1]
	if adminKey != os.Getenv("ADMIN_OP_KEY") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Admin key is not correct"})
		return
	}
	switch AdminOperation(opId) {
	case UPDATE_METRIC_TABLE:
		lambda.CollectPRs()
		lambda.UpdateCommitDataAllRepos()
		db.UpdateMetricTable()
		c.JSON(http.StatusOK, gin.H{"message": "Updating metric table success"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Admin Operation"})
	}
}
func getRepos(c *gin.Context) {
	id := c.Param("id")
	var repo db.Repository
	if err := db.FindRepository(&repo, id); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if repo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository " + id + " doesn't exist"})
	} else {
		c.JSON(http.StatusOK, RepoInfoGet{
			ID:          repo.ID,
			Name:        repo.Name,
			Url:         repo.Url,
			Status:      repo.Status,
			LastUpdated: repo.CtlModifiedDate,
		})
	}
}

func getListRepos(c *gin.Context) {
	var repos []db.Repository
	if err := db.GetListRepository(&repos); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
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

	password := strings.TrimSpace(reqInfo.Password)
	if password != "" { // if password empty, dont ecrypt to use default token later
		password = cypher.EncryptString(password, configuration.CurConfiguration.Cypher.PassPhase)
	}

	createdRepos := db.Repository{
		Name:                 reqInfo.Name,
		Url:                  reqInfo.Url,
		Status:               reqInfo.Status,
		UserName:             reqInfo.User,
		Password:             password,
		CtlCreatedBy:         reqInfo.User,
		CtlCreatedDate:       time.Now(),
		CtlModifiedBy:        reqInfo.User,
		CtlModifiedDate:      time.Now(),
		CtlLastMetricUpdated: time.Now(),
	}
	if err := db.CreateRepository(&createdRepos); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
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
		return
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
		return
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
		return
	}

	to, err := strconv.Atoi(c.Query("date_to"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var repo db.Repository
	if err := db.FindRepository(&repo, idRepository); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if repo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository " + idRepository + " doesn't exist"})
		return
	} else {
		result, err := db.GetMetricBaseOnType(idRepository, metricTypeVal, int64(from), int64(to))
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
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

	admin := ginCont.Group(gwEndPointAdmin)
	{
		admin.POST(gwAdminOp, posAdminOperation)
	}

	repoApi := ginCont.Group(gwEndPointRepository)
	{
		repoApi.Use(AuthMiddleware)
		repoApi.GET(gwRepoPost, getListRepos)
		repoApi.GET(gwRepoGetPutDel, getRepos)
		repoApi.POST(gwRepoPost, postRepos)
		repoApi.PUT(gwRepoGetPutDel, putRepos)
		repoApi.DELETE(gwRepoGetPutDel, delRepos)
		repoApi.GET(gwRepoStats, getStats)
	}

	return ginCont
}
