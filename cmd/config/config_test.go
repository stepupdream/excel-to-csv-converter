package config

import "testing"

func TestLoad(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Load1",
			args: args{
				fileName: "test_config.json",
			},
			wantErr: false,
		},
		{
			name: "Load2",
			args: args{
				fileName: "test_config2.json",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Load(tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
