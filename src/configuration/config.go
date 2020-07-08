package configuration

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Configurations exported
type Configurations struct {
	Server   ServerConfigurations
	Database DatabaseConfigurations
	Auth     AuthConfigurations
	Cypher   CypherConfigurations
	Endpoint EndpointConfigurations
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
	Frontend string
}

var CurConfiguration Configurations

func ReadConfiguration() {
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	deployEnv := viper.GetString(gwDeployEnv)
	var gwDbPasswordEnv string
	var cypherPassPhaseEnv string
	// Set the file name of the configurations file
	switch deployEnv {
	case devEnvironment:
		viper.SetConfigName(configDev)
		gwDbPasswordEnv = gwDbPasswordDev
		cypherPassPhaseEnv = cypherPassPhaseDev
	case qaEnvironment:
		viper.SetConfigName(configQA)
		gwDbPasswordEnv = gwDbPasswordQA
		cypherPassPhaseEnv = cypherPassPhaseQA
	case prodEnvironment:
		viper.SetConfigName(configPROD)
		gwDbPasswordEnv = gwDbPasswordPROD
		cypherPassPhaseEnv = cypherPassPhasePROD
	default:
		viper.SetConfigName(configLocal)
		gwDbPasswordEnv = gwDbPasswordLocal
		cypherPassPhaseEnv = cypherPassPhaseLocal
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
}
