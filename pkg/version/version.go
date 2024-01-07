package version

import (
	"fmt"
	"io"
	"runtime"
)

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	GitSource   string
	GitTag      string
	GitBranch   string
	GitHash     string
	GoBuildTime string
)

func Print(w io.Writer) {
	if GitSource != "" {
		fmt.Fprintf(w, "  Url: https://%v\n", GitSource)
	}
	if GitTag != "" || GitBranch != "" {
		fmt.Fprintf(w, "  Version: %v (branch: %q hash:%q)\n", GitTag, GitBranch, GitHash)
	}
	if GoBuildTime != "" {
		fmt.Fprintf(w, "  Built: %v\n", GoBuildTime)
	}
	fmt.Fprintf(w, "  Go: %v (%v/%v)\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
