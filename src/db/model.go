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
	PRS_OPENED
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
	"prs_opened":    PRS_OPENED,
}

var MapTypeMetricToName = map[MetricsType]string{
	LOC:           "loc",
	LINES_ADDED:   "lines_added",
	LINES_REMOVED: "lines_removed",
	COMMITS:       "commits",
	PRS_CREATED:   "prs_created",
	PRS_MERGED:    "prs_merged",
	PRS_REJECTED:  "prs_rejected",
	PRS_OPENED:    "prs_opened",
}

const (
	tableRepository      = "repository"
	tableMetric          = "metric"
	tableContributor     = "commit_data"
	tableContributorFile = "file_stat_data"
	tableUser            = "repository_user"
)

type Repository struct {
	ID                   int       `gorm:"column:id;primary_key" json:"id"`
	RepoFullName         string    `gorm:"column:repo_full_name;index:repo_full_name" json:"repo_full_name"`
	Name                 string    `gorm:"column:name" json:"name"`
	Url                  string    `gorm:"column:url" json:"url"`
	Status               string    `gorm:"column:status" json:"status"`
	AccessToken          string    `gorm:"column:access_token" json:"access_token"`
	Branches             string    `gorm:"column:branches" json:"branches"`
	NumRef               int       `gorm:"column:num_ref" json:"num_ref"`
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
	ID      int                    `json:"id"`
	Name    string                 `json:"name"`
	Url     string                 `json:"url"`
	Status  string                 `json:"status"`
	Metrics map[string][]MetricDTO `json:"metric"`
}

type User struct {
	UserEmail    string `gorm:"column:user_email;primary_key" json:"user_email"`
	RepoId       int    `gorm:"column:repo_id" json:"repo_id"`
	RepoFullName string `gorm:"column:repo_full_name" json:"repo_full_name"`
	AccessToken  string `gorm:"column:access_token" json:"access_token"`
	Branches     string `gorm:"column:branches" json:"branches"`
}

func (User) TableName() string {
	return tableUser
}

type Metric struct {
	ID               int         `gorm:"column:id;primary_key" json:"id"`
	RepositoryID     int         `gorm:"column:repository_id" json:"repository_id"`
	BranchName       string      `gorm:"column:branch;index:branch" json:"branch"`
	Type             MetricsType `gorm:"column:type" json:"type"`
	Value            uint64      `gorm:"column:value" json:"value"`
	ContributorEmail string      `gorm:"column:contributor_email" json:"contributor_email"`
	Year             int         `gorm:"column:year" json:"year"`
	Month            int         `gorm:"column:month" json:"month"`
	Day              int         `gorm:"column:day" json:"day"`
	Hour             int         `gorm:"column:hour" json:"hour"`
}

type MetricDTO struct {
	BranchName string      `json:"branch"`
	Type       MetricsType `json:"type"`
	Value      uint64      `json:"value"`
	AsOfDate   string      `json:"as_of_date"`
}

func (Metric) TableName() string {
	return tableMetric
}

type ContributorStats struct {
	RepositoryID int     `gorm:"column:repository_id" json:"repository_id"`
	Email        string  `gorm:"column:author_email" json:"author_email"`
	Name         string  `gorm:"column:author_name" json:"author_name"`
	Commits      int     `gorm:"column:commits" json:"commits"`
	AdditionLoc  int     `gorm:"column:addition_loc" json:"addition_loc"`
	DeletionLoc  int     `gorm:"column:deletion_loc" json:"deletion_loc"`
	NumFiles     int     `gorm:"column:num_files" json:"num_files"`
	LOCPercent   float32 `gorm:"column:loc_percent" json:"loc_percent"`
	Date         string  `gorm:"column:date" json:"date"`
}

type Contributor struct {
	Email string `gorm:"column:author_email" json:"author_email"`
	Name  string `gorm:"column:author_name" json:"author_name"`
}

func (ContributorStats) TableName() string {
	return tableContributor
}

func (Contributor) TableName() string {
	return tableContributor
}
