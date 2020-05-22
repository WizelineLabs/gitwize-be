package configuration

import (
	"fmt"
	"github.com/spf13/viper"
)

// Configurations exported
type Configurations struct {
	Server   ServerConfigurations
	Database DatabaseConfigurations
	Auth     AuthConfigurations
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

// AuthConfigurations exported
type AuthConfigurations struct {
	AuthDisable string
}

var CurConfiguration Configurations

func ReadConfiguration() {
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	deployEnv := viper.GetString(gwDeployEnv)
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

	err := viper.Unmarshal(&CurConfiguration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
	CurConfiguration.Database.GwDbPassword = viper.GetString(gwDbPasswordEnv)
}
