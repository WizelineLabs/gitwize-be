package db

import (
	"github.com/jinzhu/gorm"
	"gitwize-be/src/utils"
	"log"
	"strconv"
)

func GetReposUser(userEmail string, repos *[]Repository) error {
	users := make([]User, 0)
	if err := gormDB.Where("user_email = ?", userEmail).Find(&users).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
		return err
	}
	if len(users) > 0 {
		for _, user := range users {
			var repository Repository
			if err := findRepositoryBaseId(&repository, strconv.Itoa(user.RepoId)); err != nil {
				log.Println(utils.GetFuncName() + ": Error " + err.Error())
			} else {
				repository.Branches = user.Branches
				*repos = append(*repos, repository)
			}
		}
	}
	return nil
}

func GetOneRepoUser(userEmail string, repoId string, repo *Repository) error {
	user := User{}
	if err := gormDB.Where("user_email = ? AND repo_id = ?", userEmail, repoId).Find(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
		return err
	}

	if err := findRepositoryBaseId(repo, strconv.Itoa(user.RepoId)); err != nil {
		return err
	}
	repo.Branches = user.Branches

	return nil
}
func CreateRepoUser(userEmail string, repo *Repository) error {
	existedRepo := Repository{}
	user := User{
		UserEmail:    userEmail,
		RepoFullName: repo.RepoFullName,
		AccessToken:  repo.AccessToken,
		Branches:     repo.Branches,
	}

	if err := findRepositoryBaseRepoName(&existedRepo, repo.RepoFullName); err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return err
		}
		if err := createRepository(repo); err != nil {
			utils.Trace()
			return err
		}
		user.RepoId = repo.ID
		existedRepo = *repo
	}
	repo.ID = existedRepo.ID
	repo.Status = existedRepo.Status
	user.RepoId = existedRepo.ID

	if err := gormDB.Create(&user).Error; err != nil {
		utils.Trace()
		return err
	}

	existedRepo.NumRef += 1
	if err := updateRepository(&existedRepo); err != nil {
		utils.Trace()
		return err
	}
	return nil
}

func IsRepoBelongToUser(userEmail string, repoId string) (bool, error) {
	var user User
	if err := gormDB.Where("user_email = ? AND repo_id = ?", userEmail, repoId).
		Find(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
func IsRepoUserExist(userEmail string, repoFullName string) (bool, error) {
	var user User
	if err := gormDB.Where("user_email = ? AND repo_full_name = ?", userEmail, repoFullName).
		Find(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func DeleteRepoUser(userEmail string, repo *Repository) error {
	user := User{}
	if err := gormDB.Debug().Where("user_email = ? AND repo_full_name = ?", userEmail, repo.RepoFullName).
		Delete(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
		return err
	}

	if repo.NumRef == 1 {
		if err := DeleteMetricsInOneRepo(repo.ID); err != nil {
			return err
		}
		if err := deleteRepository(repo); err != nil {
			return err
		}
	} else {
		repo.NumRef -= 1
		if err := updateRepository(repo); err != nil {
			return err
		}
	}
	return nil
}
