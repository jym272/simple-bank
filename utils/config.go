package utils

import "github.com/spf13/viper"

type Config struct {
	DBDriver      string `default:"postgres" mapstructure:"DB_DRIVER"`
	DBSource      string `default:"postgres://postgres:postgres@localhost:8080/simple_bank?sslmode=disable" mapstructure:"DB_SOURCE"`
	ServerAddress string `default:":8081" mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return

}
