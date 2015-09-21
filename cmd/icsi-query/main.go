package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ebfe/icsi"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <hash>\n", os.Args[0])
		os.Exit(1)
	}

	hash, err := hex.DecodeString(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: invalid hash: %s\n", os.Args[0], err)
		os.Exit(1)
	}

	r, err := icsi.Query(hash)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: query error: %s\n", os.Args[0], err)
		os.Exit(1)
	}

	if r != nil {
		fmt.Printf("version: %d\n", r.Version)
		fmt.Printf("first_seen: %s\n", r.FirstSeen.Format("2006-01-02"))
		fmt.Printf("last_seen: %s\n", r.LastSeen.Format("2006-01-02"))
		fmt.Printf("times_seen: %d\n", r.TimesSeen)
		fmt.Printf("validated: %t\n", r.Validated)
	}
}
