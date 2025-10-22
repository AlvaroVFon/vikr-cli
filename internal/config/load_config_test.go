package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSetDefaults(t *testing.T) {
	testCases := []struct {
		name           string
		config         *Config
		expectedConfig *Config
	}{
		{
			name: "default config",
			config: &Config{
				ProjectName: "TestProject",
				Version:     "1.0.0",
				Author:      "Test Author",
				License:     "MIT",
				Debug:       true,
				Scaffold: ScaffoldConfig{
					Type:       "api",
					Language:   "go",
					OutputDir:  "output",
					IncludeGit: true,
				},
			},
			expectedConfig: &Config{
				ProjectName: "TestProject",
				Version:     "1.0.0",
				Author:      "Test Author",
				License:     "MIT",
				Debug:       true,
				Scaffold: ScaffoldConfig{
					Type:       "api",
					Language:   "go",
					OutputDir:  "output",
					IncludeGit: true,
				},
			},
		},
	}

	for _, tc := range testCases {
		viper.Reset()
		t.Run(tc.name, func(t *testing.T) {
			SetDefaults(tc.config)
			assert.Equal(t, tc.expectedConfig, tc.config)
		})
	}
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		name      string
		config    *Config
		expectErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				ProjectName: "TestProject",
				Version:     "1.0.0",
				Author:      "Test Author",
				License:     "MIT",
				Debug:       true,
				Scaffold: ScaffoldConfig{
					Type:       "api",
					Language:   "go",
					OutputDir:  "output",
					IncludeGit: true,
				},
			},
			expectErr: false,
		},
		{
			name: "missing project name",
			config: &Config{
				Version: "1.0.0",
				Author:  "Test Author",
				License: "MIT",
				Scaffold: ScaffoldConfig{
					Type:     "api",
					Language: "go",
				},
			},
			expectErr: true,
		},
		{
			name: "missing scaffold type",
			config: &Config{
				ProjectName: "TestProject",
				Version:     "1.0.0",
				Author:      "Test Author",
				License:     "MIT",
				Scaffold: ScaffoldConfig{
					Language: "go",
				},
			},
			expectErr: true,
		},
		{
			name: "missing scaffold language",
			config: &Config{
				ProjectName: "TestProject",
				Version:     "1.0.0",
				Author:      "Test Author",
				License:     "MIT",
				Scaffold: ScaffoldConfig{
					Type: "api",
				},
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(tc.config)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
