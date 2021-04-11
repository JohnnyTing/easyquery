package config

import (
	"github.com/ztrue/tracerr"

	"github.com/spf13/viper"
)

type Properties struct {
	Database Database `json:"database"`
	App      App      `json:"app"`
}

type Database struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Port     string `json:"port"`
	Host     string `json:"host"`
	Database string `json:"database"`
}

type App struct {
	Mode string `json:"mode"`
	Port string `json:"port"`
}

var GlobalProperties *Properties

func init() {
	var (
		err        error
		profile    string
		properties Properties
	)

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.BindEnv("PROFILE")

	profile = viper.GetString("PROFILE")
	if profile == "" {
		profile = "dev"
	}

	// Set the file name of the configurations file
	viper.SetConfigName(profile)

	// Set the path to look for the configurations file
	viper.AddConfigPath(ConfPath)

	viper.SetConfigType("yml")

	if err = viper.ReadInConfig(); err != nil {
		err = tracerr.Errorf("Error reading configs file, %s", err)
		panic(err)
	}

	err = viper.Unmarshal(&properties)
	if err != nil {
		err = tracerr.Errorf("Unable to decode into struct, %v", err)
		panic(err)
	}

	GlobalProperties = &properties
}
