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
	PRS_OPEN
	PRS_MERGED
	PRS_REJECTED
)

var MapTypeMetric = map[string]MetricsType{
	"ALL":           ALL,
	"loc":           LOC,
	"lines_added":   LINES_ADDED,
	"lines_removed": LINES_REMOVED,
	"commits":       COMMITS,
	"prs_open":      PRS_OPEN,
	"prs_merged":    PRS_MERGED,
	"prs_rejected":  PRS_REJECTED,
}

const (
	tableRepository = "repository"
	tableMetric     = "metric"
)

type Repository struct {
	ID              uint      `gorm:"column:id;primary_key" json:"id"`
	Name            string    `gorm:"column:name" json:"name"`
	Url             string    `gorm:"column:url" json:"url"`
	Status          string    `gorm:"column:status" json:"status"`
	UserName        string    `gorm:"column:username" json:"username"`
	Password        string    `gorm:"column:password" json:"password"`
	CtlCreatedDate  time.Time `gorm:"column:ctl_created_date" json:"ctl_created_date"`
	CtlCreatedBy    string    `gorm:"column:ctl_created_by" json:"ctl_created_by"`
	CtlModifiedDate time.Time `gorm:"column:ctl_modified_date" json:"ctl_modified_date"`
	CtlModifiedBy   string    `gorm:"column:ctl_modified_by" json:"ctl_modified_by"`
	Metrics         []Metric
}

func (Repository) TableName() string {
	return tableRepository
}

type Metric struct {
	ID           uint        `gorm:"column:id;primary_key" json:"id"`
	RepositoryID uint        `gorm:"column:repository_id" json:"repository_id"`
	BranchName   string      `gorm:"column:branch;index:branch" json:"branch"`
	Type         MetricsType `gorm:"column:type" json:"type"`
	Value        uint        `gorm:"column:value" json:"value"`
	AsOfDate     time.Time   `gorm:"column:as_of_date" json:"as_of_date"`
}

func (Metric) TableName() string {
	return tableMetric
}
