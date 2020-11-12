package installable

import (
	"path/filepath"
	"testing"
)

func TestGit_Install(t *testing.T) {
	type fields struct {
		URL    string
		SHA    string
		Branch string
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
			name: "kibana--versioned-sha",
			fields: fields{
				URL: "https://github.com/elastic/helm-charts",
				SHA: "3fb0c8267e146ef9ae8d8de7f836bb775c03e960",
			},
			wantErr: false,
		},
		{
			name: "kibana--versioned-branch",
			fields: fields{
				URL: "https://github.com/elastic/helm-charts",
				SHA: "master",
			},
			wantErr: false,
		},
		{
			name: "kibana-latest",
			fields: fields{
				URL: "https://github.com/elastic/helm-charts",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Git{
				URL:    tt.fields.URL,
				SHA:    tt.fields.SHA,
				Branch: tt.fields.Branch,
			}
			t.Cleanup(func() {
				cleanup(g)
			})
			if err := g.Install(); (err != nil) != tt.wantErr {
				t.Errorf("Git.Install() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGit_GetInstallPath(t *testing.T) {
	type fields struct {
		URL    string
		SHA    string
		Branch string
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
			name: "kibana--versioned-sha",
			fields: fields{
				URL: "https://github.com/elastic/helm-charts",
				SHA: "3fb0c8267e146ef9ae8d8de7f836bb775c03e960",
			},
			want:    filepath.Join(installDirName, "github.com", "elastic", "helm-charts", "3fb0c8267e146ef9ae8d8de7f836bb775c03e960"),
			wantErr: false,
		},
		{
			name: "kibana--versioned-branch",
			fields: fields{
				URL: "https://github.com/elastic/helm-charts",
				SHA: "master",
			},
			want:    filepath.Join(installDirName, "github.com", "elastic", "helm-charts", "master"),
			wantErr: false,
		},
		{
			name: "kibana-latest",
			fields: fields{
				URL: "https://github.com/elastic/helm-charts",
			},
			want:    filepath.Join(installDirName, "github.com", "elastic", "helm-charts", "latest"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Git{
				URL:    tt.fields.URL,
				SHA:    tt.fields.SHA,
				Branch: tt.fields.Branch,
			}
			got, err := g.GetInstallPath()
			if (err != nil) != tt.wantErr {
				t.Errorf("Git.GetInstallPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Git.GetInstallPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGit_Validate(t *testing.T) {
	type fields struct {
		URL    string
		SHA    string
		Branch string
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
			name: "no-url",
			fields: fields{
				URL: "",
				SHA: "foobar",
			},
			wantErr: true,
		},
		{
			name: "kibana--versioned-sha",
			fields: fields{
				URL: "https://github.com/elastic/helm-charts",
				SHA: "3fb0c8267e146ef9ae8d8de7f836bb775c03e960",
			},
			wantErr: false,
		},
		{
			name: "kibana--versioned-branch",
			fields: fields{
				URL: "https://github.com/elastic/helm-charts",
				SHA: "master",
			},
			wantErr: false,
		},
		{
			name: "kibana-latest",
			fields: fields{
				URL: "https://github.com/elastic/helm-charts",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Git{
				URL:    tt.fields.URL,
				SHA:    tt.fields.SHA,
				Branch: tt.fields.Branch,
			}
			if err := g.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Git.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_clone(t *testing.T) {
	type args struct {
		repo string
		dir  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := clone(tt.args.repo, tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("clone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkout(t *testing.T) {
	type args struct {
		repo   string
		target string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkout(tt.args.repo, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("checkout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
