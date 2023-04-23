package main

import (
	"github.com/stepupdream/go-support-tool/logger"
)

func main() {
	logger.Setting("excel_to_csv_converter_error.log", true)
	if err := load("excel_to_csv_converter_config.json"); err != nil {
		logger.Fatal(err)
	}

}
