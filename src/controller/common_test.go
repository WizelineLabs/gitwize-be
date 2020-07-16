package controller

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetWeekRange(t *testing.T) {
	str := "Wed, 15 Jul 2020 15:04:05 -0700"
	testDate, _ := time.Parse(time.RFC1123Z, str)
	timeRange := getWeekRange(testDate.UTC())
	begin, end := timeRange.from, timeRange.to
	assert.Equal(t, time.Monday, begin.Weekday())
	assert.Equal(t, time.Sunday, end.Weekday())
	assert.Equal(t, testDate.Day()-2, begin.Day())
	assert.Equal(t, testDate.Day()+4, end.Day())
	assert.Equal(t, begin.Hour(), 0)
	assert.Equal(t, end.Hour(), 23)
}
