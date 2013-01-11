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
	"Simple object": TestCase{
		Raw: "{\"foo\":\"bar\"}",
		Result: map[string]interface{}{
			"foo": "bar",
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
			t.Fatalf("Error decoding '%v'", desc)
		}
		if !reflect.DeepEqual(decode, testcase.Result) {
			t.Fatalf(
				"Problem decoding '%v' Expected: %v, Got %v",
				desc, testcase.Result, decode)
		}
	}
}
