package db

import "github.com/jinzhu/gorm"

//GetChartContributorStats get statistics by contributor
func GetChartContributorStats(id string, since string, until string) ([]ContributorStats, error) {

	contributorStats := make([]ContributorStats, 0)

	err := gormDB.
		Where("repository_id = ? AND commit_time_stamp >= ? AND commit_time_stamp <= ?", id, since, until).
		Select(
			"repository_id, " +
				"author_email, " +
				"COUNT(*) as commits, " +
				"CONCAT_WS(\", \", author_name) as author_name, " +
				"SUM(addition_loc) as addition_loc, " +
				"SUM(deletion_loc) as deletion_loc, " +
				"SUM(num_files) as num_files, " +
				"Date(commit_time_stamp) as date").
		Group("repository_id, author_email, year, month, day").
		Find(&contributorStats).Error

	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}
	return contributorStats, nil
}

// GetListContributors get list contributor email and name, name will be concatnate for same email
func GetListContributors(id string, since string, until string) ([]ContributorName, error) {
	contributorNames := make([]ContributorName, 0)
	err := gormDB.
		Where("repository_id = ? AND commit_time_stamp >= ? AND commit_time_stamp <= ?", id, since, until).
		Select("distinct(author_email) as author_email, CONCAT_WS(\", \", author_name) as author_name").
		Group("author_email").Find(&contributorNames).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}
	return contributorNames, nil
}
