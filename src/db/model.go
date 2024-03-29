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
	tableNetChange       = "commit_data"
	tablePullRequest     = "pull_request"
	tableFileChurn       = "file_stat_data"
	tableCommitDuration  = "commit_data"
	tableModification    = "file_stat_data"
	tableSonarqube       = "sonarqube"
	tableFileDetail      = "file_stat_data"
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
	Email        string  `gorm:"column:author_email" json:"email"`
	Name         string  `gorm:"column:author_name" json:"name"`
	Commits      int     `gorm:"column:commits" json:"commits"`
	AdditionLoc  int     `gorm:"column:addition_loc" json:"additions"`
	DeletionLoc  int     `gorm:"column:deletion_loc" json:"deletions"`
	NumFiles     int     `gorm:"column:num_files" json:"filesChange"`
	LOCPercent   float32 `gorm:"column:loc_percent" json:"changePercent"`
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

type CodeVelocityDBEntity struct {
	Month     int    `gorm:"column:month"`
	Addtions  string `gorm:"column:additions"`
	Deletions string `gorm:"column:deletions"`
	NoCommits string `gorm:"column:no_commits"`
}

func (CodeVelocityDBEntity) TableName() string {
	return tableNetChange
}

type QuarterlyTrends struct {
	PercentageRejectedPR map[string]int `json:"percentageRejectedPR"`
	AveragePRTime        map[string]int `json:"averagePRTime"`
	AveragePRSize        map[string]int `json:"averagePRSize"`
}

type RejectedMergedPR struct {
	Month int `gorm:"column:month"`
	Value int `gorm:"column:value"`
}

func (RejectedMergedPR) TableName() string {
	return tableMetric
}

type DurationSizePR struct {
	Month    int `gorm:"column:closed_month"`
	Addition int `gorm:"column:additions"`
	Deletion int `gorm:"column:deletions"`
	Duration int `gorm:"column:review_duration"`
}

func (DurationSizePR) TableName() string {
	return tablePullRequest
}

type FileChurn struct {
	FileName string `gorm:"column:file_name" json:"fileName"`
	Value    int    `gorm:"column:count" json:"value"`
}

func (FileChurn) TableName() string {
	return tableFileChurn
}

type DurationStat struct {
	ActiveDays   int `gorm:"column:active_days"`
	TotalCommits int `gorm:"column:total_commits"`
	NumFiles     int `gorm:"column:num_files"`
	Addtions     int `gorm:"column:additions"`
	Deletions    int `gorm:"column:deletions"`
	Insertions   int `gorm:"column:insertion_point"`
}

func (DurationStat) TableName() string {
	return tableCommitDuration
}

type ModificationStat struct {
	Modifications int `gorm:"column:modifications"`
	Additions     int `gorm:"column:additions"`
	Deletions     int `gorm:"column:deletions"`
}

func (ModificationStat) TableName() string {
	return tableModification
}

type PullRequestInfo struct {
	Title          string `gorm:"column:title" json:"title"`
	Url            string `gorm:"column:url" json:"url"`
	Status         string `gorm:"column:state" json:"state"`
	Addition       int    `gorm:"column:additions" json:"additions"`
	Deletion       int    `gorm:"column:deletions" json:"deletions"`
	ReviewDuration int    `gorm:"column:review_duration" json:"review_duration"`
	CreatedHour    int    `gorm:"column:created_hour" json:"created_hour"`
	ClosedHour     int    `gorm:"column:closed_hour" json:"closed_hour"`
	CreatedBy      string `gorm:"column:created_by" json:"created_by"`
}

func (PullRequestInfo) TableName() string {
	return tablePullRequest
}

type SonarQube struct {
	UserEmail             string    `gorm:"column:user_email" json:"user_email"`
	RepoId                string    `gorm:"column:repository_id" json:"repository_id"`
	ProjectKey            string    `gorm:"column:project_key" json:"project_key"`
	Token                 string    `gorm:"column:token" json:"token"`
	Branch                string    `gorm:"column:branch" json:"branch"`
	QualityGates          string    `gorm:"column:quality_gates" json:"quality_gates"`
	Bugs                  int       `gorm:"column:bugs" json:"bugs"`
	BugsRating            string    `gorm:"column:bugs_rating" json:"bugs_rating"`
	Vulnerabilities       int       `gorm:"column:vulnerabilities" json:"vulnerabilities"`
	VulnerabilitiesRating string    `gorm:"column:vulnerabilities_rating" json:"vulnerabilities_rating"`
	CodeSmells            int       `gorm:"column:code_smells" json:"code_smells"`
	Coverage              float64   `gorm:"column:coverage" json:"coverage"`
	Duplications          float64   `gorm:"column:duplications" json:"duplications"`
	DuplicationsBlocks    int       `gorm:"column:duplication_blocks" json:"duplication_blocks"`
	CognitiveComplexity   int       `gorm:"column:cognitive_complexity" json:"cognitive_complexity"`
	CyclomaticComplexity  int       `gorm:"column:cyclomatic_complexity" json:"cyclomatic_complexity"`
	SecurityHotspots      int       `gorm:"column:security_hotspots" json:"security_hotspots"`
	TechnicalDebt         int       `gorm:"column:technical_debt" json:"technical_debt"`
	TechnicalDebtRating   string    `gorm:"column:technical_debt_rating" json:"technical_debt_rating"`
	LastUpdated           time.Time `gorm:"column:last_updated" json:"last_updated"`
}

func (SonarQube) TableName() string {
	return tableSonarqube
}

type FileDetail struct {
	FileName  string `gorm:"column:file_name" json:"fileName"`
	Additions int    `gorm:"column:addition_loc" json:"additions"`
	Deletions int    `gorm:"column:deletion_loc" json:"deletions"`
}

func (FileDetail) TableName() string {
	return tableFileDetail
}
