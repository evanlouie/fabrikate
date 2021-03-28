package generatable

import (
	"path/filepath"
	"testing"
)

func TestStatic_Validate(t *testing.T) {
	type fields struct {
		ComponentPath []string
		ManifestPath  string
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
			name:    "empty-manifest-path",
			fields:  fields{[]string{"foo", "bar"}, ""},
			wantErr: true,
		},
		{
			name:    "empty-component-path",
			fields:  fields{nil, "./"},
			wantErr: true,
		},
		{
			name:    "basic",
			fields:  fields{[]string{"foo", "bar"}, filepath.Join("testdata", "static")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Static{
				Pathable: Pathable{
					ComponentPath: tt.fields.ComponentPath,
				},
				ManifestPath: tt.fields.ManifestPath,
			}
			if err := s.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Static.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStatic_Generate(t *testing.T) {
	type fields struct {
		ComponentPath []string
		ManifestPath  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name:    "empty",
			fields:  fields{},
			want:    0,
			wantErr: true,
		},
		{
			name: "basic",
			fields: fields{
				ComponentPath: []string{"foo", "bar"},
				ManifestPath:  filepath.Join("testdata", "static"),
			},
			want:    514,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Static{
				Pathable: Pathable{
					ComponentPath: tt.fields.ComponentPath,
				},
				ManifestPath: tt.fields.ManifestPath,
			}
			got, err := s.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Static.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Static.Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
