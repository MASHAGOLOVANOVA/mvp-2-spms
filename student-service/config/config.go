package config

import (
	"encoding/json"
	"os"
	"reflect"
)

type ServerConfiguration struct {
	Address   string `json:"addr" env:"SERVER_ADDRESS"`
	Port      string `json:"port" env:"SERVER_PORT"`
	ReturnUrl string `json:"return_url" env:"RETURN_URL"`
}

func ReadConfigFromFile(filename string) (ServerConfiguration, error) {
	f, err := os.Open(filename)
	if err != nil {
		return ServerConfiguration{}, err
	}

	var config ServerConfiguration
	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&config)
	if err != nil {
		return ServerConfiguration{}, err
	}

	return config, err
}

func SetConfigEnvVars(config ServerConfiguration) error {
	ct := reflect.TypeOf(config)
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Field(i)
		fieldValue := reflect.ValueOf(config).Field(i)
		tag, _ := field.Tag.Lookup("env")
		err := os.Setenv(tag, fieldValue.String())
		if err != nil {
			return err
		}
	}
	return nil
}
