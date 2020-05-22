package db

func FindRepository(repo *Repository, id string) error {
	return gormDB.First(&repo, id).Error
}

func CreateRepository(repo *Repository) error {
	return gormDB.Create(&repo).Error
}

func UpdateRepository(repo *Repository) error {
	return gormDB.Save(&repo).Error
}

func DeleteRepository(repo *Repository) error {
	return gormDB.Delete(&repo).Error
}
