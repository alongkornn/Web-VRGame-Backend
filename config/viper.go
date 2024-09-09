package config

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/spf13/viper"
// )

// func InitConfig() {
// 	viper.SetConfigName("config")    // name of config file (without extension)
// 	viper.SetConfigType("yaml")      // REQUIRED if the config file does not have the extension in the name
// 	viper.AddConfigPath("./config/") // path to look for the config file in the current working directory
// 	viper.AutomaticEnv()             // read in environment variables that match
// 	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

// 	// Read the config file
// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		panic(fmt.Errorf("fatal error config file: %s", err))
// 	}
// }

// // GetEnv fetches the environment variable value for the given key
// func GetEnv(key string) string {
// 	return viper.GetString(key)
// }
