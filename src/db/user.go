package db

import (
	"gitwize-be/src/utils"
	"log"
	"strconv"
)

func GetReposUser(userEmail string, repos *[]Repository) error {
	users := make([]User, 0)
	if err := gormDB.Where("user_email = ?", userEmail).Find(&users).Error; err != nil {
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
		return err
	}
	if user.RepoId > 0 {
		if err := findRepositoryBaseId(repo, strconv.Itoa(user.RepoId)); err != nil {
			return err
		}
		repo.Branches = user.Branches
	}
	return nil
}
func CreateRepoUser(userEmail string, repo *Repository) error {
	var existedRepo Repository
	user := User{
		UserEmail:    userEmail,
		RepoFullName: repo.RepoFullName,
		AccessToken:  repo.AccessToken,
		Branches:     repo.Branches,
	}

	if err := findRepositoryBaseRepoName(&existedRepo, repo.RepoFullName); err != nil {
		return err
	}
	if existedRepo.ID == 0 {
		// Not created  before
		if err := createRepository(repo); err != nil {
			return err
		}
		user.RepoId = repo.ID
		existedRepo = *repo
	} else {
		repo.ID = existedRepo.ID
		repo.Status = existedRepo.Status
		user.RepoId = existedRepo.ID
	}

	if err := gormDB.Create(&user).Error; err != nil {
		return err
	}

	existedRepo.NumRef += 1
	return updateRepository(&existedRepo)
}

func IsRepoUserExist(userEmail string, repoFullName string) (bool, error) {
	user := User{}
	if err := gormDB.Where("user_email = ? AND repo_full_name = ?", userEmail, repoFullName).Find(&user).Error; err != nil {
		return false, err
	}
	if user.RepoId > 0 {
		return true, nil
	}
	return false, nil
}

func UpdateRepoUser(userEmail string, repoFullName string, repo *Repository) error {
	return nil
}
func DeleteRepoUser(userEmail string, repoFullName string) error {

	return nil
}
