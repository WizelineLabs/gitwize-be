package db

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestGetUnusualFiles(t *testing.T) {
	gormDB, _ = gorm.Open("mysql", LocalDBConnString) // need to init gormDB again

	from, _ := time.Parse("2006-01-02", "2020-07-01")
	to, _ := time.Parse("2006-01-02", "2020-06-01")
	nodata, _ := GetUnusualFiles("1", from, to)
	assert.Empty(t, nodata)

	to, _ = time.Parse("2006-01-02", "2020-07-08")
	data, _ := GetUnusualFiles("1", from, to)
	log.Println("data", from, to, data)
	assert.Equal(t, data[0].FileName, "unusualfile")
	assert.Equal(t, data[0].Additions, 1001)
	gormDB.Close()
}

func TestGetUnusualFilesErr(t *testing.T) {
	gormDB, _ = gorm.Open("mysql", InvalidLocalDBConnString) // need to init gormDB again
	from, _ := time.Parse("2006-01-02", "2020-06-19")
	to, _ := time.Parse("2006-01-02", "2020-06-21")
	_, err := GetUnusualFiles("1", from, to)
	assert.NotEmpty(t, err)
	gormDB.Close()
}
