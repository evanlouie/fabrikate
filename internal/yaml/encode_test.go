package yaml

import (
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
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"zero",
			args{},
			"",
			false,
		},
		{
			"empty slice",
			args{
				[]interface{}{},
			},
			"",
			false,
		},
		{
			"simple values",
			args{
				[]interface{}{
					1, "foo", "bar", true, false,
				},
			},
			`1
---
foo
---
bar
---
true
---
false`,
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
			`
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
				t.Errorf("Encode() = %s, want %s", got, []byte(tt.want))
			}
		})
	}
}
