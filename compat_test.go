// Copyright 2012 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package json

import (
	"reflect"
	"testing"
)

type TestCase struct {
	Raw    string
	Result interface{}
}

var cases = map[string]TestCase{
	"Number": TestCase{
		Raw: "1234",
		Result: int64(1234),
	},
	"Number - negative": TestCase{
		Raw: "-1234",
		Result: int64(-1234),
	},
	"Number - float": TestCase{
		Raw: "1234.5678",
		Result: float64(1234.5678),
	},
	"Number - negative float": TestCase{
		Raw: "-1234.5678",
		Result: float64(-1234.5678),
	},
	"String": TestCase{
		Raw: "\"foobar\"",
		Result: "foobar",
	},
	"String with encoded UTF-8": TestCase{
		Raw: "\"\\u6211\\u7231\\u4f60\"",
		Result: "我爱你",
	},
	"String with big-U encoded multibyte UTF-8": TestCase{
		Raw: "\"\\U0001D11E\"",
		Result: "𝄞",
	},
	"String with small-U encoded multibyte UTF-8": TestCase{
		Raw: "\"\\uD834\\uDD1E\"",
		Result: "𝄞",
	},
	"String with hex encoded multibyte UTF-8": TestCase{
		Raw: "\"\\xF0\\x9D\\x84\\x9E\"",
		Result: "𝄞",
	},
	"String with encoded UTF-8 and backslash": TestCase{
		Raw: "\"10\\\\10 ~ \\u2764\"",
		Result: "10\\10 ~ ❤",
	},
	"String with backslash": TestCase{
		Raw: "\"10\\\\10\"",
		Result: "10\\10",
	},
	"String with backslash and tab": TestCase{
		Raw: "\"10\\\\\t10\"",
		Result: "10\\	10",
	},
	"String with backslash and backspace": TestCase{
		Raw: "\"10\\\\\b10\"",
		Result: "10\\\b10",
	},
	"String with escaped forward slash": TestCase{
		Raw: "\"\\\\\\/\"",
		Result: "\\/",
	},
	"Object": TestCase{
		Raw: "{\"foo\":\"bar\"}",
		Result: map[string]interface{}{
			"foo": "bar",
		},
	},
	"Object with spaces": TestCase{
		Raw: "{ \"foo\" : \"bar\" }",
		Result: map[string]interface{}{
			"foo": "bar",
		},
	},
	"Object with UTF-8 value": TestCase{
		Raw: "{ \"foo\" : \"\\u6211\" }",
		Result: map[string]interface{}{
			"foo": "我",
		},
	},
	"Object with tabs": TestCase{
		Raw: "{	\"foo\"	:	\"bar\"	}",
		Result: map[string]interface{}{
			"foo": "bar",
		},
	},
	"Object with empty nested object": TestCase{
		Raw: "{ \"foo\": {}}",
		Result: map[string]interface{}{
			"foo": map[string]interface{}{},
		},
	},
	"Object with empty nested array": TestCase{
		Raw: "{\"foo\": []}",
		Result: map[string]interface{}{
			"foo": []interface{}{},
		},
	},
	"Array": TestCase{
		Raw: "[1234,\"foobar\"]",
		Result: []interface{}{
			int64(1234),
			"foobar",
		},
	},
	"Array with spaces": TestCase{
		Raw: "[ 1234 , \"foobar\" ]",
		Result: []interface{}{
			int64(1234),
			"foobar",
		},
	},
	"Array with tabs": TestCase{
		Raw: "[	1234	,	\"foobar\"	]",
		Result: []interface{}{
			int64(1234),
			"foobar",
		},
	},
	"Array with multiple tabs": TestCase{
		Raw: "[				1234,\"foobar\"]",
		Result: []interface{}{
			int64(1234),
			"foobar",
		},
	},
	"Array with no contents": TestCase{
		Raw: "[]",
		Result: []interface{}{},
	},
	"Array with empty object": TestCase{
		Raw: "[{}]",
		Result: []interface{}{
			map[string]interface{}{},
		},
	},
}

func TestCases(t *testing.T) {
	var (
		err    error
		decode interface{}
	)
	for desc, testcase := range cases {
		if err = Unmarshal([]byte(testcase.Raw), &decode); err != nil {
			t.Fatalf("Error decoding '%v': %v", desc, err)
		}
		if !reflect.DeepEqual(decode, testcase.Result) {
			t.Logf("%v\n", reflect.TypeOf(decode))
			t.Logf("%v\n", reflect.TypeOf(testcase.Result))
			if reflect.TypeOf(decode) == reflect.TypeOf("") {
				t.Logf("Decode: %v\n", []byte(decode.(string)))
				t.Logf("Expected: %v\n",
					[]byte(testcase.Result.(string)))
			}
			t.Fatalf("Problem decoding '%v' Expected: %v, Got %v",
				desc, testcase.Result, decode)
		}
	}
}
