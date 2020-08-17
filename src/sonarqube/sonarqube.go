package sonarqube

import (
	"encoding/json"
	"fmt"
	"gitwize-be/src/configuration"
	"gitwize-be/src/db"
	"gitwize-be/src/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mux sync.Mutex
var metrics = []string{
	"bugs",
	"coverage",
	"vulnerabilities",
	"code_smells",
	"duplicated_lines_density",
	"ncloc",
	"sqale_index",
	"duplicated_blocks",
	"cognitive_complexity",
	"complexity",
	"security_hotspots",
}

var metricRatings = []string{
	"alert_status",
	"sqale_rating",
	"reliability_rating",
	"security_rating",
}

var metricRatingsMapping = map[string]string{
	"alert_status":       "quality gate",
	"sqale_rating":       "maintainability",
	"reliability_rating": "reliability",
	"security_rating":    "security",
}

func DelSonarQubeProj(userEmail, repoId string) error {
	utils.Trace()
	sonarQubeInt := db.SonarQube{}
	if err := db.GetSonarQubeInstance(userEmail, repoId, &sonarQubeInt); err != nil {
		return err
	}
	if len(sonarQubeInt.Token) == 0 { // This sonarQube instance has been deleted or not created before
		return nil
	}

	if _, err := performHttpRequest(configuration.CurConfiguration.Endpoint.SonarQubeServer +
		sonarQubeAPIProjectDel + "project=" + sonarQubeInt.ProjectKey); err != nil {
		log.Printf(utils.Trace() + ": Error: " + err.Error())
		return err
	}

	return nil
}

func SetupSonarQube(userEmail, repoId, mainBranch string) error {
	utils.Trace()
	sonarQubeInt := db.SonarQube{}
	if err := db.GetSonarQubeInstance(userEmail, repoId, &sonarQubeInt); err != nil {
		return err
	}
	if sonarQubeInt.Token != "" { // This sonarQube instance has been created before
		return nil
	}
	projectName := strings.Replace(userEmail, "@", "_", -1) + "_" + repoId + "_" + strconv.Itoa(int(time.Now().Unix()))
	if respCreatePrj, err := performHttpRequest(configuration.CurConfiguration.Endpoint.SonarQubeServer +
		sonarQubeAPIProjectCreate + "name=" + projectName + "&project=" + projectName); err != nil {
		log.Printf(utils.Trace() + ": Error: " + err.Error())
		return err
	} else {
		defer respCreatePrj.Body.Close()
		if respCreatePrj.StatusCode == http.StatusOK {
			if respCreateToken, errToken := performHttpRequest(configuration.CurConfiguration.Endpoint.SonarQubeServer +
				sonarQubeAPITokenCreate + "name=" + projectName); errToken != nil {
				log.Printf(utils.Trace() + ": Error: " + err.Error())
				return errToken
			} else {
				defer respCreateToken.Body.Close()
				token := &SonarQubeToken{}
				if err := json.NewDecoder(respCreateToken.Body).Decode(token); err != nil {
					log.Printf(utils.Trace() + ": Error: " + err.Error())
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
	utils.Trace()
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
		"sonar.sources=" + configuration.CurConfiguration.SonarQube.BaseDirectory + repository.Name + "\n" +
		"sonar.projectBaseDir=" + configuration.CurConfiguration.SonarQube.BaseDirectory + repository.Name + "\n" +
		"sonar.login=" + sonarQubeInt.Token + "\n" +
		"sonar.projectKey=" + sonarQubeInt.ProjectKey + "\n"
	oFile.WriteString(sonarConfig)
	oFile.Close()

	cmdLine := configuration.CurConfiguration.SonarQube.ScannerPath
	command := exec.Command(cmdLine)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		log.Println(utils.Trace() + ": Error: " + err.Error())
		return err
	}
	// Unlock mutex ======
	mux.Unlock()

	//Get result into sonarQubeInt
	measureNotReady := "Measure has not been found"
	for _, metricRating := range metricRatings {
		for {
			resp, err := performHttpRequest(configuration.CurConfiguration.Endpoint.SonarQubeServer +
				sonarQubeAPIMetricRating + "project=" + sonarQubeInt.ProjectKey + "&metric=" + metricRating)
			if err != nil {
				log.Println(utils.Trace() + ": Error: " + err.Error())
				return err
			}
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			log.Println("=============\n" + string(body))

			if strings.Contains(string(body), measureNotReady) {
				time.Sleep(5 * time.Second)
				continue
			}
			metricName := metricRatingsMapping[metricRating]
			metricValue := strings.Split(strings.Split(strings.Split(string(body), metricName+"</text>")[2], "</text>")[0], ">")[1]
			switch metricName {
			case "quality gate":
				sonarQubeInt.QualityGates = metricValue
			case "maintainability":
				sonarQubeInt.TechnicalDebtRating = metricValue
			case "reliability":
				sonarQubeInt.BugsRating = metricValue
			case "security":
				sonarQubeInt.VulnerabilitiesRating = metricValue
			}
			break
		}
	}
	//https://sonarqube.gitwize.net/api/measures/component?component=tester_wizeline.com_30_1597326143&metricKeys=code_smells
	metric := strings.Join(metrics[:], ",")
	resp, err := performHttpRequest(configuration.CurConfiguration.Endpoint.SonarQubeServer+
		sonarQubeAPIGetComponentMetric+"component="+sonarQubeInt.ProjectKey+"&metricKeys="+metric, "GET")
	if err != nil {
		log.Println(utils.Trace() + ": Error: " + err.Error())
		return err
	}
	allComponentMetric := Component{}
	json.NewDecoder(resp.Body).Decode(&allComponentMetric)
	log.Printf("=== All metrics : %+v\n", allComponentMetric)
	for _, metricType := range allComponentMetric.AllMeasures.Measure {
		switch metricType.Type {
		case "bugs":
			sonarQubeInt.Bugs, _ = strconv.Atoi(metricType.Value)
		case "coverage":
			sonarQubeInt.Coverage, _ = strconv.ParseFloat(metricType.Value, 64)
		case "vulnerabilities":
			sonarQubeInt.Vulnerabilities, _ = strconv.Atoi(metricType.Value)
		case "code_smells":
			sonarQubeInt.CodeSmells, _ = strconv.Atoi(metricType.Value)
		case "alert_status":
			sonarQubeInt.QualityGates = metricType.Value
		case "sqale_index":
			sonarQubeInt.TechnicalDebt, _ = strconv.Atoi(metricType.Value)
		case "duplicated_lines_density":
			sonarQubeInt.Duplications, _ = strconv.ParseFloat(metricType.Value, 64)
		case "duplicated_blocks":
			sonarQubeInt.DuplicationsBlocks, _ = strconv.Atoi(metricType.Value)
		case "cognitive_complexity":
			sonarQubeInt.CognitiveComplexity, _ = strconv.Atoi(metricType.Value)
		case "complexity":
			sonarQubeInt.CyclomaticComplexity, _ = strconv.Atoi(metricType.Value)
		case "security_hotspots":
			sonarQubeInt.SecurityHotspots, _ = strconv.Atoi(metricType.Value)
		}
	}

	sonarQubeInt.LastUpdated = time.Now()
	if err := db.UpdateSonarQubeInstance(&sonarQubeInt); err != nil {
		return err
	}
	return nil
}
