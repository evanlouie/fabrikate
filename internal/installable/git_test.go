package installable

import (
	"path/filepath"
	"testing"
)

func TestGit_Install(t *testing.T) {
	type fields struct {
		URL                 string
		SHA                 string
		Branch              string
		PersonalAccessToken string
	}
	const elasticGitRepo = "https://github.com/elastic/helm-charts"
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
				URL:    "https://github.com/elastic/helm-charts",
				Branch: "7.11",
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
		// Warn that certain tests require an internet connection clone
		if tt.fields.URL == elasticGitRepo {
			t.Logf(`WARNING: this test require an active internet connection to clone %s and may begin to fail if target branches/SHA are removed`, elasticGitRepo)
		}

		t.Run(tt.name, func(t *testing.T) {
			g := Git{
				URL:                 tt.fields.URL,
				SHA:                 tt.fields.SHA,
				Branch:              tt.fields.Branch,
				PersonalAccessToken: tt.fields.PersonalAccessToken,
			}
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
