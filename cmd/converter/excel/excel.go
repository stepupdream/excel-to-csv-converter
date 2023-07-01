package excel

import (
	"os"
	"strings"

	"github.com/stepupdream/go-support-tool/directory"

	"github.com/pkg/errors"
	"github.com/stepupdream/go-support-tool/array"
	"github.com/stepupdream/go-support-tool/excel"
)

// GetFilePath Get the target Excel file path.
func GetFilePath(targetDirPath string) ([]string, error) {
	paths, err := excel.GetFilePathRecursive(targetDirPath)
	if err != nil {
		return nil, err
	}

	if err = verifyExcelPath(paths); err != nil {
		return nil, err
	}

	return paths, nil
}

// verifyExcelPath Verify the Excel file path.
func verifyExcelPath(excelFilePaths []string) error {
	pathSeparator := string(os.PathSeparator)
	outputFileExtension := os.Getenv("OutputFileExtension")
	productionDirectoryPath := os.Getenv("ProductionDirectoryPath") + pathSeparator + outputFileExtension
	productionVersion := directory.MaxFileName(productionDirectoryPath)
	developExcelDirectoryPath := os.Getenv("DevelopDirectoryPath") + pathSeparator + "excel"

	for _, excelFilePath := range excelFilePaths {
		var loadType string
		tmp := excelFilePath[len(developExcelDirectoryPath)+1:]

		if strings.Contains(excelFilePath, pathSeparator+"env"+pathSeparator) {
			loadType = strings.Split(tmp, pathSeparator)[2]
			if !array.Contains([]string{"insert", "update", "delete"}, loadType) {
				return errors.New("EditDirectory configuration is incorrect.")
			}
			continue
		}

		versionText := strings.Split(tmp, pathSeparator)[0]
		loadType = strings.Split(tmp, pathSeparator)[1]
		if !array.Contains([]string{"insert", "update", "delete"}, loadType) {
			return errors.New("EditDirectory configuration is incorrect.")
		}

		if versionText > productionVersion {
			continue
		}

		return errors.New("[WARNING] Conversion failed due to a lower version than the production version :" + excelFilePath)
	}

	return nil
}
