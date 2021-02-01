package env_utils

import "os"

func GetEnvVarWithDefaults(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
