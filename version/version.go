package version
import _ "embed"

var (
	// BuildTime is a time label of the moment when the binary was built
	BuildTime = "unset"
	// Commit is a last commit hash at the moment when the binary was built
	Commit = "unset"
	// Release is a semantic Version of current build
	//go:embed version
	Release string
)

func Show() string {
	return "hello couchbase_exporter\n" +
		"Release: " + Release + "\n" +
		"Commit: " + Commit + "\n" +
		"Build Time: " + BuildTime
}
