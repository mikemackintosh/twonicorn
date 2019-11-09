package config

import (
	"fmt"
	"reflect"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

var testValidateEntry = []struct {
	got  []byte
	want error
}{
	{
		got: []byte(`
target: test
owner: owner@email.com
url: https://localhost
method: POST
translate:
  field: test
`),
		want: nil,
	},
	{
		got: []byte(`
target: test
url: https://localhost
method: POST
owner: owner@email.com
`),
		want: fmt.Errorf("invalid translate provided"),
	},
	{
		got: []byte(`
target: test
method: POST
owner: owner@email.com
`),
		want: fmt.Errorf("invalid url provided"),
	},
}

func TestValidateEntry(t *testing.T) {
	for _, test := range testValidateEntry {
		var e Entry
		err := yaml.Unmarshal(test.got, &e)
		if err != nil {
			t.Fatal(err)
		}

		var validateResult = e.Validate()
		if !reflect.DeepEqual(validateResult, test.want) {
			t.Fatalf("got '%s' (%T), wanted '%s' (%T)", validateResult, validateResult, test.want, test.want)
		}
	}
}
