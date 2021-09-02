package download

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDownloadURL(t *testing.T) {
	Convey("Given a DownloadSpec object", t, func() {
		spec := &DownloadSpec{
			Version: "5.0.2",
			Arch:    "x86_64",
		}
		Convey("When platform is Mac", func() {
			spec.Platform = "osx"
			Convey("Then GetDownloadURL builds the right url", func() {
				url, err := spec.GetDownloadURL()

				So(err, ShouldBeNil)
				So(url, ShouldEqual, "https://fastdl.mongodb.org/osx/mongodb-macos-x86_64-5.0.2.tgz")
			})
		})

		Convey("When platform is Linux", func() {
			spec.Platform = "linux"
			Convey("And a linux id is provided", func() {
				spec.OSName = "ubuntu2004"
				Convey("Then GetDownloadURL builds the right url", func() {
					url, err := spec.GetDownloadURL()

					So(err, ShouldBeNil)
					So(url, ShouldEqual, "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-ubuntu2004-5.0.2.tgz")
				})
			})
			Convey("And no linux id is provided", func() {
				spec.OSName = ""
				Convey("Then an error is thrown", func() {
					url, err := spec.GetDownloadURL()

					So(url, ShouldBeBlank)
					So(err, ShouldBeError)
					So(err, ShouldResemble, errors.New("invalid spec: OS name not provided"))
				})
			})
		})

		Convey("When platform is not supported", func() {
			spec.Platform = "win32"

			Convey("Then an error is thrown", func() {
				url, err := spec.GetDownloadURL()

				So(url, ShouldBeBlank)
				So(err, ShouldBeError)
				So(err, ShouldResemble, errors.New("invalid spec: unsupported platform win32"))
			})
		})
	})
}

func TestMakeDownloadSpec(t *testing.T) {
	var originalGoOs = goOS
	var originalGoArch = goArch

	Convey("Given a valid MongoDB version", t, func() {
		version := "5.0.4"

		Convey("When running on a x64 architecture", func() {
			goArch = "amd64"
			Convey("And on Mac", func() {
				goOS = "darwin"
				Convey("Then the returned spec is correct", func() {
					spec, err := MakeDownloadSpec(version)

					So(err, ShouldBeNil)
					So(spec, ShouldResemble, &DownloadSpec{
						Version:  "5.0.4",
						Arch:     "x86_64",
						Platform: "osx",
					})
				})
			})

			Convey("And on Linux", func() {
				goOS = "linux"

				tests := map[string]struct {
					linuxId      string
					linuxVersion string
					expectedSpec *DownloadSpec
					expectedErr  error
				}{
					"Ubuntu 20.04": {
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
						linuxId:      "ubuntu",
						linuxVersion: "14.04",
						expectedErr:  &UnsupportedSystemError{msg: "invalid ubuntu version 14 (min 16)"},
					},
					"Debian 10": {
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
						linuxId:      "debian",
						linuxVersion: "8.1",
						expectedErr:  &UnsupportedSystemError{msg: "invalid debian version 8 (min 9)"},
					},
					"Other Linux": {
						linuxId:      "fedora",
						linuxVersion: "17",
						expectedErr:  &UnsupportedSystemError{msg: "invalid linux version 'fedora'"},
					},
					"Invalid linux version": {
						linuxId:      "fedora",
						linuxVersion: "vvv111",
						expectedErr:  &UnsupportedSystemError{msg: "invalid version number vvv111"},
					},
				}
				for name, tc := range tests {
					Convey(name, func() {
						getOsReleaseContent = func() (map[string]string, error) {
							return map[string]string{
								"ID":         tc.linuxId,
								"VERSION_ID": tc.linuxVersion,
							}, nil
						}
						Convey("Then the returned spec is correct", func() {
							spec, err := MakeDownloadSpec(version)
							if tc.expectedErr != nil {
								So(err, ShouldBeError)
								So(err, ShouldResemble, tc.expectedErr)
								So(err, ShouldHaveSameTypeAs, tc.expectedErr)
							} else {
								So(err, ShouldBeNil)
							}

							if tc.expectedSpec != nil {
								So(spec, ShouldResemble, tc.expectedSpec)
							} else {
								So(spec, ShouldBeNil)
							}
						})
					})
				}

				Convey("When there is an error reading the os-release file", func() {
					expectedErr := errors.New("could not read file")
					getOsReleaseContent = func() (map[string]string, error) {
						return nil, expectedErr
					}
					Convey("Then an error is returned", func() {
						spec, err := MakeDownloadSpec(version)
						So(err, ShouldEqual, expectedErr)
						So(spec, ShouldBeNil)
					})
				})
			})

			Convey("And on a non supported platform", func() {
				goOS = "win32"
				Convey("Then an error is returned", func() {
					spec, err := MakeDownloadSpec(version)
					So(spec, ShouldBeNil)
					So(err, ShouldBeError)
					expectedError := &UnsupportedSystemError{msg: "OS " + goOS + " not supported"}
					So(err, ShouldResemble, expectedError)
					So(err, ShouldHaveSameTypeAs, expectedError)
				})
			})

			Reset(func() {
				goOS = originalGoOs
			})
		})

		Convey("When running on a non supported architecture", func() {
			goArch = "386"
			Convey("Then an error is returned", func() {
				spec, err := MakeDownloadSpec(version)
				So(spec, ShouldBeNil)
				So(err, ShouldBeError)
				expectedError := &UnsupportedSystemError{msg: "architecture " + goArch + " not supported"}
				So(err, ShouldResemble, expectedError)
				So(err, ShouldHaveSameTypeAs, expectedError)
			})
		})

		Reset(func() {
			goArch = originalGoArch
		})
	})

	Convey("Given an invalid MongoDB version", t, func() {
		Convey("Without periods", func() {
			version := "version"
			Convey("Then an error is returned", func() {
				spec, err := MakeDownloadSpec(version)
				So(spec, ShouldBeNil)
				So(err, ShouldBeError)
				expectedError := &UnsupportedMongoVersionError{
					version: version,
					msg:     "MongoDB version number must be in the form x.y.z",
				}
				So(err, ShouldResemble, expectedError)
				So(err, ShouldHaveSameTypeAs, expectedError)
			})
		})
		Convey("With less than 2 periods", func() {
			version := "1.2"
			Convey("Then an error is returned", func() {
				spec, err := MakeDownloadSpec(version)
				So(spec, ShouldBeNil)
				So(err, ShouldBeError)
				expectedError := &UnsupportedMongoVersionError{
					version: version,
					msg:     "MongoDB version number must be in the form x.y.z",
				}
				So(err, ShouldResemble, expectedError)
				So(err, ShouldHaveSameTypeAs, expectedError)
			})
		})
		Convey("With more than 2 periods", func() {
			version := "2.1.0.a"
			Convey("Then an error is returned", func() {
				spec, err := MakeDownloadSpec(version)
				So(spec, ShouldBeNil)
				So(err, ShouldBeError)
				expectedError := &UnsupportedMongoVersionError{
					version: version,
					msg:     "MongoDB version number must be in the form x.y.z",
				}
				So(err, ShouldResemble, expectedError)
				So(err, ShouldHaveSameTypeAs, expectedError)
			})
		})
		Convey("With an invalid major version", func() {
			version := "a.1.0"
			Convey("Then an error is returned", func() {
				spec, err := MakeDownloadSpec(version)
				So(spec, ShouldBeNil)
				So(err, ShouldBeError)
				expectedError := &UnsupportedMongoVersionError{
					version: version,
					msg:     "could not parse major version",
				}
				So(err, ShouldResemble, expectedError)
				So(err, ShouldHaveSameTypeAs, expectedError)
			})
		})
		Convey("With an invalid minor version", func() {
			version := "4.minor.0"
			Convey("Then an error is returned", func() {
				spec, err := MakeDownloadSpec(version)
				So(spec, ShouldBeNil)
				So(err, ShouldBeError)
				expectedError := &UnsupportedMongoVersionError{
					version: version,
					msg:     "could not parse minor version",
				}
				So(err, ShouldResemble, expectedError)
				So(err, ShouldHaveSameTypeAs, expectedError)
			})
		})
		Convey("With an invalid patch version", func() {
			version := "4.7.pp"
			Convey("Then an error is returned", func() {
				spec, err := MakeDownloadSpec(version)
				So(spec, ShouldBeNil)
				So(err, ShouldBeError)
				expectedError := &UnsupportedMongoVersionError{
					version: version,
					msg:     "could not parse patch version",
				}
				So(err, ShouldResemble, expectedError)
				So(err, ShouldHaveSameTypeAs, expectedError)
			})
		})
		Convey("With a non-supported old version", func() {
			version := "4.2.15"
			Convey("Then an error is returned", func() {
				spec, err := MakeDownloadSpec(version)
				So(spec, ShouldBeNil)
				So(err, ShouldBeError)
				expectedError := &UnsupportedMongoVersionError{
					version: version,
					msg:     "only version 4.4 and above are supported",
				}
				So(err, ShouldResemble, expectedError)
				So(err, ShouldHaveSameTypeAs, expectedError)
			})
		})
	})
}
