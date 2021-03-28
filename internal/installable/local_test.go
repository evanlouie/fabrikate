package installable

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var cwd string // gets set during init

func TestLocal_Install(t *testing.T) {
	type fields struct {
		Root string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "relative-directory",
			fields: fields{
				Root: "testdata/local",
			},
			wantErr: false,
		},
		{
			name: "relative-file",
			fields: fields{
				Root: "testdata/local/1-random-file.txt",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Local{
				Root: tt.fields.Root,
			}
			t.Cleanup(func() {
				cleanup(l)
			})
			if err := l.Install(); (err != nil) != tt.wantErr {
				t.Errorf("Local.Install() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocal_GetInstallPath(t *testing.T) {
	type fields struct {
		Root string
	}

	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:    "empty",
			fields:  fields{},
			want:    "",
			wantErr: true,
		},
		{
			name: "relative-file",
			fields: fields{
				Root: filepath.Join("testdata", "local", "nested", "3-random-file.txt"),
			},
			want:    filepath.Join(installDirName, localRoot, "testdata", "local", "nested"),
			wantErr: false,
		},
		{
			name: "relative",
			fields: fields{
				Root: filepath.Join("testdata", "local"),
			},
			want:    filepath.Join(installDirName, localRoot, "testdata", "local"),
			wantErr: false,
		},
		{
			name: "absolute",
			fields: fields{
				Root: filepath.Join(cwd, "testdata", "local"),
			},
			want:    filepath.Join(installDirName, localRoot, "testdata", "local"),
			wantErr: false,
		},
		{
			name: "relative-dots",
			fields: fields{
				Root: filepath.Join("testdata", "local", "does-not-exist", ".."),
			},
			want:    filepath.Join(installDirName, localRoot, "testdata", "local"),
			wantErr: false,
		},
		{
			name: "absolute-dots",
			fields: fields{
				Root: filepath.Join(cwd, "testdata", "local", "does-not-exist", ".."),
			},
			want:    filepath.Join(installDirName, localRoot, "testdata", "local"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Local{
				Root: tt.fields.Root,
			}
			got, err := l.GetInstallPath()
			if (err != nil) != tt.wantErr {
				t.Errorf("Local.GetInstallPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Local.GetInstallPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocal_Validate(t *testing.T) {
	type fields struct {
		Root string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "relative-exists",
			fields: fields{
				Root: filepath.Join("testdata", "local"),
			},
			wantErr: false,
		},
		{
			name: "absolute-exists",
			fields: fields{
				Root: filepath.Join(cwd, "testdata", "local"),
			},
			wantErr: false,
		},
		{
			name: "relative-dots-exists",
			fields: fields{
				Root: filepath.Join("testdata", "local", "does-not-exist", ".."),
			},
			wantErr: false,
		},
		{
			name: "absolute-dots-exists",
			fields: fields{
				Root: filepath.Join(cwd, "testdata", "local", "does-not-exist", ".."),
			},
			wantErr: false,
		},
		{
			name: "relative-not-exists",
			fields: fields{
				Root: filepath.Join("not", "exist", "local", "does-not-exist", ".."),
			},
			wantErr: true,
		},
		{
			name: "absolute-not-exists",
			fields: fields{
				Root: filepath.Join(cwd, "..", "some", "random", "dir", "foo", "..", "bar"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Local{
				Root: tt.fields.Root,
			}
			if err := l.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Local.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func init() {
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		log.Fatal(fmt.Errorf(`error computing current working directory for running local installable tests: %w`, err))
	}
}
