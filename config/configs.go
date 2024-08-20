package config

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

var MapLevel = map[string]slog.Level{
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError,
}

var Confs Config

type Config struct {
	Server server `yaml:"server"`
	Logger logger `yaml:"logger"`
}

type server struct {
	Port uint `yaml:"port"`
}

type logger struct {
	AddSource bool   `yaml:"add_source"`
	Level     string `yaml:"level"`
}

func Load(path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(f, &Confs)
}
