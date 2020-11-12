package core

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestInstall(t *testing.T) {
	type args struct {
		startPath string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "yaml",
			args: args{filepath.Join("testdata", "install", "yaml")},
			want: []string{"monitoring", "monitoring/prometheus", "monitoring/grafana"},
		},
		// {
		// 	name: "json",
		// 	args: args{filepath.Join("testdata", "install", "json")},
		// 	want: []string{"monitoring", "monitoring/prometheus", "monitoring/grafana"},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cwd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			t.Log(cwd)
			got, err := Install(tt.args.startPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Install() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Install() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_install(t *testing.T) {
	type args struct {
		queue   []Component
		visited []Component
	}
	tests := []struct {
		name    string
		args    args
		want    []Component
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := install(tt.args.queue, tt.args.visited)
			if (err != nil) != tt.wantErr {
				t.Errorf("install() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("install() = %v, want %v", got, tt.want)
			}
		})
	}
}
