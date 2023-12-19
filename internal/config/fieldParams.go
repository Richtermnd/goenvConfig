package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type FieldParams struct {
	Name         string
	Required     bool
	DefaultValue string
	Type         reflect.Type
}

func (fp *FieldParams) GetValue() (reflect.Value, error) {
	value, err := fp.getEnv()
	if err != nil {
		return reflect.Value{}, err
	}
	return fp.upcastValue(value)
}

func (fp *FieldParams) getEnv() (string, error) {
	value := os.Getenv(fp.Name)
	if value != "" {
		return value, nil
	}
	if !fp.Required {
		if fp.DefaultValue == "" {
			return "", fmt.Errorf("%w %s", ErrEmptyVariable, fp.Name)
		}
		return fp.DefaultValue, nil
	}
	return "", fmt.Errorf("%w %s", ErrRequiredVariable, fp.Name)
}

func (fp *FieldParams) upcastValue(value string) (reflect.Value, error) {
	if fp.Type.Kind() == reflect.String {
		return reflect.ValueOf(value), nil
	}
	// cringe, but anyway.
	switch fp.Type.Kind() {
	case reflect.Int:
		a, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int(a)), nil
	case reflect.Uint:
		a, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint(a)), nil
	case reflect.Float32:
		a, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(float32(a)), nil
	case reflect.Float64:
		a, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(a), nil
	default:
		return reflect.Value{}, fmt.Errorf("%w %s", ErrUnsupportedType, fp.Type)
	}
}

func ParseField(field reflect.StructField) (FieldParams, error) {
	// Splitting params from tag
	tag := field.Tag.Get("config")
	params := strings.Split(tag, ",")

	// Result var
	var res FieldParams
	res.Type = field.Type
	// First item is name
	res.Name = params[0]

	// Parsing every param
	for _, param := range params[1:] {
		// Key value params
		if strings.Contains(param, "=") {
			splittedParam := strings.Split(param, "=")
			if len(splittedParam) != 2 {
				return FieldParams{}, fmt.Errorf("wrong format")
			}
			res.DefaultValue = splittedParam[1]
			continue
		}
		param := strings.TrimSpace(param)
		switch param {
		case "required":
			res.Required = true
		default:
			return FieldParams{}, fmt.Errorf("%w %s", ErrUnknownParam, param)
		}
	}
	return res, nil
}
