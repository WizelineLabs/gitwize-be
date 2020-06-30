package db

func findRepositoryBaseId(repo *Repository, id string) error {
	return gormDB.First(&repo, id).Error
}

func findRepositoryBaseRepoName(repo *Repository, repoFullName string) error {
	return gormDB.Where("repo_full_name = ?", repoFullName).Find(repo).Error
}
func createRepository(repo *Repository) error {
	return gormDB.Create(&repo).Error
}

func updateRepository(repo *Repository) error {
	return gormDB.Save(&repo).Error
}

func deleteRepository(repo *Repository) error {
	return gormDB.Delete(&repo).Error
}
