package download

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/acobaugh/osrelease"
)

// We define these as package vars so we can override it in tests
var goOS = runtime.GOOS
var goArch = runtime.GOARCH
var getOsReleaseContent = func() (map[string]string, error) {
	return osrelease.ReadFile("/etc/os-release")
}

// DownloadSpec specifies what copy of MongoDB to download
type DownloadSpec struct {
	// Version is what version of MongoDB to download
	Version string

	// Platform is "osx" or "linux"
	Platform string

	// Arch
	Arch string

	// OSName is one of:
	// - ubuntu2004
	// - ubuntu1804
	// - ubuntu1604
	// - debian10
	// - debian92
	// - "" for MacOS
	OSName string
}

// MakeDownloadSpec returns a DownloadSpec for the current operating system
func MakeDownloadSpec(version string) (*DownloadSpec, error) {
	parsedVersion, versionErr := parseVersion(version)
	if versionErr != nil {
		return nil, versionErr
	}

	platform, platformErr := detectPlatform()
	if platformErr != nil {
		return nil, platformErr
	}

	arch, archErr := detectArch()
	if archErr != nil {
		return nil, archErr
	}

	osName, osErr := detectLinuxId(parsedVersion)
	if osErr != nil {
		return nil, osErr
	}

	return &DownloadSpec{
		Version:  version,
		Arch:     arch,
		Platform: platform,
		OSName:   osName,
	}, nil
}

// GetDownloadURL returns the download URL to download the binary
// from the MongoDB website
func (spec *DownloadSpec) GetDownloadURL() (string, error) {
	archiveName := "mongodb-"

	switch spec.Platform {
	case "linux":
		archiveName += "linux-" + spec.Arch

		if spec.OSName != "" {
			archiveName += "-" + spec.OSName
		}

	case "osx":
		archiveName += "macos-" + spec.Arch
	default:
		return "", fmt.Errorf("invalid spec: unsupported platform " + spec.Platform)
	}

	return fmt.Sprintf(
		"https://fastdl.mongodb.org/%s/%s-%s.tgz",
		spec.Platform,
		archiveName,
		spec.Version,
	), nil
}

// Parse a version string into an array [major, minor, patch]
func parseVersion(version string) ([]int, error) {
	versionParts := strings.Split(version, ".")
	if len(versionParts) < 3 {
		return nil, &UnsupportedMongoVersionError{
			version: version,
			msg:     "MongoDB version number must be in the form x.y.z",
		}
	}

	majorVersion, majErr := strconv.Atoi(versionParts[0])
	if majErr != nil {
		return nil, &UnsupportedMongoVersionError{
			version: version,
			msg:     "could not parse major version",
		}
	}

	minorVersion, minErr := strconv.Atoi(versionParts[1])
	if minErr != nil {
		return nil, &UnsupportedMongoVersionError{
			version: version,
			msg:     "could not parse minor version",
		}
	}

	patchVersion, patchErr := strconv.Atoi(versionParts[2])
	if patchErr != nil {
		return nil, &UnsupportedMongoVersionError{
			version: version,
			msg:     "could not parse patch version",
		}
	}

	if (majorVersion < 4) || (majorVersion == 4 && minorVersion < 4) {
		return nil, &UnsupportedMongoVersionError{
			version: version,
			msg:     "only version 4.4 and above are supported",
		}
	}

	return []int{majorVersion, minorVersion, patchVersion}, nil
}

func versionGTE(a []int, b []int) bool {
	if a[0] > b[0] {
		return true
	}

	if a[0] < b[0] {
		return false
	}

	if a[1] > b[1] {
		return true
	}

	if a[1] < b[1] {
		return false
	}

	return a[2] >= b[2]
}

func detectPlatform() (string, error) {
	switch goOS {
	case "darwin":
		return "osx", nil
	case "linux":
		return "linux", nil
	default:
		return "", &UnsupportedSystemError{msg: "OS " + goOS + " not supported"}
	}
}

func detectArch() (string, error) {
	switch goArch {
	case "amd64":
		return "x86_64", nil
	default:
		return "", &UnsupportedSystemError{msg: "architecture " + goArch + " not supported"}
	}
}

func detectLinuxId(mongoVersion []int) (string, error) {
	if goOS != "linux" {
		// Not on Linux
		return "", nil
	}

	osRelease, osReleaseErr := getOsReleaseContent()
	if osReleaseErr != nil {
		return "", osReleaseErr
	}

	id := osRelease["ID"]
	versionString := strings.Split(osRelease["VERSION_ID"], ".")[0]
	version, versionErr := strconv.Atoi(versionString)
	if versionErr != nil {
		return "", &UnsupportedSystemError{msg: "invalid version number " + versionString}
	}
	switch id {
	case "ubuntu":
		if version >= 20 {
			return "ubuntu2004", nil
		}
		if version >= 18 {
			return "ubuntu1804", nil
		}
		if version >= 16 {
			return "ubuntu1604", nil
		}
		return "", &UnsupportedSystemError{msg: "invalid ubuntu version " + versionString + " (min 16)"}
	case "debian":
		if version >= 10 {
			return "debian10", nil
		}
		if version >= 9 {
			return "debian92", nil
		}
		return "", &UnsupportedSystemError{msg: "invalid debian version " + versionString + " (min 9)"}
	default:
		return "", &UnsupportedSystemError{msg: "invalid linux version '" + id + "'"}
	}
}
