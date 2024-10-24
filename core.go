package core

import (
	"runtime"

	"github.com/gzjjjfree/gzv2ray/common/serial"
)

var (
	version  = "4.37.3"
	build    = "Custom"
	codename = "V2Fly, a community-driven edition of V2Ray."
	intro    = "To learn."
)

// Version returns V2Ray's version as a string, in the form of "x.y.z" where x, y and z are numbers.
// ".z" part may be omitted in regular releases.
func Version() string {
	return version
}

// VersionStatement returns a list of strings representing the full version info.
func VersionStatement() []string {
	return []string{
		serial.Concat("V2Ray ", Version(), " (", codename, ") ", build, " (", runtime.Version(), " ", runtime.GOOS, "/", runtime.GOARCH, ")"),
		intro,
	}
}
