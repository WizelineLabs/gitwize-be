package db

import (
	"gitwize-be/src/utils"
	"log"
	"strconv"
)

func GetReposUser(userEmail string, repos *[]Repository) error {
	users := make([]User, 0)
	if err := gormDB.Where("userEmail = ?", userEmail).Find(&users).Error; err != nil {
		return err
	}
	if len(users) > 0 {
		for _, user := range users {
			var repository Repository
			if err := findRepository(&repository, strconv.Itoa(user.RepoId)); err != nil {
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
	if err := gormDB.Where("userEmail = ? AND repo_id = ?", userEmail, repoId).Find(&user).Error; err != nil {
		return err
	}
	if user.RepoId > 0 {
		if err := findRepository(repo, strconv.Itoa(user.RepoId)); err != nil {
			return err
		}
		repo.Branches = user.Branches
	}
	return nil
}
func CreateRepoUser(userEmail string, repo *Repository) (int, error) {
	if err := createRepository(repo); err != nil {
		return 0, err
	}
	return 0, nil
}

func IsRepoUserExist(userEmail string, repoFullName string) (bool, error) {
	return false, nil
}

func UpdateRepoUser(userEmail string, repoFullName string, repo *Repository) error {
	return nil
}
func DeleteRepoUser(userEmail string, repoFullName string) error {

	return nil
}
