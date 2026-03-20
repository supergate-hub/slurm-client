package slurmclient

// Version represents a Slurm REST API version.
type Version string

const (
	V0039 Version = "v0.0.39"
	V0040 Version = "v0.0.40"
	V0041 Version = "v0.0.41"
	V0042 Version = "v0.0.42"
	V0043 Version = "v0.0.43"
	V0044 Version = "v0.0.44"
)

// String returns the version string.
func (v Version) String() string {
	return string(v)
}

// IsValid reports whether the version is a known Slurm API version.
func (v Version) IsValid() bool {
	switch v {
	case V0039, V0040, V0041, V0042, V0043, V0044:
		return true
	}
	return false
}

// supportedVersions lists versions in descending order (newest first).
var supportedVersions = []Version{V0044, V0043, V0042, V0041, V0040, V0039}
