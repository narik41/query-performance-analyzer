package config

import "testing"

func TestConfigInit(t *testing.T) {
	Default().GetString("test")
}
