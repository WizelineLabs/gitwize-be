package controller

import (
	"gitwize-be/src/auth"
	"gitwize-be/src/configuration"
	"gitwize-be/src/cypher"
	"gitwize-be/src/db"
	"gitwize-be/src/githubapi"
	"gitwize-be/src/lambda"
	"gitwize-be/src/sonarqube"
	"gitwize-be/src/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func extractUserInfo(c *gin.Context) string {
	userId := c.Request.Header.Get("AuthenticatedUser")
	if userId == "" && configuration.CurConfiguration.Auth.AuthDisable == "true" {
		userId = "tester@wizeline.com"
	} else if userId == "" {
		c.JSON(ErrCodeNotAuthenticatedUser, RestErr{
			ErrorKey:     ErrKeyNotAuthenticatedUser,
			ErrorMessage: ErrMsgNotAuthenticatedUser})
		return ""
	}
	return userId
}

func posAdminOperation(c *gin.Context) {
	defer utils.TimeTrack(time.Now(), utils.GetFuncName())
	// TODO: will implement this feature later
}

func getRepos(c *gin.Context) {
	defer utils.TimeTrack(time.Now(), utils.GetFuncName())

	userId := extractUserInfo(c)
	if userId == "" {
		return
	}
	id := c.Param("id")
	repo := db.Repository{}
	if err := db.GetOneRepoUser(userId, id, &repo); err != nil {
		c.JSON(http.StatusInternalServerError, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}

	branches := make([]string, 0)
	if len(repo.Branches) > 0 {
		branches = strings.Split(repo.Branches, ",")
	}
	c.JSON(http.StatusOK, RepoInfoGet{
		ID:          repo.ID,
		Name:        repo.Name,
		Url:         repo.Url,
		Status:      repo.Status,
		Branches:    branches,
		LastUpdated: repo.CtlModifiedDate,
	})
}

func getListRepos(c *gin.Context) {
	defer utils.TimeTrack(time.Now(), utils.GetFuncName())

	userId := extractUserInfo(c)
	if userId == "" {
		return
	}
	repos := make([]db.Repository, 0)
	if err := db.GetReposUser(userId, &repos); err != nil {
		c.JSON(http.StatusInternalServerError, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}

	repoInfos := make([]RepoInfoGet, 0)
	for _, repo := range repos {
		branches := make([]string, 0)
		if len(repo.Branches) > 0 {
			branches = strings.Split(repo.Branches, ",")
		}
		repoInfos = append(repoInfos, RepoInfoGet{
			ID:          repo.ID,
			Name:        repo.Name,
			Url:         repo.Url,
			Status:      repo.Status,
			Branches:    branches,
			LastUpdated: repo.CtlModifiedDate,
		})
	}
	c.JSON(http.StatusOK, repoInfos)
}

func postRepos(c *gin.Context) {
	defer utils.TimeTrack(time.Now(), utils.GetFuncName())
	var reqInfo RepoInfoPost
	var err error
	var branches []string
	var owner, repoName string

	userId := extractUserInfo(c)
	if userId == "" {
		return
	}
	if err = c.BindJSON(&reqInfo); err != nil {
		c.JSON(ErrCodeBadJsonFormat, RestErr{
			ErrorKey:     ErrKeyBadJsonFormat,
			ErrorMessage: err.Error(),
		})
		return
	}

	reqInfo.Url = strings.Replace(reqInfo.Url, "git@github.com:", "https://github.com/", 1)
	if owner, repoName, err = githubapi.ParseGithubUrl(reqInfo.Url); err != nil {
		c.JSON(ErrCodeRepoInvalidUrl, RestErr{
			ErrKeyRepoInvalidUrl,
			ErrMsgRepoInvalidUrl,
		})
		return
	}

	if duplicated, err := db.IsRepoUserExist(userId, owner+"/"+repoName); err != nil {
		utils.Trace()
		c.JSON(http.StatusInternalServerError, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	} else if duplicated {
		c.JSON(ErrCodeRepoExisted, RestErr{
			ErrKeyRepoExisted,
			ErrMsgRepoExisted,
		})
		return
	}

	if branches, err = githubapi.GetListBranches(owner, repoName, reqInfo.AccessToken); err != nil {
		if strings.Contains(err.Error(), "Bad credentials") ||
			strings.Contains(err.Error(), "Resource protected by organization SAML") {
			c.JSON(ErrCodeRepoBadCredential, RestErr{
				ErrKeyRepoBadCredential,
				ErrMsgRepoBadCredential,
			})
		} else if strings.Contains(err.Error(), "Not Found") {
			c.JSON(ErrCodeRepoNotFound, RestErr{
				ErrKeyRepoNotFound,
				ErrMsgRepoNotFound,
			})
		} else {
			c.JSON(http.StatusBadRequest, RestErr{
				ErrKeyUnknownIssue,
				err.Error(),
			})
		}
		return
	}

	accessToken := strings.TrimSpace(reqInfo.AccessToken)
	if accessToken != "" { // if access token is empty, not encrypt and use default token
		accessToken = cypher.EncryptString(accessToken, configuration.CurConfiguration.Cypher.PassPhase)
	}

	createdRepo := db.Repository{
		Name:                 reqInfo.Name,
		RepoFullName:         owner + "/" + repoName,
		Url:                  reqInfo.Url,
		Status:               statusDataLoading,
		AccessToken:          accessToken,
		Branches:             strings.Join(branches, ","),
		NumRef:               0,
		CtlCreatedBy:         userId,
		CtlCreatedDate:       time.Now(),
		CtlModifiedBy:        userId,
		CtlModifiedDate:      time.Now(),
		CtlLastMetricUpdated: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	if err := db.CreateRepoUser(userId, &createdRepo); err != nil {
		c.JSON(http.StatusInternalServerError, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}

	repoPayload := lambda.RepoPayload{
		RepoID:   createdRepo.ID,
		RepoName: createdRepo.Name,
		URL:      createdRepo.Url,
		RepoPass: strings.TrimSpace(reqInfo.AccessToken), // sending non-decrypted value, it will be decrypted on lambda
		Branch:   "",
	}
	lambda.Trigger(repoPayload, lambda.GetLoadFullRepoLambdaFunc(), "ap-southeast-1")

	repoInfo := RepoInfoGet{
		ID:          createdRepo.ID,
		Name:        createdRepo.Name,
		Url:         createdRepo.Url,
		Status:      createdRepo.Status,
		Branches:    branches,
		LastUpdated: createdRepo.CtlModifiedDate,
	}

	go func() { // Start setup sonarQube & Scan
		if err := sonarqube.SetupSonarQube(userId, strconv.Itoa(createdRepo.ID), "master"); err != nil {
			log.Println(utils.Trace() + ": Error:" + err.Error())
		} else {
			sonarqube.ScanAndUpdateResult(userId, strconv.Itoa(createdRepo.ID))
		}
	}()

	c.JSON(http.StatusCreated, repoInfo)
}

func delRepos(c *gin.Context) {
	defer utils.TimeTrack(time.Now(), utils.GetFuncName())
	userId := extractUserInfo(c)
	if userId == "" {
		return
	}

	id := c.Param("id")
	repo := db.Repository{}
	if !hasErrUnknown(c, db.GetOneRepoUser(userId, id, &repo)) &&
		!hasErrUnknown(c, sonarqube.DelSonarQubeProj(userId, strconv.Itoa(repo.ID))) &&
		!hasErrUnknown(c, db.DelSonarQubeIntance(userId, strconv.Itoa(repo.ID))) &&
		!hasErrUnknown(c, db.DeleteRepoUser(userId, &repo)) {
		c.JSON(http.StatusNoContent, nil)
	}

	return
}

func getStats(c *gin.Context) {
	defer utils.TimeTrack(time.Now(), utils.GetFuncName())
	userId := extractUserInfo(c)
	if userId == "" {
		return
	}

	idRepository := c.Param("id")
	metricTypeName := c.DefaultQuery("metric_type", "ALL")
	metricTypeVal, ok := db.MapNameToTypeMetric[metricTypeName]
	if !ok {
		metricTypeVal = db.ALL
	}

	from, err := strconv.Atoi(c.Query("date_from"))
	if err != nil {
		c.JSON(http.StatusBadRequest, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}

	to, err := strconv.Atoi(c.Query("date_to"))
	if err != nil {
		c.JSON(http.StatusBadRequest, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}

	repo := db.Repository{}
	if err := db.GetOneRepoUser(userId, idRepository, &repo); err != nil {
		c.JSON(http.StatusInternalServerError, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
		return
	}

	result, err := db.GetMetricBaseOnType(idRepository, metricTypeVal, int64(from), int64(to))
	if err != nil {
		c.JSON(http.StatusInternalServerError, RestErr{
			ErrKeyUnknownIssue,
			err.Error(),
		})
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

// authMiddleware checks for valid access token
func authMiddleware(c *gin.Context) {
	if !(configuration.CurConfiguration.Auth.AuthDisable == "true") {
		authorized, emailUser := auth.IsAuthorized(nil, c.Request)
		if !authorized {
			c.AbortWithStatusJSON(ErrCodeUnauthorized, RestErr{
				ErrorKey:     ErrKeyUnauthorized,
				ErrorMessage: ErrMsgUnauthorized},
			)
		} else {
			idRepository := c.Param("id")
			if idRepository != "" {
				if isBelongTo, err := db.IsRepoBelongToUser(emailUser, idRepository); err != nil {
					c.JSON(http.StatusInternalServerError, RestErr{
						ErrKeyUnknownIssue,
						err.Error(),
					})
				} else if !isBelongTo {
					c.JSON(ErrCodeEntityNotFound, RestErr{
						ErrorKey:     ErrKeyEntityNotFound,
						ErrorMessage: ErrMsgEntityNotFound})
					return
				}
			}
		}
	}
	c.Next()
}

func corsHandler() gin.HandlerFunc {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{configuration.CurConfiguration.Endpoint.Frontend},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "OPTIONS", "HEAD", "DELETE"},
		AllowedHeaders:   []string{"Origin", "Content-Length", "Content-Type", "Authorization", "AuthenticatedUser"},
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
		repoApi.Use(authMiddleware)
		repoApi.GET("", getListRepos)
		repoApi.GET(gwRepoGetPutDel, getRepos)
		repoApi.POST("", postRepos)
		repoApi.DELETE(gwRepoGetPutDel, delRepos)
		repoApi.GET(gwRepoStats, getStats)
		repoApi.GET(gwContributorStats, getContributorStats)
		repoApi.GET(gwWeeklyImpact, getWeeklyImpact)
		repoApi.GET(gwCodeVelocity, getCodeChangeVelocity)
		repoApi.GET(gwQuarterlyTrend, getStatsQuarterlyTrends)
		repoApi.GET(gwPullRequestSize, getPullRequestSize)
		repoApi.GET(gwCodeQuality, getCodeQuality)
	}

	return ginCont
}
