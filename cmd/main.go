package main

import (
	"github.com/stepupdream/excel-to-csv-converter/cmd/config"
	"github.com/stepupdream/go-support-tool/logger"
)

func main() {
	// Load configuration file.
	logger.Setting("excel_to_csv_converter_error.log", true)
	if err := config.Load("excel_to_csv_converter_config.json"); err != nil {
		logger.Fatal(err)
	}
}
