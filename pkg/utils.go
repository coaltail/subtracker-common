package pkg

import "os"

func GetEnvOrDefault(env string, def string) string {
	env, ok := os.LookupEnv(env)
	if !ok {
		return def
	}
	return env
}
