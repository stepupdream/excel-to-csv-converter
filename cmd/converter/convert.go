package converter

import (
	"github.com/stepupdream/excel-to-csv-converter/cmd/converter/excel"
	"github.com/stepupdream/go-support-tool/array"
	"github.com/stepupdream/go-support-tool/console"
	"github.com/stepupdream/go-support-tool/delimited"
	"os"
	"path/filepath"
	"reflect"

	"github.com/xuri/excelize/v2"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/stepupdream/go-support-tool/file"
)

// RunRoot is the main function of the converter.
func RunRoot(targetDirectory string) (changedFilePaths []string, err error) {
	color.HiYellow("[CONVERT START]")
	excelFilePaths, err := excel.GetFilePath(targetDirectory)
	progressBar := console.StartProgressBar(len(excelFilePaths))
	for _, filePath := range excelFilePaths {
		isSkip, e := Run(filePath)
		if e != nil {
			return nil, e
		}
		if isSkip {
			continue
		}
		changedFilePaths = append(changedFilePaths, filePath)
		color.Green("[COMPLETE] " + filePath)
		progressBar.Increment()
	}
	progressBar.Finish()

	return changedFilePaths, nil
}

// Run is the main function of the converter.
// If there is a progress bar, the test will not work properly, so the functions are separated.
func Run(filePath string) (bool, error) {
	outputPath, e := prepare(filePath)
	if e != nil {
		return false, e
	}
	convertedRows, e := convert(filePath)
	if e != nil {
		return false, e
	}
	isSkip, e := createFile(outputPath, convertedRows)
	if e != nil {
		return false, e
	}

	return isSkip, nil
}

// createFile Create a file.
func createFile(outputPath string, convertedRows [][]string) (bool, error) {
	if file.Exists(outputPath) {
		existingRows, err := delimited.Load(outputPath, false, false)
		if err != nil {
			return true, err
		}
		if reflect.DeepEqual(convertedRows, existingRows) {
			color.Green("[  SKIP  ] " + outputPath)
			return true, nil
		}
	}

	// Embed special words to determine if the process was completed correctly.
	// Remove special characters when the process is completely finished.
	// Inserting the special character prevents the conversion from being skipped at an unintended timing.
	columnCount := len(convertedRows[0])
	var validationText []string
	for i := 0; i < columnCount; i++ {
		validationText = append(validationText, "#unverified#")
	}
	beforeValidationRows := append(convertedRows, validationText)
	if err := delimited.CreateNewFile(outputPath, beforeValidationRows); err != nil {
		return false, err
	}

	return false, nil
}

// prepare Prepare the output file.
func prepare(filePath string) (string, error) {
	outputFileExtension := os.Getenv("OutputFileExtension")
	delimitedDirectoryPath := os.Getenv("DevelopDirectoryPath") + string(os.PathSeparator) + outputFileExtension
	excelDirectoryPath := os.Getenv("DevelopDirectoryPath") + string(os.PathSeparator) + "excel"

	directoryPath := filepath.Dir(filePath[len(excelDirectoryPath):])
	filePathWithoutExtension := file.BaseFileName(filePath)
	if err := os.MkdirAll(delimitedDirectoryPath+directoryPath, 0755); err != nil {
		return "", err
	}

	pathSeparator := string(os.PathSeparator)
	outputPath := delimitedDirectoryPath + directoryPath + pathSeparator + filePathWithoutExtension + "." + outputFileExtension

	return outputPath, nil
}

// convert Convert Excel.
func convert(excelPath string) (result [][]string, err error) {
	excelFile, err := excelize.OpenFile(excelPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := excelFile.Close(); err == nil {
			err = closeErr
		}
	}()

	sheetName := os.Getenv("SheetName")
	rowsExcel, err := excelFile.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rowsExcel) == 0 {
		return nil, errors.New("Excel with empty data. : " + excelPath)
	}

	if err = verifyColumnUnique(rowsExcel[0]); err != nil {
		return nil, err
	}

	result = createNewRows(rowsExcel)

	return result, nil
}

// verifyColumnUnique Verify that the column name is unique.
func verifyColumnUnique(row []string) error {
	if !array.IsUnique(row) {
		return errors.New("The column names in the first row must be unique.")
	}

	return nil
}

// createNewRows Create a new row.
// In data read from Excel, if the last column is blank, the data is not read correctly and should be corrected.
func createNewRows(rowsExcel [][]string) (newRowsExcel [][]string) {
	columnCount := len(rowsExcel[0])

	for _, row := range rowsExcel {
		diffCount := columnCount - len(row)
		if diffCount != 0 {
			for i := 0; i < diffCount; i++ {
				row = append(row, "")
			}
		}
		newRowsExcel = append(newRowsExcel, row)
	}
	return newRowsExcel
}
