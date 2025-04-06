package utils

import (
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func ParseEnvConfig[T string | int | bool | []string](key string, defaultValue T) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	switch any(defaultValue).(type) {
	case string:
		return any(value).(T)
	case int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			logrus.Panic(err)
		}
		return any(intValue).(T)
	case bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			logrus.Panic(err)
		}
		return any(boolValue).(T)
	case []string:
		return any(ParseComma(value)).(T)
	default:
		panic("unknown type of defaultValue")
	}
}

func ParseComma(value string) []string {
	result := strings.FieldsFunc(value, func(r rune) bool {
		return r == ','
	})

	var parsed []string
	for _, item := range result {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			parsed = append(parsed, trimmed)
		}
	}

	return parsed
}

func ParseUrl(u string) url.URL {
	parsed, err := url.Parse(u)
	if err != nil {
		logrus.Fatal(err)
	}
	return *parsed
}

func GetPortFromAddress(addr string) string {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		logrus.Fatal(err)
	}
	return port
}
