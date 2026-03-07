package util
import (
	"github.com/spf13/viper"
)
//config stores all configrations of the application 
// The values are read by viper from a config file or environment variables.
type Config struct {
DBDriver string `mapstructure:"DB_DRIVER"`
DBSource string  `mapstructure:"DB_SOURCE"`
ServerAdress string `mapstructure:"SERVER_ADDRESS"`
}

//LoadConfig reads configirations from file or environment vraiables
func LoadConfig(path string) (Config Config ,err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&Config)
	return
}
