package db

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	userEmail    = "testerGitwize@wizeline.com"
	repositoryId = "100"
)

func Test_SonarQubeDB_OK(t *testing.T) {

	gormDB, _ = gorm.Open("mysql", LocalDBConnString)
	defer gormDB.Close()

	sonarQubeInt := SonarQube{}
	// Get sonarQube with empty value
	err := GetSonarQubeInstance(userEmail, repositoryId, &sonarQubeInt)
	assert.Nil(t, err)

	// Create new sonarQubeInt
	sonarQubeInt.UserEmail = userEmail
	sonarQubeInt.RepoId = repositoryId
	err = CreateSonarQubeInstance(&sonarQubeInt)
	assert.Nil(t, err)

	// Get error when trying creating the same sonarQube again
	err = CreateSonarQubeInstance(&sonarQubeInt)
	assert.NotNil(t, err)

	// Get sonarQubeInt
	sonarResult := SonarQube{}
	err = GetSonarQubeInstance(userEmail, repositoryId, &sonarResult)
	assert.Nil(t, err)
	assert.Equal(t, sonarResult.UserEmail, userEmail)
	assert.Equal(t, sonarResult.RepoId, repositoryId)

	// Update sonarQubeInt
	sonarQubeInt.Bugs = 5
	err = UpdateSonarQubeInstance(&sonarQubeInt)
	assert.Nil(t, err)

	// Get again, then comparing
	err = GetSonarQubeInstance(userEmail, repositoryId, &sonarResult)
	assert.Nil(t, err)
	assert.Equal(t, sonarQubeInt.Bugs, sonarResult.Bugs)

	// Delete sonar Instance
	err = DelSonarQubeIntance(userEmail, repositoryId)
	assert.Nil(t, err)
	err = GetSonarQubeInstance(userEmail, repositoryId, &sonarResult)
	assert.Nil(t, err)
}
