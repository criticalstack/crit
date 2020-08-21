// Package buildinfo contains any build time information that needs to be
// available at run time.
package buildinfo

import "runtime"

var (
	Date string

	GitSHA string

	GitTreeState = "clean"

	GoVersion = runtime.Version()

	Version string
)
