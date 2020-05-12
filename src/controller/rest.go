package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gitwize-be/src/db"
	"net/http"
	"time"
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

func Initialize() *gin.Engine {
	dbConnect = db.Initilize()

	ginCont := gin.Default()
	ginCont.GET(gwEndPointGetPutDel, getRepos)
	ginCont.POST(gwEndPointPost, postRepos)
	ginCont.PUT(gwEndPointGetPutDel, putRepos)
	ginCont.DELETE(gwEndPointGetPutDel, delRepos)
	ginCont.GET(statsEndPoint, getStats)
	return ginCont
}
