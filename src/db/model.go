package db

import (
	"time"
)

type MetricsType uint

const (
	ALL MetricsType = iota
	LOC
	LINES_ADDED
	LINES_REMOVED
	COMMITS
	PRS
)

type Repository struct {
	ID              uint      `json:"id" gorm:"primary_key"`
	Name            string    `json:"name"`
	Url             string    `json:"url"`
	Status          string    `json:"status"`
	UserName        string    `json:"username"`
	CtlCreatedDate  time.Time `json:"ctl_created_date"`
	CtlCreatedBy    string    `json:"ctl_created_by"`
	CtlModifiedDate time.Time `json:"ctl_modified_date"`
	CtlModifiedBy   string    `json:"ctl_modified_by"`
	Metrics         []Metric
}

type Metric struct {
	ID         uint        `json:"id" gorm:"primary_key"`
	ReposID    uint        `json:"repos_id"`
	BranchName string      `json:"branch" gorm:"index:branch"`
	Type       MetricsType `json:"metric_type"`
	Value      uint        `json:"value"`
	AsOfDate   time.Time   `json:"as_of_date"`
}
