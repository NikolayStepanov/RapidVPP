package config

import (
	"os"
	"path/filepath"
)

const (
	defaultHttpHost = ""
	defaultHttpPort = "8080"

	defaultConfigDir            = "configs/"
	defaultLoggerConfigFileName = "/configs/logger.yml"
	defaultLoggerLevel          = "debug"
	defaultLoggerOutputPaths    = "stdout"
	dbConnEnvironment           = "DATA_SOURCE_NAME"
)

type (
	Config struct {
		HTTP   HTTPConfig
		Logger LoggerConfig
		VPP    VPPConfig
	}
	HTTPConfig struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}
	LoggerConfig struct {
		NameConfigFile string   `mapstructure:"name_config_file"`
		Level          string   `mapstructure:"level" yaml:"level"`
		OutputPaths    []string `mapstructure:"output_paths" yaml:"output_paths"`
		EncoderTime    string   `mapstructure:"encode_time" yaml:"encode_time"`
	}
	VPPConfig struct {
		Socket         string `yaml:"socket"`
		StreamPoolSize int    `yaml:"stream_pool_size"`
	}
)

func Init() *Config {
	cfg := Config{}
	cfg.getEnvironmentVariables()
	cfg.HTTP.Port = defaultHttpPort
	cfg.HTTP.Host = defaultHttpHost

	configDir := os.Getenv("CONFIG_DIR")
	if configDir == "" {
		configDir = defaultConfigDir
	}

	loggerConfigPath := filepath.Join(configDir, "logger.yml")
	cfg.Logger.NameConfigFile = loggerConfigPath
	cfg.Logger.Level = defaultLoggerLevel
	cfg.Logger.OutputPaths = []string{
		defaultLoggerOutputPaths,
	}
	cfg.VPP.Socket = "/run/vpp/api.sock"
	return &cfg
}

func (c *Config) getEnvironmentVariables() {
}
