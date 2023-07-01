package main

import (
	"os"
	"testing"
)

func Test_targetDirectoryPath(t *testing.T) {
	tests := []struct {
		name                 string
		ExecutionType        string
		OutputFileExtension  string
		DevelopDirectoryPath string
		CurrentDirectoryPath string
		want                 string
		wantErr              bool
	}{
		{
			name:                 "targetDirectoryPath",
			ExecutionType:        "pull",
			OutputFileExtension:  "csv",
			DevelopDirectoryPath: "converter/testdata/main",
			CurrentDirectoryPath: "/converter/testdata/main/csv",
			want:                 "/converter/testdata/main/excel",
			wantErr:              false,
		},
		{
			name:                 "targetDirectoryPath2",
			ExecutionType:        "push",
			OutputFileExtension:  "csv",
			DevelopDirectoryPath: "converter/testdata/main",
			CurrentDirectoryPath: "/converter/testdata/main/excel",
			want:                 "/converter/testdata/main/excel",
			wantErr:              false,
		},
	}
	currentDirectory, _ := os.Getwd()
	for _, tt := range tests {
		_ = os.Setenv("ExecutionType", tt.ExecutionType)
		_ = os.Setenv("OutputFileExtension", tt.OutputFileExtension)
		_ = os.Setenv("DevelopDirectoryPath", tt.DevelopDirectoryPath)
		_ = os.Chdir(currentDirectory + tt.CurrentDirectoryPath)

		t.Run(tt.name, func(t *testing.T) {
			got, err := targetDirectoryPath()
			if (err != nil) != tt.wantErr {
				t.Errorf("targetDirectoryPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			want := currentDirectory + tt.want
			if got != want {
				t.Errorf("targetDirectoryPath() got = %v, want %v", got, want)
			}
		})
	}
}
