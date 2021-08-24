package download

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDownloadURL(t *testing.T) {

	tests := map[string]struct {
		spec          *DownloadSpec
		expectedURL   string
		expectedError string
	}{
		"mac": {
			spec: &DownloadSpec{
				Version:  "5.0.2",
				Platform: "osx",
				Arch:     "x86_64",
				OSName:   "irrelevant",
			},
			expectedURL: "https://fastdl.mongodb.org/osx/mongodb-macos-x86_64-5.0.2.tgz",
		},
		"linux-os": {
			spec: &DownloadSpec{
				Version:  "5.0.2",
				Platform: "linux",
				Arch:     "x86_64",
				OSName:   "ubuntu2004",
			},
			expectedURL: "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-ubuntu2004-5.0.2.tgz",
		},
		"linux-no-os": {
			spec: &DownloadSpec{
				Version:  "4.0.26",
				Platform: "linux",
				Arch:     "x86_64",
				OSName:   "",
			},
			expectedURL: "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-4.0.26.tgz",
		},
		"win32": {
			spec: &DownloadSpec{
				Version:  "5.0.2",
				Platform: "win32",
				Arch:     "x86_64",
				OSName:   "",
			},
			expectedError: "invalid spec: unsupported platform win32",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			spec := &DownloadSpec{
				Version:  test.spec.Version,
				Platform: test.spec.Platform,
				Arch:     test.spec.Arch,
				OSName:   test.spec.OSName,
			}

			actualURL, err := spec.GetDownloadURL()

			if test.expectedError != "" {
				require.Equal(t, test.expectedError, err.Error())
			}

			if test.expectedURL != "" {
				require.Equal(t, test.expectedURL, actualURL)
			}
		})
	}
}

func TestMakeDownloadSpec(t *testing.T) {

	tests := map[string]struct {
		version      string
		os           string
		arch         string
		linuxId      string
		linuxVersion string
		expectedSpec *DownloadSpec
		expectedErr  string
	}{
		"Mac": {
			version: "5.0.4",
			os:      "darwin",
			arch:    "amd64",
			expectedSpec: &DownloadSpec{
				Version:  "5.0.4",
				Arch:     "x86_64",
				Platform: "osx",
				OSName:   "",
			},
		},
		"Ubuntu 20.04": {
			version:      "5.0.4",
			os:           "linux",
			arch:         "amd64",
			linuxId:      "ubuntu",
			linuxVersion: "20.04",
			expectedSpec: &DownloadSpec{
				Version:  "5.0.4",
				Arch:     "x86_64",
				Platform: "linux",
				OSName:   "ubuntu2004",
			},
		},
		"Ubuntu 20.10": {
			version:      "5.0.4",
			os:           "linux",
			arch:         "amd64",
			linuxId:      "ubuntu",
			linuxVersion: "20.10",
			expectedSpec: &DownloadSpec{
				Version:  "5.0.4",
				Arch:     "x86_64",
				Platform: "linux",
				OSName:   "ubuntu2004",
			},
		},
		"Ubuntu 18.04": {
			version:      "5.0.4",
			os:           "linux",
			arch:         "amd64",
			linuxId:      "ubuntu",
			linuxVersion: "18.04",
			expectedSpec: &DownloadSpec{
				Version:  "5.0.4",
				Arch:     "x86_64",
				Platform: "linux",
				OSName:   "ubuntu1804",
			},
		},
		"Ubuntu 16.04": {
			version:      "5.0.4",
			os:           "linux",
			arch:         "amd64",
			linuxId:      "ubuntu",
			linuxVersion: "16.04",
			expectedSpec: &DownloadSpec{
				Version:  "5.0.4",
				Arch:     "x86_64",
				Platform: "linux",
				OSName:   "ubuntu1604",
			},
		},
		"Old Ubuntu": {
			version:      "5.0.4",
			os:           "linux",
			arch:         "amd64",
			linuxId:      "ubuntu",
			linuxVersion: "14.04",
			expectedErr:  "unsupported system: invalid ubuntu version 14 (min 16)",
		},
		"Debian 10": {
			version:      "5.0.4",
			os:           "linux",
			arch:         "amd64",
			linuxId:      "debian",
			linuxVersion: "10",
			expectedSpec: &DownloadSpec{
				Version:  "5.0.4",
				Arch:     "x86_64",
				Platform: "linux",
				OSName:   "debian10",
			},
		},
		"Debian 9.2": {
			version:      "5.0.4",
			os:           "linux",
			arch:         "amd64",
			linuxId:      "debian",
			linuxVersion: "9.2",
			expectedSpec: &DownloadSpec{
				Version:  "5.0.4",
				Arch:     "x86_64",
				Platform: "linux",
				OSName:   "debian92",
			},
		},
		"Old Debian": {
			version:      "5.0.4",
			os:           "linux",
			arch:         "amd64",
			linuxId:      "debian",
			linuxVersion: "8.1",
			expectedErr:  "unsupported system: invalid debian version 8 (min 9)",
		},
		"Other Linux": {
			version:      "5.0.4",
			os:           "linux",
			arch:         "amd64",
			linuxId:      "fedora",
			linuxVersion: "17",
			expectedErr:  "unsupported system: invalid linux version 'fedora'",
		},
		"Old Mongo version": {
			version:      "4.2.15",
			os:           "linux",
			arch:         "amd64",
			linuxId:      "ubuntu",
			linuxVersion: "20.04",
			expectedErr:  "unsupported MongoDB version \"4.2.15\": only version 4.4 and above are supported",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(tt *testing.T) {

			goOS = tc.os
			goArch = tc.arch

			getOsReleaseContent = func() (map[string]string, error) {
				return map[string]string{
					"ID":         tc.linuxId,
					"VERSION_ID": tc.linuxVersion,
				}, nil
			}

			spec, err := MakeDownloadSpec(tc.version)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedSpec, spec)
			}

			if tc.expectedSpec != nil {
				require.Equal(t, tc.expectedSpec, spec)
			} else {
				require.Nil(t, spec)
			}
		})
	}
}
