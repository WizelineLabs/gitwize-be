package db

func GetListRepository(repos *[]Repository) error {
	return gormDB.Find(repos).Error
}
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
