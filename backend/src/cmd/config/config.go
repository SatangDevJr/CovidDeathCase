package config

import (
	"os"
	"reflect"
)

type Configuration struct {
	PORT        string `env:"PORT"`
	AppVersion  string `env:"APP_VERSION"`
	ELSURL      string `env:"ELS_URL"`
	ELSUsername string `env:"ELS_USERNAME"`
	ELSPassword string `env:"ELS_PASSWORD"`
	ELSIndex    string `env:"ELS_INDEX"`
	Stage       string `env:"STAGE"`
	DCCAPI 		string `env:"DDC_API"`
	Origin      string `env:"HOST_WEB"`
}

func New() Configuration {
	conf := Configuration{}
	v := reflect.ValueOf(&conf).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		envKey := fieldType.Tag.Get("env")
		envValue, ok := os.LookupEnv(envKey)

		switch ok {
		case true:
			field.SetString(envValue)
		case false:
			field.SetString(fieldType.Tag.Get("default"))
		}
	}

	return conf
}
