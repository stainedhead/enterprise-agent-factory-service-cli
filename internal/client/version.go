package client

import (
	"fmt"
	"strconv"
	"strings"
)

func majorMinor(v string) (int, int, bool) {
	parts := strings.Split(strings.TrimSpace(v), ".")
	if len(parts) < 2 {
		return 0, 0, false
	}
	a, e1 := strconv.Atoi(parts[0])
	b, e2 := strconv.Atoi(parts[1])
	return a, b, e1 == nil && e2 == nil
}

// contractNotice returns the one-per-session notice for a server contract version, or "".
func contractNotice(ours, server string) string {
	om, on, ok1 := majorMinor(ours)
	sm, sn, ok2 := majorMinor(server)
	if !ok1 || !ok2 {
		return ""
	}
	switch {
	case sm != om:
		return fmt.Sprintf("Notice: the server speaks API contract %s but this afs was built for contract/%d.%d; commands may not work. Install a matching afs.", server, om, on)
	case sn > on:
		return fmt.Sprintf("Notice: the server speaks API contract %s, newer than contract/%d.%d of this afs; new fields and operations are unavailable until you upgrade afs.", server, om, on)
	}
	return ""
}
