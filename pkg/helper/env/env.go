package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func GetString(key, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	return val
}

func GetStringSlice(key string, def []string) []string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	return strings.Split(val, " ")
}

func GetInt(key string, def int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return def
	}

	return i
}

func GetBool(key string, def bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	b, err := strconv.ParseBool(val)
	if err != nil {
		return def
	}

	return b
}

func GetDuration(key string, def time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	d, err := time.ParseDuration(val)
	if err != nil {
		return def
	}

	return d
}
