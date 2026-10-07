package config

import "sync"

type Config interface {
	build()
	IsSet(key string) bool
	GetString(key string) string
	GetInt(key string) int
	GetFloat(key string) float64
	GetBool(key string) bool
	GetIntSlice(key string) []int
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringSlice(key string) []string
}

var once sync.Once

var (
	instance Config
)

func Default() Config {
	once.Do(func() {
		instance = ViperConfig{}
		instance.build()
	})
	return instance
}
