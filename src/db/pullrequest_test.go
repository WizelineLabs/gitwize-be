package db

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func Test_GetPullRequestInfo_OK(t *testing.T) {
	gormDB, _ = gorm.Open("mysql", LocalDBConnString)
	defer gormDB.Close()

	from := time.Date(2020, 5, 1, 0, 0, 0, 0, time.Local).Unix()
	to := time.Date(2020, 5, 5, 23, 59, 59, 0, time.Local).Unix()
	nodata, _ := GetPullRequestInfo("1", from, to)
	assert.Empty(t, nodata)

	to = time.Date(2020, 5, 12, 23, 59, 59, 0, time.Local).Unix()
	data, _ := GetPullRequestInfo("1", from, to)
	log.Println("data", from, to, data)
	assert.Equal(t, data[0].Title, "GWZ-23 verifies access token")
	assert.Equal(t, data[0].Status, "merged")
	assert.Equal(t, data[0].Addition, 113)
	assert.Equal(t, data[0].Deletion, 7)
	assert.Equal(t, data[0].Status, "merged")
	assert.Equal(t, data[0].ReviewDuration, 322766)
	assert.Equal(t, data[0].Url, "https://github.com/wizeline/gitwize-be/pull/1")
	assert.Equal(t, data[0].CreatedHour, 2020050711)
	assert.Equal(t, data[0].ClosedHour, 2020051104)
}

func Test_GetPullRequestInfo_Err(t *testing.T) {
	gormDB, _ = gorm.Open("mysql", InvalidLocalDBConnString)
	defer gormDB.Close()

	from := time.Date(2020, 5, 1, 0, 0, 0, 0, time.Local).Unix()
	to := time.Date(2020, 5, 5, 23, 59, 59, 0, time.Local).Unix()
	nodata, err := GetPullRequestInfo("1", from, to)
	assert.Empty(t, nodata)
	assert.NotEmpty(t, err)
}
