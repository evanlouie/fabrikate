package yaml

import (
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	type args struct {
		docs []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"empty",
			args{},
			"",
			false,
		},
		{
			"no values",
			args{
				[]interface{}{},
			},
			"",
			false,
		},
		{
			"some values",
			args{
				[]interface{}{
					1, "foo", "bar",
				},
			},
			`1
---
foo
---
bar`,
			false,
		},
		{
			"complex values",
			args{
				[]interface{}{
					map[string]interface{}{
						"foo": map[string]interface{}{
							"bar": "baz",
							"list": []interface{}{
								1,
								2,
								"a string",
							},
						},
					},
					[]interface{}{
						1, "foo", true,
					},
				},
			},
			`
foo:
  bar: baz
  list:
    - 1
    - 2
    - a string
---
- 1
- foo
- true`,
			false,
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
				// print as %s instead of %v so the error is human readable
				t.Errorf("Encode() = %s, want %s", got, []byte(tt.want))
			}
		})
	}
}
