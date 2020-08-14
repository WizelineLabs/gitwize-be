package configuration

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Configurations exported
type Configurations struct {
	Server    ServerConfigurations
	Database  DatabaseConfigurations
	Auth      AuthConfigurations
	Cypher    CypherConfigurations
	Endpoint  EndpointConfigurations
	SonarQube SonarQubeConfigurations
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Port string
}

type CypherConfigurations struct {
	PassPhase string
}

// DatabaseConfigurations exported
type DatabaseConfigurations struct {
	GwDbName     string
	GwDbUser     string
	GwDbPassword string
	GwDbHost     string
	GwDbPort     int
}

// AuthConfigurations exported
type AuthConfigurations struct {
	AuthDisable string
}

// Endpoint FE exported
type EndpointConfigurations struct {
	Frontend        string
	SonarQubeServer string
}

type SonarQubeConfigurations struct {
	AdminSecret    string
	ScannerPath    string
	PropertiesPath string
}

var CurConfiguration Configurations

func ReadConfiguration() {
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	deployEnv := viper.GetString(gwDeployEnv)
	var gwDbPasswordEnv string
	var cypherPassPhaseEnv string
	var sonarQubeAdminSecret string
	// Set the file name of the configurations file
	switch deployEnv {
	case devEnvironment:
		viper.SetConfigName(configDev)
		gwDbPasswordEnv = gwDbPasswordDev
		cypherPassPhaseEnv = cypherPassPhaseDev
		sonarQubeAdminSecret = sonarQubeAdminSecretDev
	case qaEnvironment:
		viper.SetConfigName(configQA)
		gwDbPasswordEnv = gwDbPasswordQA
		cypherPassPhaseEnv = cypherPassPhaseQA
		sonarQubeAdminSecret = sonarQubeAdminSecretQA
	case prodEnvironment:
		viper.SetConfigName(configPROD)
		gwDbPasswordEnv = gwDbPasswordPROD
		cypherPassPhaseEnv = cypherPassPhasePROD
		sonarQubeAdminSecret = sonarQubeAdminSecretPROD
	default:
		viper.SetConfigName(configLocal)
		gwDbPasswordEnv = gwDbPasswordLocal
		cypherPassPhaseEnv = cypherPassPhaseLocal
		sonarQubeAdminSecret = sonarQubeAdminSecretLocal
	}

	// Set the path to look for the configurations file
	viper.AddConfigPath(configPathFromRootDir)
	viper.AddConfigPath(configPathFromSubModules)

	viper.SetConfigType(configTypeYaml)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	// override with Env variables if any
	for _, k := range viper.AllKeys() {
		envValue := os.Getenv(k)
		if envValue != "" {
			fmt.Printf("Config key overridden by env variable: %s", k)
			viper.Set(k, envValue)
		}
	}

	err := viper.Unmarshal(&CurConfiguration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
	CurConfiguration.Database.GwDbPassword = viper.GetString(gwDbPasswordEnv)
	CurConfiguration.Cypher.PassPhase = viper.GetString(cypherPassPhaseEnv)
	CurConfiguration.SonarQube.AdminSecret = viper.GetString(sonarQubeAdminSecret)
}
