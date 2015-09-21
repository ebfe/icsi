package icsi

import (
	"encoding/hex"
	"reflect"
	"testing"
	"time"
)

func TestParseResponse(t *testing.T) {
	txt := "version=1 first_seen=15387 last_seen=15646 times_seen=260 validated=1"
	out := &Response{
		Version:   1,
		FirstSeen: time.Date(2012, time.February, 17, 0, 0, 0, 0, time.UTC),
		LastSeen:  time.Date(2012, time.November, 2, 0, 0, 0, 0, time.UTC),
		TimesSeen: 260,
		Validated: true,
	}

	r, err := parseResponse(txt)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(out, r) {
		t.Fatalf("expected: %+v\ngot: %+v\n", out, r)
	}
}

func checkStatus(t *testing.T, hash []byte, expected Status) {
	status, err := QueryStatus(hash)
	if err != nil {
		t.Fatal(err)
	}
	if status != expected {
		t.Fatalf("unexpected status %v %x", status, hash)
	}
}

func TestQueryStatus(t *testing.T) {
	shash, _ := hex.DecodeString("C1956DC8A7DFB2A5A56934DA09778E3A11023358")
	checkStatus(t, shash, Seen)
	checkStatus(t, make([]byte, 16), Unknown)
}
