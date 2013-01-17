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
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"strconv"
)

const (
	STRING = iota
	NUMBER
	MAP
	ARRAY
	ENDARRAY
	ESCAPE
	BOOL
	NULL
)

type Event struct {
	Type  int
	Index int
}

type State struct {
	data   []byte
	i      int
	v      interface{}
	events []Event
}

func (s *State) Read() (err error) {
	var t int = s.nextType()
	switch t {
	case STRING:
		err = s.readString()
	case NUMBER:
		err = s.readNumber()
	case MAP:
		err = s.readMap()
	case ARRAY:
		err = s.readArray()
	case ENDARRAY:
		s.i++
		err = EndArray{}
	case BOOL:
		err = s.readBool()
	case NULL:
		err = s.readNull()
	case ESCAPE:
		err = fmt.Errorf("JSON should not start with escape")
	default:
		b := string(s.data[s.i-10 : s.i])
		c := string(s.data[s.i : s.i+1])
		e := string(s.data[s.i+1 : s.i+10])
		err = fmt.Errorf("Unrecognized type in %v -->%v<-- %v", b, c, e)
	}
	return
}

func (s *State) nextType() int {
	for {
		c := s.data[s.i]
		switch {
		case c == ' ':
			fallthrough
		case c == '\t':
			s.i++
			break
		case c == '"':
			return STRING
		case '0' <= c && c <= '9' || c == '-':
			return NUMBER
		case c == '[':
			return ARRAY
		case c == ']':
			return ENDARRAY
		case c == '{':
			return MAP
		case c == 't' || c == 'T' || c == 'f' || c == 'F':
			return BOOL
		case c == 'n':
			return NULL
		}
	}
	return -1
}

func (s *State) readString() (err error) {
	var (
		c       byte
		start   int
		buf     *bytes.Buffer
		atstart bool = false
		more    bool = true
		utf     bool = false
	)
	for atstart == false {
		c = s.data[s.i]
		switch {
		case c == ' ':
			fallthrough
		case c == '\t':
			s.i++
		case c == '"':
			atstart = true
			break
		case c == '}':
			s.i++
			return EndMap{}
		case c == ']':
			s.i++
			return EndArray{}
		}
	}
	s.i++
	start = s.i
	buf = new(bytes.Buffer)
	for more {
		c = s.data[s.i]
		switch {
		case c == '\\':
			buf.Write(s.data[start:s.i])
			switch {
			case len(s.data) > s.i+8 && s.data[s.i+1] == 'U':
				fallthrough
			case len(s.data) > s.i+6 && s.data[s.i+1] == 'u':
				fallthrough
			case len(s.data) > s.i+4 && s.data[s.i+1] == 'x':
				utf = true
				buf.WriteString("\\")
			}
			s.i++
			start = s.i
		case c == '"':
			more = false
		case s.i >= len(s.data)-1:
			return fmt.Errorf("No string terminator")
		}
		s.i++
	}
	buf.Write(s.data[start : s.i-1])
	s.v = buf.String()
	if utf == true {
		s.v, err = strconv.Unquote(fmt.Sprintf("\"%v\"", s.v))
	}
	return
}

func (s *State) readNumber() (err error) {
	var c byte
	var val int64 = 0
	var valf float64 = 0
	var mult int64 = 1
	if s.data[s.i] == '-' {
		mult = -1
		s.i++
	}
	var more = true
	var places int = 0
	for more {
		c = s.data[s.i]
		switch {
		case '0' <= c && c <= '9':
			if places != 0 {
				places *= 10
			}
			val = val*10 + int64(c-'0')
		case '}' == c:
			err = EndMap{}
			more = false
		case ']' == c:
			err = EndArray{}
			more = false
		case ',' == c:
			s.i--
			more = false
		case ' ' == c || '\t' == c:
			more = false
		case '.' == c:
			valf = float64(val)
			val = 0
			places = 1
		default:
			return fmt.Errorf("Bad num char: %v", string([]byte{c}))
		}
		if s.i >= len(s.data)-1 {
			more = false
		}
		s.i++
	}
	if places > 0 {
		s.v = valf + (float64(val)/float64(places))*float64(mult)
	} else {
		s.v = val * mult
	}
	return
}

type EndMap struct{}

func (e EndMap) Error() string {
	return "End of map structure encountered."
}

type EndArray struct{}

func (e EndArray) Error() string {
	return "End of array structure encountered."
}

func (s *State) readComma() (err error) {
	var more = true
	for more {
		switch {
		case s.data[s.i] == ',':
			more = false
		case s.data[s.i] == '}':
			s.i++
			return EndMap{}
		case s.data[s.i] == ']':
			s.i++
			return EndArray{}
		case s.i >= len(s.data)-1:
			return fmt.Errorf("No comma")
		}
		s.i++
	}
	return nil
}

func (s *State) readColon() (err error) {
	var more = true
	for more {
		switch {
		case s.data[s.i] == ':':
			more = false
		case s.i >= len(s.data)-1:
			return fmt.Errorf("No colon")
		}
		s.i++
	}
	return nil
}

func (s *State) readMap() (err error) {
	s.i++
	var (
		m   map[string]interface{}
		key string
	)
	m = make(map[string]interface{})
	for {
		if err = s.readString(); err != nil {
			return
		}
		key = s.v.(string)
		if err = s.readColon(); err != nil {
			return
		}
		if err = s.Read(); err != nil {
			if _, ok := err.(EndMap); !ok {
				return
			}
		}
		m[key] = s.v
		if _, ok := err.(EndMap); ok {
			break
		}
		if err = s.readComma(); err != nil {
			if _, ok := err.(EndMap); ok {
				break
			}
			return
		}
	}
	s.v = m
	return nil
}

func (s *State) readArray() (err error) {
	s.i++
	var (
		a []interface{}
	)
	a = make([]interface{}, 0, 10)
	for {
		if err = s.Read(); err != nil {
			if _, ok := err.(EndArray); !ok {
				return
			}
		}
		a = append(a, s.v)
		if _, ok := err.(EndArray); ok {
			break
		}
		if err = s.readComma(); err != nil {
			if _, ok := err.(EndArray); ok {
				break
			}
			return
		}
	}
	s.v = a
	return nil
}

func (s *State) readBool() (err error) {
	if strings.ToLower(string(s.data[s.i:s.i+4])) == "true" {
		s.i += 4
		s.v = true
	} else if strings.ToLower(string(s.data[s.i:s.i+5])) == "false" {
		s.i += 5
		s.v = false
	} else {
		err = fmt.Errorf("Could not parse boolean")
	}
	return
}

func (s *State) readNull() (err error) {
	if strings.ToLower(string(s.data[s.i:s.i+4])) == "null" {
		s.i += 4
		s.v = nil
	} else {
		err = fmt.Errorf("Could not parse null")
	}
	return
}

func Unmarshal(data []byte, v interface{}) error {
	state := &State{data, 0, v, make([]Event, 0, 10)}
	if err := state.Read(); err != nil {
		return err
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("Need a pointer, got %v", reflect.TypeOf(v))
	}
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	sv := reflect.ValueOf(state.v)
	for sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	var (
		rvt = rv.Type()
		svt = sv.Type()
	)
	if !svt.AssignableTo(rvt) {
		if rv.Kind() != reflect.Slice && sv.Kind() != reflect.Slice {
			return fmt.Errorf("Cannot assign %v to %v", svt, rvt)
		}
		var (
			mapi  map[string]interface{}
			mapt  = reflect.TypeOf(mapi)
			svte  = svt.Elem()
			rvte  = rvt.Elem()
			ismap bool
		)
		_, ismap = sv.Index(0).Interface().(map[string]interface{})
		if !(ismap && mapt.AssignableTo(rvte)) {
			return fmt.Errorf("Cannot assign %v to %v", svte, rvte)
		}
		var (
			ssv = reflect.MakeSlice(rvt, sv.Len(), sv.Cap())
		)
		for i := 0; i < sv.Len(); i++ {
			v := sv.Index(i).Interface().(map[string]interface{})
			ssv.Index(i).Set(reflect.ValueOf(v))
		}
		sv = ssv
	}
	rv.Set(sv)
	return nil
}
