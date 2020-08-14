package db

import "github.com/jinzhu/gorm"

func CreateSonarQubeInstance(sonarQubeInt *SonarQube) error {

	return gormDB.Create(sonarQubeInt).Error
}

func GetSonarQubeInstance(userEmail, repoId string, sonarQubeInt *SonarQube) error {
	err := gormDB.Where("user_email = ? AND repository_id = ?", userEmail, repoId).Find(sonarQubeInt).Error
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return nil
	}

	return err
}

func UpdateSonarQubeInstance(sonarQube *SonarQube) error {

	return gormDB.Model(sonarQube).Where("user_email = ? AND repository_id = ?",
		sonarQube.UserEmail, sonarQube.RepoId).Update(sonarQube).Error
}
