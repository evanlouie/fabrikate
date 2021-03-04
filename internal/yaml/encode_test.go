package yaml

import (
	"reflect"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	type args struct {
		docs []interface{}
	}

	// TODO figure out a good way to test maps with multiple entries.
	// String comparision does not work because ordering for map entries is not
	// ensured to be stable.
	// TODO figure out a how to trigger an encode error
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "zero",
			args:    args{},
			want:    nil,
			wantErr: false,
		},
		{
			name: "empty slice",
			args: args{
				[]interface{}{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "simple values",
			args: args{
				[]interface{}{
					1, "foo", "bar", true, false,
				},
			},
			want: []byte(`1
---
foo
---
bar
---
true
---
false`),
			wantErr: false,
		},
		{
			name: "complex values",
			args: args{
				[]interface{}{
					map[string]interface{}{
						"foo": map[string]interface{}{
							"bar": "baz",
							"list": []interface{}{
								1,
								2,
								"a string",
								map[string]interface{}{
									"nested": "map",
								},
							},
						},
					},
					[]interface{}{
						1, "foo", true,
					},
				},
			},
			want: []byte(`
foo:
  bar: baz
  list:
    - 1
    - 2
    - a string
    - nested: map
---
- 1
- foo
- true`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.docs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Do string comparison with all white space removed to make the test
			// simpler.
			// TODO see if we can make the tests precise enough to not need white
			// space removal
			gotStr := strings.ReplaceAll(string(got), " ", "")
			gotStr = strings.ReplaceAll(gotStr, "\n", "")
			wantStr := strings.ReplaceAll(string(tt.want), " ", "")
			wantStr = strings.ReplaceAll(wantStr, "\n", "")
			if gotStr != wantStr {
				t.Errorf("Encode() = %s, want %s", got, []byte(tt.want))
			}
		})
	}
}

func TestEncodeNoFailFast(t *testing.T) {
	type args struct {
		docs []interface{}
	}
	// TODO figure out how to trigger an encode error
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "zero",
			args:    args{},
			want:    nil,
			wantErr: false,
		},
		{
			name: "empty doc slice",
			args: args{
				docs: []interface{}{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "single valid value",
			args: args{
				docs: []interface{}{
					123,
				},
			},
			want:    []byte("123\n"),
			wantErr: false,
		},
		{
			name: "multiple valid value",
			args: args{
				docs: []interface{}{
					123,
					"foobar",
				},
			},
			want:    []byte("123\n---\nfoobar\n"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeNoFailFast(tt.args.docs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeNoFailFast() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				// print as %s instead of %v for human readability
				t.Errorf("EncodeNoFailFast() = %s, want %s", got, tt.want)
			}
		})
	}
}
