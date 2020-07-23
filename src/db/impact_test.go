package db

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestGetCommitDurationStat(t *testing.T) {
	gormDB, _ = gorm.Open("mysql", LocalDBConnString) // need to init gormDB again

	from, _ := time.Parse("2006-01-02", "2020-06-19")
	to, _ := time.Parse("2006-01-02", "2020-06-01")
	nodata, _ := GetCommitDurationStat("1", from, to)
	assert.Empty(t, nodata)

	to, _ = time.Parse("2006-01-02", "2020-06-22")
	data, _ := GetCommitDurationStat("1", from, to)
	log.Println("data", from, to, data)
	assert.Equal(t, data.ActiveDays, 2)
	assert.Equal(t, data.TotalCommits, 2)
	assert.Equal(t, data.Insertions, 20)
	gormDB.Close()
}

func TestGetModificationStat(t *testing.T) {
	gormDB, _ = gorm.Open("mysql", LocalDBConnString) // need to init gormDB again

	from, _ := time.Parse("2006-01-02", "2020-06-19")
	to, _ := time.Parse("2006-01-02", "2020-06-01")
	nodata, _ := GetModificationStat("1", from, to)
	assert.Empty(t, nodata)

	to, _ = time.Parse("2006-01-02", "2020-06-22")
	data, _ := GetModificationStat("1", from, to)
	log.Println("data", from, to, data)
	assert.Equal(t, data.TableName(), tableModification)
	assert.Equal(t, data.Modifications, 7)
	gormDB.Close()
}

func TestImpactErr(t *testing.T) {
	gormDB, _ = gorm.Open("mysql", InvalidLocalDBConnString) // need to init gormDB again
	from, _ := time.Parse("2006-01-02", "2020-06-19")
	to, _ := time.Parse("2006-01-02", "2020-06-01")
	_, err := GetCommitDurationStat("1", from, to)
	assert.NotEmpty(t, err)
	_, err = GetModificationStat("1", from, to)
	assert.NotEmpty(t, err)
	gormDB.Close()
}
