package cmd

import "os"

var envVars = []string{
	"DEBUG",
}

type EnvVars struct {
	values map[string]string
}

func (ev *EnvVars) Get(key, defaultValue string) string {
	if ev.values[key] == "" {
		return defaultValue
	}
	return ev.values[key]
}

func InitEnvVars() *EnvVars {
	ev := &EnvVars{
		values: make(map[string]string),
	}
	for _, v := range envVars {
		ev.values[v] = os.Getenv(v)
	}
	return ev
}
