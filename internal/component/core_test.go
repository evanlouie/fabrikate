package component

import (
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
			want: []string{
				"test-component",
				"test-component/prometheus",
				"test-component/grafana",
				"test-component/traefik",
				"test-component/my-svc",
				"test-component/azure voting",
				"test-component/traefik/traefik",
				"test-component/my-svc/random-svc",
			},
		},
		// {
		// 	name: "json",
		// 	args: args{filepath.Join("testdata", "install", "json")},
		// 	want: []string{"monitoring", "monitoring/prometheus", "monitoring/grafana"},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// cwd, err := os.Getwd()
			// if err != nil {
			// 	t.Fatal(err)
			// }
			// t.Log(cwd)
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
