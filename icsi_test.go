package icsi

import (
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
