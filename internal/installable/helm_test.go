package installable

import (
	"path/filepath"
	"testing"
)

func TestHelm_Install(t *testing.T) {
	type fields struct {
		URL     string
		Chart   string
		Version string
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
			name: "prometheus--latest",
			fields: fields{
				URL:   "https://prometheus-community.github.io/helm-charts",
				Chart: "prometheus",
			},
			wantErr: false,
		},
		{
			name: "grafana--versioned",
			fields: fields{
				URL:     "https://grafana.github.io/helm-charts",
				Chart:   "grafana",
				Version: "6.2.1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Helm{
				URL:     tt.fields.URL,
				Chart:   tt.fields.Chart,
				Version: tt.fields.Version,
			}
			if err := h.Install(); (err != nil) != tt.wantErr {
				t.Errorf("Helm.Install() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHelm_GetInstallPath(t *testing.T) {
	type fields struct {
		URL     string
		Chart   string
		Version string
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
			name: "prometheus--latest",
			fields: fields{
				URL:   "https://prometheus-community.github.io/helm-charts",
				Chart: "prometheus",
			},
			want:    filepath.Join(installDirName, "prometheus-community.github.io", "helm-charts", "prometheus", "latest"),
			wantErr: false,
		},
		{
			name: "grafana--versioned",
			fields: fields{
				URL:     "https://grafana.github.io/helm-charts",
				Chart:   "grafana",
				Version: "6.2.1",
			},
			want:    filepath.Join(installDirName, "grafana.github.io", "helm-charts", "grafana", "6.2.1"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Helm{
				URL:     tt.fields.URL,
				Chart:   tt.fields.Chart,
				Version: tt.fields.Version,
			}
			got, err := h.GetInstallPath()
			if (err != nil) != tt.wantErr {
				t.Errorf("Helm.GetInstallPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Helm.GetInstallPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHelm_Validate(t *testing.T) {
	type fields struct {
		URL     string
		Chart   string
		Version string
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
			name: "invalid: /wo chart /wo version",
			fields: fields{
				URL: "https://prometheus-community.github.io/helm-charts",
			},
			wantErr: true,
		},
		{
			name: "invalid: /wo chart /w version",
			fields: fields{
				URL:     "https://grafana.github.io/helm-charts",
				Version: "6.2.1",
			},
			wantErr: true,
		},
		{
			name: "valid: /w chart /wo version",
			fields: fields{
				URL:   "https://prometheus-community.github.io/helm-charts",
				Chart: "prometheus",
			},
			wantErr: false,
		},
		{
			name: "valid: /w chart /w version",
			fields: fields{
				URL:     "https://grafana.github.io/helm-charts",
				Chart:   "grafana",
				Version: "6.2.1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Helm{
				URL:     tt.fields.URL,
				Chart:   tt.fields.Chart,
				Version: tt.fields.Version,
			}
			if err := h.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Helm.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
