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
	PRS_CREATED
	PRS_MERGED
	PRS_REJECTED
)

var MapNameToTypeMetric = map[string]MetricsType{
	"ALL":           ALL,
	"loc":           LOC,
	"lines_added":   LINES_ADDED,
	"lines_removed": LINES_REMOVED,
	"commits":       COMMITS,
	"prs_created":   PRS_CREATED,
	"prs_merged":    PRS_MERGED,
	"prs_rejected":  PRS_REJECTED,
}

var MapTypeMetricToName = map[MetricsType]string{
	LOC:           "loc",
	LINES_ADDED:   "lines_added",
	LINES_REMOVED: "lines_removed",
	COMMITS:       "commits",
	PRS_CREATED:   "prs_created",
	PRS_MERGED:    "prs_merged",
	PRS_REJECTED:  "prs_rejected",
}

const (
	tableRepository = "repository"
	tableMetric     = "metric"
)

type Repository struct {
	ID                   uint      `gorm:"column:id;primary_key" json:"id"`
	Name                 string    `gorm:"column:name" json:"name"`
	Url                  string    `gorm:"column:url" json:"url"`
	Status               string    `gorm:"column:status" json:"status"`
	UserName             string    `gorm:"column:username" json:"username"`
	Password             string    `gorm:"column:password" json:"password"`
	CtlCreatedDate       time.Time `gorm:"type:timestamp;column:ctl_created_date" json:"ctl_created_date"`
	CtlCreatedBy         string    `gorm:"column:ctl_created_by" json:"ctl_created_by"`
	CtlModifiedDate      time.Time `gorm:"type:timestamp;column:ctl_modified_date" json:"ctl_modified_date"`
	CtlModifiedBy        string    `gorm:"column:ctl_modified_by" json:"ctl_modified_by"`
	CtlLastMetricUpdated time.Time `gorm:"type:timestamp;column:ctl_last_metric_updated" json:"ctl_last_metric_updated"`
}

func (Repository) TableName() string {
	return tableRepository
}

type RepositoryDTO struct {
	ID      uint                `json:"id"`
	Name    string              `json:"name"`
	Url     string              `json:"url"`
	Status  string              `json:"status"`
	Metrics map[string][]Metric `json:"metric"`
}

type Metric struct {
	ID               uint        `gorm:"column:id;primary_key" json:"id"`
	RepositoryID     int         `gorm:"column:repository_id" json:"repository_id"`
	BranchName       string      `gorm:"column:branch;index:branch" json:"branch"`
	Type             MetricsType `gorm:"column:type" json:"type"`
	Value            uint64      `gorm:"column:value" json:"value"`
	ContributorEmail string      `gorm:"column:contributor_email" json:"contributor_email"`
	Year             uint        `gorm:"column:year" json:"year"`
	Month            uint        `gorm:"column:month" json:"month"`
	Day              uint        `gorm:"column:day" json:"day"`
	Hour             uint        `gorm:"column:hour" json:"hour"`
}

type MetricDTO struct {
	ID           uint        `json:"id"`
	RepositoryID int         `json:"repository_id"`
	BranchName   string      `json:"branch"`
	Type         MetricsType `json:"type"`
	Value        uint64      `json:"value"`
	AsOfDate     time.Time   `json:"as_of_date"`
}

func (Metric) TableName() string {
	return tableMetric
}
