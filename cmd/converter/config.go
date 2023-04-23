package main

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strconv"

	"github.com/stepupdream/go-support-tool/file"
)

// ConfigJson https://mholt.github.io/json-to-go/
type ConfigJson struct {
	ExecutionType               string `json:"execution_type"`
	ProductionCsvDirectoryPath  string `json:"production_csv_directory_path"`
	DevelopCsvDirectoryPath     string `json:"develop_csv_directory_path"`
	DevelopExcelDirectoryPath   string `json:"develop_excel_directory_path"`
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
func load(fileName string) error {
	path, err := file.RecursiveFilePathInParent(fileName)
	if err != nil {
		return err
	}
	config, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var configJson ConfigJson
	if err = json.Unmarshal(config, &configJson); err != nil {
		return err
	}

	reflection := reflect.TypeOf(ConfigJson{})
	for i := 0; i < reflection.NumField(); i++ {
		field := reflection.Field(i)
		value := reflect.Indirect(reflect.ValueOf(configJson)).FieldByName(field.Name)
		if err = setEnv(field.Name, field.Type.String(), value); err != nil {
			return err
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
