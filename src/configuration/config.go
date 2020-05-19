package configuration

import (
	"fmt"
	"github.com/spf13/viper"
)

// Configurations exported
type Configurations struct {
	Server   ServerConfigurations
	Database DatabaseConfigurations
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Port string
}

// DatabaseConfigurations exported
type DatabaseConfigurations struct {
	GwDbName     string
	GwDbUser     string
	GwDbPassword string
	GwDbHost     string
	GwDbPort     int
}

func ReadConfiguration() Configurations {
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	deployEnv := viper.GetString(gwDeployEnv)
	var configuration Configurations
	var gwDbPasswordEnv string
	// Set the file name of the configurations file
	switch deployEnv {
	case devEnvironment:
		viper.SetConfigName(configDev)
		gwDbPasswordEnv = gwDbPasswordDev
	default:
		viper.SetConfigName(configLocal)
		gwDbPasswordEnv = gwDbPasswordLocal
	}

	// Set the path to look for the configurations file
	viper.AddConfigPath(configPathFromRootDir)
	viper.AddConfigPath(configPathFromSubModules)

	viper.SetConfigType(configTypeYaml)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
	configuration.Database.GwDbPassword = viper.GetString(gwDbPasswordEnv)

	return configuration
}
