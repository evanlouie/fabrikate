package installable

import (
	"path/filepath"
	"testing"
)

func TestHTTP_Install(t *testing.T) {
	type fields struct {
		URL string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty",
			wantErr: true,
		},
		{
			name: "azure voting all in one",
			fields: fields{
				URL: "https://raw.githubusercontent.com/Azure-Samples/azure-voting-app-redis/master/azure-vote-all-in-one-redis.yaml",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HTTP{
				URL: tt.fields.URL,
			}
			if err := h.Install(); (err != nil) != tt.wantErr {
				t.Errorf("HTTP.Install() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHTTP_GetInstallPath(t *testing.T) {
	type fields struct {
		URL string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:    "empty",
			wantErr: true,
		},
		{
			name:   "valid",
			fields: fields{"https://raw.githubusercontent.com/Azure-Samples/azure-voting-app-redis/master/azure-vote-all-in-one-redis.yaml"},
			want:   filepath.Join(installDirName, "raw.githubusercontent.com", "Azure-Samples", "azure-voting-app-redis", "master", "azure-vote-all-in-one-redis.yaml"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HTTP{
				URL: tt.fields.URL,
			}
			got, err := h.GetInstallPath()
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTP.GetInstallPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HTTP.GetInstallPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTP_Validate(t *testing.T) {
	type fields struct {
		URL string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty",
			wantErr: true,
		},
		{
			name:    "non-http",
			fields:  fields{"ftp://foo@bar.com"},
			wantErr: true,
		},
		{
			name:    "http",
			fields:  fields{"http://raw.githubusercontent.com/Azure-Samples/azure-voting-app-redis/master/azure-vote-all-in-one-redis.yaml"},
			wantErr: false,
		},
		{
			name:    "https",
			fields:  fields{"https://raw.githubusercontent.com/Azure-Samples/azure-voting-app-redis/master/azure-vote-all-in-one-redis.yaml"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HTTP{
				URL: tt.fields.URL,
			}
			if err := h.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("HTTP.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
