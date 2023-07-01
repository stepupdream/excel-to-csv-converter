package converter

import (
	"github.com/stepupdream/go-support-tool/delimited"
	"os"
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Run1",
			args: args{
				filePath: "/testdata/convert/excel/1_0_0_0/insert/tests.xlsx",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Run2",
			args: args{
				filePath: "/testdata/convert/excel/1_0_0_0/update/users.xlsx",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Run3",
			args: args{
				filePath: "/testdata/convert/excel/1_0_0_0/insert/tests2.xlsx",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Run4",
			args: args{
				filePath: "/testdata/convert/excel/1_0_0_0/insert/tests3.xlsx",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Run5",
			args: args{
				filePath: "/testdata/convert/excel/1_0_0_0/insert/tests4.xlsx",
			},
			want:    false,
			wantErr: false,
		},
	}
	currentDirectory, _ := os.Getwd()
	for _, tt := range tests {
		_ = os.Setenv("OutputFileExtension", "csv")
		_ = os.Setenv("DevelopDirectoryPath", currentDirectory+"/testdata/convert")
		_ = os.Setenv("SheetName", "data")

		t.Run(tt.name, func(t *testing.T) {
			got, err := Run(currentDirectory + tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Run() got = %v, want %v", got, tt.want)
			}
		})
	}

	fileNames := []string{"tests.csv", "tests3.csv", "tests4.csv"}
	for _, fileName := range fileNames {
		dataA, _ := delimited.Load(currentDirectory+"/testdata/convert/csv/1_0_0_0/insert/"+fileName, false, false)
		dataB, _ := delimited.Load(currentDirectory+"/testdata/convert/diff/"+fileName, false, false)
		if !reflect.DeepEqual(dataA, dataB) {
			t.Errorf("Diff Error: %v", fileName)
		}
	}
}
