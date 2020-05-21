package controller

import (
	"gitwize-be/src/auth"
	"gitwize-be/src/db"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var dbConnect *gorm.DB

func getRepos(c *gin.Context) {
	id := c.Param("id")
	var repo db.Repository
	if err := dbConnect.First(&repo, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	if repo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository " + id + " doesn't exist"})
	} else {
		c.JSON(http.StatusOK, repo)
	}
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
		CtlCreatedBy:    reqInfo.User,
		CtlCreatedDate:  time.Now(),
		CtlModifiedBy:   reqInfo.User,
		CtlModifiedDate: time.Now(),
	}
	if err := dbConnect.Create(&createdRepos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusCreated, gin.H{
		"id": createdRepos.ID,
	})
}

func putRepos(c *gin.Context) {
	var reqInfo RepoInfoPost
	if err := c.BindJSON(&reqInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	var repo db.Repository
	if err := dbConnect.First(&repo, id).Error; err != nil {
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
		if err := dbConnect.Save(&repo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		c.JSON(http.StatusOK, repo)
	}
}

func delRepos(c *gin.Context) {
	id := c.Param("id")
	var repo db.Repository
	if err := dbConnect.First(&repo, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	if repo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository " + id + " doesn't exist"})
	} else {
		if err := dbConnect.Delete(&repo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		c.JSON(http.StatusNoContent, nil)
	}
}

func getStats(c *gin.Context) {
	id := c.Param("id")
	metricType := c.DefaultQuery("metric_type", "ALL")
	c.JSON(http.StatusOK, gin.H{
		"message":    "Statistic information is read successfully",
		"id":         id,
		"metricType": metricType,
	})
}

// AuthMiddleware checks for valid access token
func AuthMiddleware(c *gin.Context) {
	authDisabled := os.Getenv("AUTH_DISABLED") == "true"
	if !authDisabled && !auth.IsAuthorized(nil, c.Request) {
		c.AbortWithStatusJSON(401, gin.H{
			"message.key": "system.unauthorized",
			"message":     "Unauthorized!",
		})
	}

	c.Next()
}

func Initialize() *gin.Engine {
	dbConnect = db.Initialize()

	ginCont := gin.Default()
	ginCont.Use(AuthMiddleware)
	ginCont.GET(gwEndPointGetPutDel, getRepos)
	ginCont.POST(gwEndPointPost, postRepos)
	ginCont.PUT(gwEndPointGetPutDel, putRepos)
	ginCont.DELETE(gwEndPointGetPutDel, delRepos)
	ginCont.GET(statsEndPoint, getStats)
	return ginCont
}
