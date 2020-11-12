package generatable

import (
	"path/filepath"
	"testing"
)

func Test_pathable_GetGeneratePath(t *testing.T) {
	type fields struct {
		ComponentPath []string
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
			name:    "basic",
			fields:  fields{[]string{"foo", "bar_", "_baz", "_zaz_"}},
			want:    filepath.Join(generateDirName, "foo_bar___baz__zaz_") + ".yaml",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Pathable{
				ComponentPath: tt.fields.ComponentPath,
			}
			got, err := n.GetGeneratePath()
			if (err != nil) != tt.wantErr {
				t.Errorf("nestable.GetGeneratePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("nestable.GetGeneratePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanup(t *testing.T) {
	type args struct {
		g Generatable
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: true,
		},
		{
			name:    "non-existent-dir",
			args:    args{Static{Pathable: Pathable{ComponentPath: []string{"I", "don't", "exist"}}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cleanup(tt.args.g); (err != nil) != tt.wantErr {
				t.Errorf("cleanup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
