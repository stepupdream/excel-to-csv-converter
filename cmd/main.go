package main

import (
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/stepupdream/excel-to-csv-converter/cmd/config"
	"github.com/stepupdream/excel-to-csv-converter/cmd/converter"
	"github.com/stepupdream/go-support-tool/directory"
	"github.com/stepupdream/go-support-tool/logger"
	"os"
	"strings"
)

func main() {
	// Load configuration file.
	logger.Setting("excel_to_csv_converter_error.log", true)
	if err := config.Load("excel_to_csv_converter_config.json"); err != nil {
		logger.Fatal(err)
	}

	// Excel ⇒ CSV
	color.HiYellow("[CONVERT START]")
	path, err := targetDirectoryPath()
	if err != nil {
		logger.Fatal(err)
	}
	_, err = converter.RunRoot(path)
	if err != nil {
		logger.Fatal(err)
	}
}

// targetDirectoryPath Get the target Excel directory path.
func targetDirectoryPath() (string, error) {
	currentDirectory, _ := os.Getwd()
	executionType := os.Getenv("ExecutionType")
	outputFileExtension := os.Getenv("OutputFileExtension")
	developDirectoryPath := os.Getenv("DevelopDirectoryPath")
	developExcelDirectoryPath := developDirectoryPath + string(os.PathSeparator) + "excel"
	developCsvDirectoryPath := developDirectoryPath + string(os.PathSeparator) + outputFileExtension

	var targetPath string
	switch executionType {
	case "pull", "PULL", "Pull":
		targetPath = strings.Replace(currentDirectory, developCsvDirectoryPath, developExcelDirectoryPath, 1)
	case "push", "Push", "PUSH":
		targetPath = currentDirectory
	default:
		return "", errors.New("the ExecutionType selection is incorrect, please select pull or push")
	}

	if !directory.Exist(targetPath) {
		return "", errors.New("the target directory does not exist")
	}

	fileInfo, err := os.Stat(targetPath)
	if err != nil {
		return "", err
	}
	if !fileInfo.IsDir() {
		return "", errors.New("incorrect directory specification. Please specify the directory")
	}

	return targetPath, nil
}
