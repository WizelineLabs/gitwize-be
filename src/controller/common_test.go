package controller

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetWeekRange(t *testing.T) {
	str := "Wed, 15 Jul 2020 15:04:05 -0700"
	testDate, _ := time.Parse(time.RFC1123Z, str)
	monday, sunday := getWeekRange(testDate)
	assert.Equal(t, time.Monday, monday.Weekday())
	assert.Equal(t, time.Sunday, sunday.Weekday())
	assert.Equal(t, testDate.Day()-2, monday.Day())
	assert.Equal(t, testDate.Day()+4, sunday.Day())
}
