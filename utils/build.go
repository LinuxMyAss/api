package utils

// Build represents this application's build data
type Build struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Branch  string `json:"branch"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

var build *Build

// GetBuild returns the stored build information
func GetBuild() *Build {
	return build
}

// SetBuild sets the stored build information
func SetBuild(b *Build) {
	build = b
}
