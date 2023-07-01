package config

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/stepupdream/go-support-tool/file"
)

// Setting https://mholt.github.io/json-to-go/
type Setting struct {
	ExecutionType               string `json:"execution_type"`
	ProductionDirectoryPath     string `json:"production_directory_path"`
	DevelopDirectoryPath        string `json:"develop_directory_path"`
	OutputFileExtension         string `json:"output_file_extension"`
	SheetName                   string `json:"sheet_name"`
	EnvName                     string `json:"env_name"`
	LastUpdatedTimeFilePath     string `json:"last_updated_time_file_path"`
	MasterDataYamlDirectoryPath string `json:"master_data_yaml_directory_path"`
	EnumYamlDirectoryPath       string `json:"enum_yaml_directory_path"`
	TargetVersion               string `json:"target_version"`
	NeedValidation              bool   `json:"need_validation"`
}

// Load
// read all data from a file at once.
// Read and return in []bytes, so Close is unnecessary.
func Load(fileName string) error {
	path, err := file.RecursiveFilePathInParent(fileName)
	if err != nil {
		return err
	}
	config, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err = verifyKey(config); err != nil {
		return err
	}

	var setting Setting
	if err = json.Unmarshal(config, &setting); err != nil {
		return err
	}

	// Set everything specified in the structure to environment variables
	reflection := reflect.TypeOf(Setting{})
	for i := 0; i < reflection.NumField(); i++ {
		field := reflection.Field(i)
		value := reflect.Indirect(reflect.ValueOf(setting)).FieldByName(field.Name)
		if err = setEnv(field.Name, field.Type.String(), value); err != nil {
			return err
		}
	}

	return nil
}

// verifyKey check if the key is missing
func verifyKey(config []byte) error {
	var raw json.RawMessage
	if err := json.Unmarshal(config, &raw); err != nil {
		return err
	}

	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &m); err != nil {
		return err
	}

	reflection := reflect.TypeOf(Setting{})
	for i := 0; i < reflection.NumField(); i++ {
		field := reflection.Field(i)
		if _, ok := m[toSnakeCase(field.Name)]; !ok {
			return errors.New("missing key: " + field.Name)
		}
	}

	return nil
}

// setEnv set environment variable
func setEnv(fieldName string, fieldType string, value reflect.Value) error {
	var strValue string
	switch fieldType {
	case "string":
		strValue = value.String()
	case "bool":
		strValue = strconv.FormatBool(value.Bool())
	default:
		return errors.New("setEnv:not supported type")
	}

	return os.Setenv(fieldName, strValue)
}

// toSnakeCase convert to a snake case
func toSnakeCase(str string) string {
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchAllCap.ReplaceAllString(str, "${1}_${2}")

	return strings.ToLower(snake)
}
