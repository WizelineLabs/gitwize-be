package controller

import (
	"github.com/gin-gonic/gin"
	"gitwize-be/src/db"
	"net/http"
	"time"
)

type Rating struct {
	Value  int    `json:"value"`
	Rating string `json:"rating"`
}

type CodeQualityDTO struct {
	QualityGates         string    `json:"qualityGates"`
	Bugs                 Rating    `json:"bugs"`
	Vulnerabilities      Rating    `json:"vulnerabilities"`
	CodeSmells           int       `json:"codeSmells"`
	Coverage             float64   `json:"codeCoveragePercentage"`
	Duplications         float64   `json:"duplicationPercentage"`
	DuplicationsBlocks   int       `json:"duplicationBlocks"`
	CognitiveComplexity  int       `json:"cognitiveComplexity"`
	CyclomaticComplexity int       `json:"cyclomaticComplexity"`
	SecurityHotspots     int       `json:"securityHotspots"`
	TechnicalDebt        Rating    `json:"technicalDebt"`
	LastUpdated          time.Time `json:"lastUpdated"`
}

func getCodeQuality(c *gin.Context) {
	userEmail := extractUserInfo(c)
	if userEmail == "" {
		return
	}
	repositoryId := c.Param("id")

	sonarQubeInt := db.SonarQube{}
	if !hasErrUnknown(c, db.GetSonarQubeInstance(userEmail, repositoryId, &sonarQubeInt)) {
		c.JSON(http.StatusOK, CodeQualityDTO{
			QualityGates:         sonarQubeInt.QualityGates,
			CodeSmells:           sonarQubeInt.CodeSmells,
			Coverage:             sonarQubeInt.Coverage,
			Duplications:         sonarQubeInt.Duplications,
			DuplicationsBlocks:   sonarQubeInt.DuplicationsBlocks,
			CognitiveComplexity:  sonarQubeInt.CognitiveComplexity,
			CyclomaticComplexity: sonarQubeInt.CyclomaticComplexity,
			SecurityHotspots:     sonarQubeInt.SecurityHotspots,
			Bugs: Rating{
				Value:  sonarQubeInt.Bugs,
				Rating: sonarQubeInt.BugsRating,
			},
			Vulnerabilities: Rating{
				Value:  sonarQubeInt.Vulnerabilities,
				Rating: sonarQubeInt.VulnerabilitiesRating,
			},
			TechnicalDebt: Rating{
				Value:  sonarQubeInt.TechnicalDebt,
				Rating: sonarQubeInt.TechnicalDebtRating,
			},
			LastUpdated: sonarQubeInt.LastUpdated,
		})
	}
}
