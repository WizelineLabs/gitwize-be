package db

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestGetListContributors(t *testing.T) {
	gormDB, _ = gorm.Open("mysql", LocalDBConnString) // need to init gormDB again

	from, _ := time.Parse("2006-01-02", "2020-06-19")
	to, _ := time.Parse("2006-01-02", "2020-06-01")
	nodata, _ := GetListContributors("1", from, to)
	assert.Empty(t, nodata)

	to, _ = time.Parse("2006-01-02", "2020-06-21")
	data, _ := GetListContributors("1", from, to)
	log.Println("data", from, to, data)
	assert.Equal(t, data[0].Email, "test@wizeline.com")
	assert.Equal(t, data[0].Name, "test")
	gormDB.Close()
}

func TestGetContributorStatsByPerson(t *testing.T) {
	gormDB, _ = gorm.Open("mysql", LocalDBConnString) // need to init gormDB again

	from, _ := time.Parse("2006-01-02", "2020-06-19")
	to, _ := time.Parse("2006-01-02", "2020-06-01")
	nodata, _ := GetContributorStatsByPerson("1", from, to)
	assert.Empty(t, nodata)

	to, _ = time.Parse("2006-01-02", "2020-06-21")
	data, _ := GetContributorStatsByPerson("1", from, to)
	log.Println("ContributorStatsByPersonData", from, to, data)
	assert.Equal(t, data[0].Email, "test@wizeline.com")
	assert.Equal(t, data[0].Name, "test")
	assert.Equal(t, data[0].Commits, 1)
	assert.Equal(t, data[0].AdditionLoc, 32)
	assert.Equal(t, data[0].DeletionLoc, 20)
	assert.Equal(t, data[0].LOCPercent, float32(0.0))
	gormDB.Close()
}

func TestGetTotalContributorStats(t *testing.T) {
	gormDB, _ = gorm.Open("mysql", LocalDBConnString) // need to init gormDB again

	from, _ := time.Parse("2006-01-02", "2020-06-19")
	to, _ := time.Parse("2006-01-02", "2020-06-01")
	nodata, _ := GetTotalContributorStats("1", from, to)
	assert.Empty(t, nodata)

	to, _ = time.Parse("2006-01-02", "2020-06-21")
	data, _ := GetTotalContributorStats("1", from, to)
	log.Println("TotalStats Data", from, to, data)
	assert.Equal(t, data[0].Commits, 1)
	assert.Equal(t, data[0].AdditionLoc, 32)
	assert.Equal(t, data[0].DeletionLoc, 20)
	assert.Equal(t, data[0].LOCPercent, float32(0.0))
	gormDB.Close()
}

func TestContributorErr(t *testing.T) {
	gormDB, _ = gorm.Open("mysql", InvalidLocalDBConnString) // need to init gormDB again
	from, _ := time.Parse("2006-01-02", "2020-06-19")
	to, _ := time.Parse("2006-01-02", "2020-06-01")
	_, err := GetListContributors("1", from, to)
	assert.NotEmpty(t, err)
	_, err = GetContributorStatsByPerson("1", from, to)
	assert.NotEmpty(t, err)
	_, err = GetTotalContributorStats("1", from, to)
	assert.NotEmpty(t, err)
	gormDB.Close()
}
