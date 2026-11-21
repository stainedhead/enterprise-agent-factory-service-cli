package output

// ColorEnabled decides whether to color: "on" always, "off" never, "auto" only on a terminal
// when NO_COLOR is unset (NFR accessibility).
func ColorEnabled(mode string, tty bool, getenv func(string) string) bool {
	switch mode {
	case "on":
		return true
	case "off":
		return false
	}
	return tty && getenv("NO_COLOR") == "" && getenv("TERM") != "dumb"
}

func bold(s string, on bool) string {
	if !on {
		return s
	}
	return "\x1b[1m" + s + "\x1b[0m"
}
