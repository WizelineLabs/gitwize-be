package db

func findRepository(repo *Repository, id string) error {
	return gormDB.First(&repo, id).Error
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
