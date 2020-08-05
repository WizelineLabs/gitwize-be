package sonarqube

import (
	"encoding/json"
	"fmt"
	"gitwize-be/src/configuration"
	"gitwize-be/src/db"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mux sync.Mutex

func SetupSonarQube(userEmail, repoId, mainBranch string) error {
	sonarQubeInt := db.SonarQube{}
	if err := db.GetSonarQubeInstance(userEmail, repoId, &sonarQubeInt); err != nil {
		return err
	}
	if sonarQubeInt.Token != "" { // This sonarQube instance has been created before
		return nil
	}
	projectName := strings.Replace(userEmail, "@", "_", -1) + "_" + repoId + "_" + strconv.Itoa(int(time.Now().Unix()))
	if respCreatePrj, err := performHttpRequest(sonarQubeAPIProjectCreate + "name=" + projectName + "&project=" + projectName); err != nil {
		return err
	} else {
		defer respCreatePrj.Body.Close()
		if respCreatePrj.StatusCode == http.StatusOK {
			if respCreateToken, errToken := performHttpRequest(configuration.CurConfiguration.Endpoint.SonarQubeServer +
				sonarQubeAPITokenCreate + "name=" + projectName); errToken != nil {
				return errToken
			} else {
				defer respCreateToken.Body.Close()
				token := &SonarQubeToken{}
				if err := json.NewDecoder(respCreateToken.Body).Decode(token); err != nil {
					fmt.Println(err.Error())
					return err
				}
				fmt.Printf("Respond token: %+v\n", token)
				sonarQubeInt = db.SonarQube{
					UserEmail:   userEmail,
					RepoId:      repoId,
					ProjectKey:  projectName,
					Token:       token.Token,
					Branch:      mainBranch,
					LastUpdated: time.Now(),
				}
				return db.CreateSonarQubeInstance(&sonarQubeInt)
			}
		}
	}

	return nil
}

func ScanAndUpdateResult(userEmail, repoId string) error {
	sonarQubeInt := db.SonarQube{}
	if err := db.GetSonarQubeInstance(userEmail, repoId, &sonarQubeInt); err != nil {
		return err
	}

	repository := db.Repository{}
	if err := db.GetOneRepoUser(userEmail, repoId, &repository); err != nil {
		return err
	}
	cloneRepo(repository.Name, repository.Url, repository.AccessToken)

	// Lock mutex ======
	mux.Lock()
	// Edit sonar property file
	os.Remove(configuration.CurConfiguration.SonarQube.PropertiesPath)

	oFile, _ := os.Create(configuration.CurConfiguration.SonarQube.PropertiesPath)
	sonarConfig := "sonar.host.url=" + configuration.CurConfiguration.Endpoint.SonarQubeServer + "\n" +
		"sonar.sources=" + curDirectory + repository.Name + "\n" +
		"sonar.projectBaseDir=" + curDirectory + repository.Name + "\n" +
		"sonar.login=" + sonarQubeInt.Token + "\n" +
		"sonar.projectKey=" + sonarQubeInt.ProjectKey + "\n"
	oFile.WriteString(sonarConfig)
	oFile.Close()

	cmdLine := configuration.CurConfiguration.SonarQube.ScannerPath
	command := exec.Command(cmdLine)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return err
	}

	// Remove repository after scanning
	removeRepoDir := curDirectory + repository.Name
	removeCmd := "rm -rf " + removeRepoDir
	command = exec.Command(removeCmd)
	if err := command.Run(); err != nil {
		return err
	}
	// Unlock mutex ======
	mux.Unlock()
	//Get result into sonarQubeInt
	//https://sonarqube.gitwize.net/api/project_badges/measure\?project\=gitwize-be\&metric\=code_smells
	var metrics = []string{"bugs", "coverage", "vulnerabilities", "code_smells"}
	var resp *http.Response
	var err error
	for _, metric := range metrics {
		resp, err = performHttpRequest(configuration.CurConfiguration.Endpoint.SonarQubeServer +
			sonarQubeAPIGetMetric + "project=" + sonarQubeInt.ProjectKey + "&metric=" + metric)
		if err != nil {
			return err
		}
		body, _ := ioutil.ReadAll(resp.Body)
		metric = strings.Replace(metric, "_", " ", -1)
		metricValue := strings.Split(strings.Split(strings.Split(string(body), metric+"</text>")[2], "</text>")[0], ">")[1]
		fmt.Printf("metric = %s\n", metricValue)
		switch metric {
		case "bugs":
			sonarQubeInt.Bugs, _ = strconv.Atoi(metricValue)
		case "coverage":
			sonarQubeInt.Coverage, _ = strconv.ParseFloat(metricValue, 64)
		case "vulnerabilities":
			sonarQubeInt.Vulnerabilities, _ = strconv.Atoi(metricValue)
		case "code smells":
			sonarQubeInt.CodeSmells, _ = strconv.Atoi(metricValue)
		}
	}
	resp.Body.Close()
	sonarQubeInt.LastUpdated = time.Now()
	if err := db.UpdateSonarQubeInstance(&sonarQubeInt); err != nil {
		return err
	}
	return nil
}
