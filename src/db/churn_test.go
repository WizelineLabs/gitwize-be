package db

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestGetFileChurn(t *testing.T) {
	gormDB, _ = gorm.Open("mysql", LocalDBConnString) // need to init gormDB again

	from, _ := time.Parse("2006-01-02", "2020-06-19")
	to, _ := time.Parse("2006-01-02", "2020-06-01")
	nodata, _ := GetFileChurn("1", from, to)
	assert.Empty(t, nodata)

	to, _ = time.Parse("2006-01-02", "2020-06-21")
	data, _ := GetFileChurn("1", from, to)
	log.Println("data", from, to, data)
	assert.Equal(t, data[0].FileName, "file1")
	assert.Equal(t, data[0].Value, 1)
	gormDB.Close()
}

func TestGetFileChurnErr(t *testing.T) {
	gormDB, _ = gorm.Open("mysql", InvalidLocalDBConnString) // need to init gormDB again
	from, _ := time.Parse("2006-01-02", "2020-06-19")
	to, _ := time.Parse("2006-01-02", "2020-06-21")
	_, err := GetFileChurn("1", from, to)
	assert.NotEmpty(t, err)
	gormDB.Close()
}
