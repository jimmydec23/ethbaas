package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var C = viper.GetViper()

func init() {
	C.AutomaticEnv()
	C.SetConfigName("config")
	C.SetConfigType("yaml")

	// current dir or parent dir
	C.AddConfigPath(".")
	C.AddConfigPath("..")

	cp := os.Getenv("CONFIG_PATH")
	if cp != "" {
		C.AddConfigPath(cp)
	}

	C.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := C.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
