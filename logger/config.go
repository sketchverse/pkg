package logger

import (
	"os"
	"path/filepath"
)

type Config struct {
	Output   string `mapstructure:"output"`
	LogFile  string `mapstructure:"log_file"`
	Level    string `mapstructure:"level"`
	MaxSize  int    `mapstructure:"max_size"`
	MaxAge   int    `mapstructure:"max_age"`
	MaxFiles int    `mapstructure:"max_files"`
	Compress bool   `mapstructure:"compress"`
}

func DefaultConfig() *Config {
	return &Config{
		Level:    "info",
		Output:   "stdout",
		LogFile:  filepath.Join(os.TempDir(), "app.log"),
		MaxSize:  100,
		MaxFiles: 3,
		MaxAge:   7,
		Compress: true,
	}
}
