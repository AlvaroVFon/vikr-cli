// Package config handles the loading and management of application configuration settings.
package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type ScaffoldConfig struct {
	Type       string `mapstructure:"type"        validate:"required,oneof=api service lib cli"`
	Language   string `mapstructure:"language"    validate:"required,oneof=go ts"`
	OutputDir  string `mapstructure:"output_dir"  validate:"required"`
	IncludeGit bool   `mapstructure:"include_git"`
}

type Config struct {
	ProjectName string         `mapstructure:"project_name" validate:"required"`
	Version     string         `mapstructure:"version"      validate:"required"`
	Author      string         `mapstructure:"author"       validate:"required"`
	License     string         `mapstructure:"license"      validate:"required"`
	Debug       bool           `mapstructure:"debug"`
	Scaffold    ScaffoldConfig `mapstructure:"scaffold"     validate:"required"`
}

var defaultConfig = &Config{
	ProjectName: "MyVikrProject",
	Version:     "0.1.0",
	Author:      "Álvaro Villamarín",
	License:     "MIT",
	Debug:       false,
	Scaffold: ScaffoldConfig{
		Type:       "api",
		Language:   "go",
		OutputDir:  "./output",
		IncludeGit: true,
	},
}

var ConfigPaths = []string{
	".",
	".vikrrc",
	"vikr.yaml",
	"vikr.yml",
	"$HOME/.local/share/vikr/vikr.yaml",
	"$HOME/.local/share/vikr/vikr.yml",
	"/etc/vikr/vikr.yaml",
	"/etc/vikr/vikr.yml",
	"/usr/local/etc/vikr/vikr.yaml",
	"/usr/local/etc/vikr/vikr.yml",
}

var C Config

func Load() error {
	// yaml config
	if err := InitYAMLConfig(); err != nil {
		return err
	}

	// set defaults
	SetDefaults(defaultConfig)

	if err := viper.Unmarshal(&C); err != nil {
		return err
	}

	// validate config
	if err := Validate(&C); err != nil {
		return err
	}

	return nil
}

func InitYAMLConfig() error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("vikr")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/vikr")
	viper.AddConfigPath("/usr/local/etc/vikr")
	viper.SetConfigFile(".vikrrc")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}
	return nil
}

func SetDefaults(c *Config) {
	viper.SetDefault("project_name", c.ProjectName)
	viper.SetDefault("version", c.Version)
	viper.SetDefault("author", c.Author)
	viper.SetDefault("license", c.License)
	viper.SetDefault("debug", c.Debug)
	viper.SetDefault("scaffold.type", c.Scaffold.Type)
	viper.SetDefault("scaffold.language", c.Scaffold.Language)
	viper.SetDefault("scaffold.output_dir", c.Scaffold.OutputDir)
	viper.SetDefault("scaffold.include_git", c.Scaffold.IncludeGit)
}

func Validate(c *Config) error {
	v := validator.New()

	if err := v.Struct(c); err != nil {
		return err
	}

	return nil
}

func GenerateDefaultConfigYAML() (string, error) {
	_, exists := ConfigExists(ConfigPaths)
	if exists {
		return "", fmt.Errorf("ya existe un archivo de configuración Vikr en las rutas predeterminadas")
	}
	viper.SetConfigType("yaml")
	SetDefaults(defaultConfig)
	if err := viper.WriteConfigAs(".vkirrc"); err != nil {
		return "", err
	}
	return ".vikrrc", nil
}

// ConfigExists checks if any of the provided config file paths exist.
func ConfigExists(paths []string) (string, bool) {
	for _, path := range paths {
		path = os.ExpandEnv(path)

		if _, err := os.Stat(path); err == nil {
			return path, true // existe
		} else if !os.IsNotExist(err) {
			fmt.Printf("⚠️ Error comprobando %s: %v\n", path, err)
			return path, true
		}
	}
	return "", false // ninguno existe
}
