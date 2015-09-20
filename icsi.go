package icsi

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const notaryDomain = "notary.icsi.berkeley.edu"

type Status int

const (
	Unknown Status = iota
	Seen
	Validated
)

var (
	ipSeen      = net.IP{127, 0, 0, 1}
	ipValidated = net.IP{127, 0, 0, 2}

	errInvalidResponse = errors.New("icsi: invalid response")
	errUnknownVersion  = errors.New("icsi: unknown version")
	errMultipleRecords = errors.New("icsi: multiple records")
)

func dnsname(sha []byte) string {
	return fmt.Sprintf("%x.%s", sha, notaryDomain)
}

func isnxdomain(err error) bool {
	if err, ok := err.(*net.DNSError); ok {
		return err.Err == "no such host"
	}
	return false
}

func QueryStatus(sha []byte) (Status, error) {
	ips, err := net.LookupIP(dnsname(sha))
	if err != nil {
		if isnxdomain(err) {
			return Unknown, nil
		}
		return Unknown, err
	}

	if len(ips) != 1 {
		return Unknown, errMultipleRecords
	}
	if bytes.Equal(ips[0], ipSeen) {
		return Seen, nil
	}
	if bytes.Equal(ips[0], ipValidated) {
		return Validated, nil
	}

	return Unknown, nil
}

type Response struct {
	Version   int
	FirstSeen time.Time
	LastSeen  time.Time
	TimesSeen int
	Validated bool
}

func parseDate(s string) (time.Time, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(0, 0).AddDate(0, 0, int(i)), nil
}

func parseResponse(txt string) (*Response, error) {
	var r Response

	tok := strings.Split(txt, " ")
	for _, t := range tok {
		pair := strings.Split(t, "=")
		if len(pair) != 2 {
			return nil, errInvalidResponse
		}
		switch pair[0] {
		case "version":
			i, err := strconv.ParseInt(pair[1], 10, 32)
			if err != nil {
				return nil, errInvalidResponse
			}
			r.Version = int(i)
		case "first_seen":
			t, err := parseDate(pair[1])
			if err != nil {
				return nil, errInvalidResponse
			}
			r.FirstSeen = t
		case "last_seen":
			t, err := parseDate(pair[1])
			if err != nil {
				return nil, errInvalidResponse
			}
			r.LastSeen = t
		case "times_seen":
			i, err := strconv.ParseInt(pair[1], 10, 32)
			if err != nil {
				return nil, errInvalidResponse
			}
			r.TimesSeen = int(i)
		case "validated":
			i, err := strconv.ParseInt(pair[1], 10, 32)
			if err != nil {
				return nil, errInvalidResponse
			}
			r.Validated = i == 1
		}
	}

	if r.Version != 1 {
		return nil, errUnknownVersion
	}

	return &r, nil

}

func Query(sha []byte) (*Response, error) {
	txts, err := net.LookupTXT(dnsname(sha))
	if err != nil {
		if isnxdomain(err) {
			return nil, nil
		}
		return nil, err
	}

	if len(txts) != 1 {
		return nil, errMultipleRecords
	}

	return parseResponse(txts[0])

}
