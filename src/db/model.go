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

const (
	tableRepository = "repository"
	tableMetric     = "metric"
)

type Repository struct {
	ID              uint      `json:"id" gorm:"primary_key"`
	Name            string    `json:"name"`
	Url             string    `json:"url"`
	Status          string    `json:"status"`
	UserName        string    `json:"username"`
	Password        string    `json:"password"`
	CtlCreatedDate  time.Time `json:"ctl_created_date"`
	CtlCreatedBy    string    `json:"ctl_created_by"`
	CtlModifiedDate time.Time `json:"ctl_modified_date"`
	CtlModifiedBy   string    `json:"ctl_modified_by"`
	Metrics         []Metric
}

func (Repository) TableName() string {
	return tableRepository
}

type Metric struct {
	ID           uint        `json:"id" gorm:"primary_key"`
	RepositoryID uint        `json:"repository_id"`
	BranchName   string      `json:"branch" gorm:"index:branch"`
	Type         MetricsType `json:"type"`
	Value        uint        `json:"value"`
	AsOfDate     time.Time   `json:"as_of_date"`
}

func (Metric) TableName() string {
	return tableMetric
}
