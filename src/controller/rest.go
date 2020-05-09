package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getRepos(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Repository is ok",
		"id":      id,
	})
}

func postRepos(c *gin.Context) {
	var reqInfo RepoInfoPost
	if err := c.BindJSON(&reqInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("Request Information for create repos is %v", reqInfo)
	c.JSON(http.StatusCreated, reqInfo)
}

func putRepos(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Repository is updated successfully",
		"id":      id,
	})
}

func delRepos(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Repository is deleted successfully",
		"id":      id,
	})
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
	ginCont := gin.Default()
	ginCont.GET(gwEndPointGetPutDel, getRepos)
	ginCont.POST(gwEndPointPost, postRepos)
	ginCont.PUT(gwEndPointGetPutDel, putRepos)
	ginCont.DELETE(gwEndPointGetPutDel, delRepos)
	ginCont.GET(statsEndPoint, getStats)
	return ginCont
}
