package config

import (
	"encoding/json"
	"os"

	"github.com/Karagar/cyprusExercise/pkg/structs"
	"github.com/Karagar/cyprusExercise/pkg/utils"
)

var config *structs.Config

// getConfig - singleton config wrapper
func New() *structs.Config {
	if config != nil {
		return config
	}

	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = "../scripts/config.json"
	}

	newConfig := &structs.Config{}
	configContent := utils.MustReadFile(configPath)
	utils.PanicOnErr(json.Unmarshal(configContent, newConfig))
	config = newConfig
	return config
}
