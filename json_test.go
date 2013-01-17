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
	"fmt"
	"testing"
)

func TestParseString(t *testing.T) {
	var (
		gold    = "Hello world"
		encoded = []byte(fmt.Sprintf("\"%v\"", gold))
		parsed  string
	)
	if err := Unmarshal(encoded, &parsed); err != nil {
		t.Fatalf("%v", err)
	}

	if gold != parsed {
		t.Fatalf("%v != %v", gold, parsed)
	}
}

func TestParseNumber(t *testing.T) {
	var (
		gold    int64 = 1234567
		encoded       = []byte(fmt.Sprintf("%v", gold))
		parsed  int64
	)
	if err := Unmarshal(encoded, &parsed); err != nil {
		t.Fatalf("%v", err)
	}
	if gold != parsed {
		t.Fatalf("%v != %v", gold, parsed)
	}
}

func TestParseNegativeNumber(t *testing.T) {
	var (
		gold    int64 = -1234567
		encoded       = []byte(fmt.Sprintf("%v", gold))
		parsed  int64
	)
	if err := Unmarshal(encoded, &parsed); err != nil {
		t.Fatalf("%v", err)
	}
	if gold != parsed {
		t.Fatalf("%v != %v", gold, parsed)
	}
}

func TestParseFloat(t *testing.T) {
	var (
		gold    float64 = 1234567.89
		encoded         = []byte("1234567.89")
		parsed  float64
	)
	if err := Unmarshal(encoded, &parsed); err != nil {
		t.Fatalf("%v", err)
	}
	if gold != parsed {
		t.Fatalf("%v != %v", gold, parsed)
	}
}

func TestParseMap(t *testing.T) {
	var (
		gold = map[string]interface{}{
			"foo": "Bar",
			"baz": 1234,
		}
		encoded = []byte("{\"foo\":\"Bar\",\"baz\":1234}")
		parsed  map[string]interface{}
	)
	if err := Unmarshal(encoded, &parsed); err != nil {
		t.Fatalf("%v", err)
	}
	if len(parsed) != len(gold) {
		t.Fatalf("Parsed len %v != gold len %v", len(parsed), len(gold))
	}
	for i, v := range parsed {
		if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", gold[i]) {
			t.Errorf("%v: %v != %v", i, v, gold[i])
		}
	}
}

func TestParseArray(t *testing.T) {
	var (
		gold    = []interface{}{1234, "Foo", 5678}
		encoded = []byte("[1234,\"Foo\",5678]")
		parsed  []interface{}
	)
	if err := Unmarshal(encoded, &parsed); err != nil {
		t.Fatalf("%v", err)
	}
	if len(parsed) != len(gold) {
		t.Fatalf("Parsed len %v != gold len %v", len(parsed), len(gold))
	}
	for i, v := range parsed {
		if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", gold[i]) {
			t.Errorf("%v: %v != %v", i, v, gold[i])
		}
	}
}

type Bucket map[string]interface{}
type BucketList []Bucket

func TestParseStruct(t *testing.T) {
	var (
		encoded = []byte("[{\"foo\":1},{\"foo\":2}]")
		gold    = BucketList{Bucket{"foo": 1}, Bucket{"foo": 2}}
		parsed  BucketList
	)
	if err := Unmarshal(encoded, &parsed); err != nil {
		t.Fatalf("%v", err)
	}
	if len(parsed) != len(gold) {
		t.Fatalf("Parsed len %v != gold len %v", len(parsed), len(gold))
	}
	for i, v := range parsed {
		if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", gold[i]) {
			t.Errorf("%v: %v != %v", i, v, gold[i])
		}
	}
}

const TWITTER_USER = `{"id":370773112,"id_str":"370773112","name":"fakekurrik","screen_name":"fakekurrik","location":"Trapped in factory","description":"I am just a testing account, following me probably won't gain you very much","url":"http:\/\/blog.roomanna.com","entities":{"url":{"urls":[{"url":"http:\/\/blog.roomanna.com","expanded_url":null,"indices":[0,24]}]},"description":{"urls":[]}},"protected":false,"followers_count":10,"friends_count":5,"listed_count":0,"created_at":"Fri Sep 09 16:13:20 +0000 2011","favourites_count":4,"utc_offset":-28800,"time_zone":"Pacific Time (US & Canada)","geo_enabled":true,"verified":false,"statuses_count":576,"lang":"en","status":{"created_at":"Thu Jan 17 19:00:46 +0000 2013","id":291983420479905792,"id_str":"291983420479905792","text":"http:\/\/t.co\/fHwmE7OI","source":"web","truncated":false,"in_reply_to_status_id":null,"in_reply_to_status_id_str":null,"in_reply_to_user_id":null,"in_reply_to_user_id_str":null,"in_reply_to_screen_name":null,"geo":null,"coordinates":null,"place":null,"contributors":null,"retweet_count":0,"entities":{"hashtags":[],"urls":[{"url":"http:\/\/t.co\/fHwmE7OI","expanded_url":"http:\/\/www.youtube.com\/watch?v=BS-FyAh9cv8","display_url":"youtube.com\/watch?v=BS-FyA\u2026","indices":[0,20]}],"user_mentions":[]},"favorited":false,"retweeted":false,"possibly_sensitive":false},"contributors_enabled":false,"is_translator":false,"profile_background_color":"C0DEED","profile_background_image_url":"http:\/\/a0.twimg.com\/profile_background_images\/616512781\/iarz5lvj7lg7zpg3zv8j.jpeg","profile_background_image_url_https":"https:\/\/si0.twimg.com\/profile_background_images\/616512781\/iarz5lvj7lg7zpg3zv8j.jpeg","profile_background_tile":true,"profile_image_url":"http:\/\/a0.twimg.com\/profile_images\/2440719659\/x47xdzkguqxr1w1gg5un_normal.png","profile_image_url_https":"https:\/\/si0.twimg.com\/profile_images\/2440719659\/x47xdzkguqxr1w1gg5un_normal.png","profile_banner_url":"https:\/\/si0.twimg.com\/profile_banners\/370773112\/1349887268","profile_link_color":"0084B4","profile_sidebar_border_color":"C0DEED","profile_sidebar_fill_color":"FFFFFF","profile_text_color":"333333","profile_use_background_image":true,"default_profile":false,"default_profile_image":false,"following":true,"follow_request_sent":false,"notifications":false}`

func TestParseTwitterUser(t *testing.T) {
	var (
		parsed map[string]interface{}
		status map[string]interface{}
	)
	if err := Unmarshal([]byte(TWITTER_USER), &parsed); err != nil {
		t.Fatalf("Could not parse Twitter user: %v", err)
	}
	if parsed["id"] != int64(370773112) {
		t.Fatalf("Could not parse 64-bit Twitter user ID.")
	}
	if parsed["name"] != "fakekurrik" {
		t.Fatalf("Could not parse Twitter user screen name.")
	}
	status = parsed["status"].(map[string]interface{})
	if status["id"] != int64(291983420479905792) {
		t.Fatalf("Could not parse nested 64-bit Tweet ID.")
	}
}
