package excel

import (
	"os"
	"reflect"
	"testing"
)

func TestGetFilePath(t *testing.T) {
	type args struct {
		targetDirPath string
	}
	tests := []struct {
		name                    string
		ProductionDirectoryPath string
		DevelopDirectoryPath    string
		OutputFileExtension     string
		args                    args
		want                    []string
		wantErr                 bool
	}{
		{
			name:                    "GetFilePath1",
			ProductionDirectoryPath: "./testdata/pattern1/production",
			DevelopDirectoryPath:    "./testdata/pattern1/develop",
			OutputFileExtension:     "csv",
			args: args{
				targetDirPath: "./testdata/pattern1/develop/excel/1_0_2_0",
			},
			want:    []string{"testdata/pattern1/develop/excel/1_0_2_0/insert/tests.xlsx"},
			wantErr: false,
		},
		{
			name:                    "GetFilePath2",
			ProductionDirectoryPath: "./testdata/pattern1/production",
			DevelopDirectoryPath:    "./testdata/pattern1/develop",
			OutputFileExtension:     "csv",
			args: args{
				targetDirPath: "./testdata/pattern1/develop/excel/1_0_3_0",
			},
			want: []string{
				"testdata/pattern1/develop/excel/1_0_3_0/insert/tests.xlsx",
				"testdata/pattern1/develop/excel/1_0_3_0/insert/tests2.xlsx",
			},
			wantErr: false,
		},
		{
			name:                    "GetFilePath3",
			ProductionDirectoryPath: "./testdata/pattern1/production",
			DevelopDirectoryPath:    "./testdata/pattern1/develop",
			OutputFileExtension:     "csv",
			args: args{
				targetDirPath: "./testdata/pattern1/develop/excel/env/Tom",
			},
			want: []string{
				"testdata/pattern1/develop/excel/env/Tom/insert/tests.xlsx",
			},
			wantErr: false,
		},
		{
			name:                    "GetFilePath4",
			ProductionDirectoryPath: "./testdata/pattern1/production",
			DevelopDirectoryPath:    "./testdata/pattern1/develop",
			OutputFileExtension:     "csv",
			args: args{
				targetDirPath: "./testdata/pattern1/develop/excel/1_0_4_0",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:                    "GetFilePath5",
			ProductionDirectoryPath: "./testdata/pattern1/production",
			DevelopDirectoryPath:    "./testdata/pattern1/develop",
			OutputFileExtension:     "csv",
			args: args{
				targetDirPath: "./testdata/pattern1/develop/excel/1_0_1_0",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		_ = os.Setenv("ProductionDirectoryPath", tt.ProductionDirectoryPath)
		_ = os.Setenv("DevelopDirectoryPath", tt.DevelopDirectoryPath)
		_ = os.Setenv("OutputFileExtension", tt.OutputFileExtension)

		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFilePath(tt.args.targetDirPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFilePath() got = %v, want %v", got, tt.want)
			}
		})
	}
}
