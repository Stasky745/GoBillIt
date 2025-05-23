package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"

	"github.com/Stasky745/go-libs/log"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/cliflagv2"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/urfave/cli/v2"
)

func loadConfig(cliContext *cli.Context) {
	configFile := cliContext.String("config")
	if configFile == "" {
		configFile = setDefaultEnv(string(APP_PREFIX)+"CONFIG", string(CONFIG_FILE))
	}

	// load base config file
	loadConfigFile(cliContext, configFile)

	configDir := cliContext.String("config_dir")
	if configDir == "" {
		configDir = setDefaultEnv(string(APP_PREFIX)+"CONFIG_DIR", string(CONFIG_DIR))
	}

	_, err := os.Stat(configDir)
	if err == nil {
		// configDir exists
		log.Info("config files directory exists. Processing files", "CONFIG_DIR", CONFIG_DIR)
		filepath.WalkDir(configDir, func(s string, d fs.DirEntry, e error) error {
			if e != nil {
				return e
			}
			if slices.Contains([]string{"json", "yaml", "yml"}, filepath.Ext(d.Name())) {
				loadConfigFile(cliContext, d.Name())
			}
			return nil
		})
	} else if errors.Is(err, fs.ErrNotExist) {
		// configDir doesn't exist
		log.Warn("config files directory does not exist", "CONFIG_DIR", CONFIG_DIR)
	} else {
		// other error
		log.Error("can't check config files directory", "error", err, "CONFIG_DIR", CONFIG_DIR)
	}

	validateRequiredFields()
}

func loadConfigFile(cliContext *cli.Context, configFile string) {
	var err error

	ext := filepath.Ext(configFile)

	switch ext {
	case ".json":
		err = k.Load(file.Provider(configFile), json.Parser())
	case ".yaml", ".yml":
		err = k.Load(file.Provider(configFile), yaml.Parser())
	}

	log.CheckErr(err, true, "error loading config file", "configFile", configFile)

	err = k.Load(env.ProviderWithValue(APP_PREFIX, ".", func(key, value string) (string, interface{}) {
		// Strip out the MYVAR_ prefix and lowercase and get the key while also replacing
		// the _ character with . in the key (koanf delimiter).
		newKey := strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(key, APP_PREFIX)), "_", ".")

		if slices.Contains(LIST_ENVS, key) {
			return newKey, strings.Split(value, ",")
		}

		return newKey, value
	}), nil)

	log.CheckErr(err, true, "can't load environment variables")

	err = k.Load(cliflagv2.Provider(cliContext, "-"), nil)
	log.CheckErr(err, true, "can't load flags")
}

// validateRequiredFields checks that all required keys are present in the config and validates different types of fields.
func validateRequiredFields() {
	var missingFields []string

	// Iterate over required fields and check if they're set
	for _, field := range REQUIRED_FIELDS {
		// Use Koanf to get the value for the field
		value := k.Get(field)

		// Check if the field is missing or invalid for different types
		switch value := value.(type) {
		case string:
			// A string is invalid if it's empty
			if value == "" {
				missingFields = append(missingFields, field)
			}
		case int:
			// An int is invalid if it's zero (could be any default invalid value)
			if value == 0 {
				missingFields = append(missingFields, field)
			}
		case bool:
			// A boolean is invalid if it's false (assuming false as default invalid value)
			if !value {
				missingFields = append(missingFields, field)
			}
		default:
			// For unsupported types, you can add more cases if necessary
			if value == nil || reflect.ValueOf(value).IsZero() {
				missingFields = append(missingFields, field)
			}
		}
	}

	if len(missingFields) > 0 {
		log.Panicf("❌ Missing required fields: %v", missingFields)
	}
}

func loadYaml(path string) error {
	return k.Load(file.Provider(path), yaml.Parser())
}
