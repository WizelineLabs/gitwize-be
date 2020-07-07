package db

import "github.com/jinzhu/gorm"

//GetContributorStatsByPerson get statistics by contributor
func GetContributorStatsByPerson(id string, since string, until string) ([]ContributorStats, error) {

	contributorStats := make([]ContributorStats, 0)

	err := gormDB.
		Where("repository_id = ? AND commit_time_stamp >= ? AND commit_time_stamp <= ?", id, since, until).
		Select(
			"repository_id, " +
				"LOWER(author_email) as author_email, " +
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

//GetTotalContributorStats get statistics for repo group by day
func GetTotalContributorStats(id string, since string, until string) ([]ContributorStats, error) {
	contributorStats := make([]ContributorStats, 0)
	err := gormDB.
		Where("repository_id = ? AND commit_time_stamp >= ? AND commit_time_stamp <= ?", id, since, until).
		Select(
			"repository_id, " +
				"COUNT(*) as commits, " +
				"SUM(addition_loc) as addition_loc, " +
				"SUM(deletion_loc) as deletion_loc, " +
				"SUM(num_files) as num_files, " +
				"Date(commit_time_stamp) as date").
		Group("repository_id, year, month, day").
		Find(&contributorStats).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}
	return contributorStats, nil
}

// GetListContributors get list contributor email and name, name will be concatnate for same email
func GetListContributors(id string, since string, until string) ([]Contributor, error) {
	contributors := make([]Contributor, 0)
	err := gormDB.
		Where("repository_id = ? AND commit_time_stamp >= ? AND commit_time_stamp <= ?", id, since, until).
		Select("DISTINCT(LOWER(author_email)) as author_email, CONCAT_WS(\", \", author_name) as author_name").
		Group("author_email").Order("author_name").Find(&contributors).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}
	return contributors, nil
}
