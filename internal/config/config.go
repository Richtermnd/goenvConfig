package config

import (
	"fmt"
	"reflect"

	"github.com/joho/godotenv"
)

func LoadConfig(cfg interface{}) error {
	cfgType := reflect.TypeOf(cfg)
	// Struct check
	if cfgType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("%s is not struct", cfgType.Kind())
	}
	cfgIndirect := reflect.Indirect(reflect.ValueOf(cfg))
	return FillConfig(&cfgIndirect, cfgType)
}

func FillConfig(cfg *reflect.Value, cfgType reflect.Type) error {
	// Explore fields
	for i := 0; i < cfg.NumField(); i++ {
		// Get field
		field := cfg.Field(i)

		// ignore private fields
		if !field.CanSet() {
			continue
		}

		// Parse field
		fieldParams, err := ParseField(cfgType.Elem().Field(i))
		if err != nil {
			return err
		}

		// Upload value from enviroment
		value, err := fieldParams.GetValue()
		if err != nil {
			return err
		}
		// Set value to field
		field.Set(value)
	}
	return nil
}

// LoadEnviroment upload env vars from .env files
// Just a wrap for godotenv.Load
func LoadEnviroment(envFiles ...string) error {
	return godotenv.Load(envFiles...)
}
