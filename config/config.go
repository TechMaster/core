package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

/*
Trả về true nếu ứng dụng đang chạy ở chế độ Debug và ngược lại
*/
func IsAppInDebugMode() bool {
	appCommand := os.Args[0]
	if strings.Contains(appCommand, "debug") || strings.Contains(appCommand, "exe") {
		return true
	}
	return false
}

func ReadConfig(configPaths ...string) {
	var configPath string
	if len(configPaths) == 0 {
		configPath = "."
	} else {
		configPath = configPaths[0]
	}
	if IsAppInDebugMode() {
		viper.SetConfigName("config.dev") // Debug - development mode
	} else {
		viper.SetConfigName("config.product") // Product mode
	}

	viper.SetConfigType("json") // REQUIRED if the config file does not have the extension in the name

	viper.AddConfigPath(configPath) // optionally look for config in the working directory
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}
