package installable

import "testing"

func Test_cleanup(t *testing.T) {
	type args struct {
		i Installable
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
			name: "non-existant",
			args: args{
				Local{
					Root: "some/random/directory/that/does/not/exist/1/23/2/123/12/31/23",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cleanup(tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("cleanup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
