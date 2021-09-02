package download

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewConfig(t *testing.T) {
	var originalGetDownloadUrl = getDownloadUrl
	var originalGetEnv = getEnv
	var originalGoOs = goOS
	Convey("Given a download url can be found", t, func() {

		Convey("And the url is valid", func() {
			filename := "mongodb-linux-x86_64-ubuntu2004-5.0.2.tgz"
			mongoUrl := "https://fastdl.mongodb.org/linux/" + filename

			getDownloadUrl = func(version string) (string, error) {
				return mongoUrl, nil
			}
			Convey("When XDG_CACHE_HOME env var is set", func() {
				getEnv = func(key string) string {
					if key == "XDG_CACHE_HOME" {
						return "/cache/home"
					}
					return ""
				}
				Convey("Then NewConfig uses the XDG_CACHE_HOME env var to determine the cache path", func() {
					cfg, err := NewConfig("version")
					So(err, ShouldBeNil)
					So(cfg.mongoUrl, ShouldEqual, mongoUrl)
					So(cfg.cachePath, ShouldEqual, "/cache/home/dp-mongodb-in-memory/"+filename+"/mongod")
				})
			})
			Convey("When XDG_CACHE_HOME env var is not set", func() {
				userHome := "/usr/home"
				getEnv = func(key string) string {
					if key == "HOME" {
						return userHome
					}
					return ""
				}
				Convey("And running on OSX", func() {
					goOS = "darwin"
					Convey("Then NewConfig determines the right home cache path", func() {
						cfg, err := NewConfig("version")
						So(err, ShouldBeNil)
						So(cfg.mongoUrl, ShouldEqual, mongoUrl)
						So(cfg.cachePath, ShouldEqual, userHome+"/Library/Caches/dp-mongodb-in-memory/"+filename+"/mongod")
					})
				})
				Convey("And running on Linux", func() {
					goOS = "linux"
					Convey("Then NewConfig determines the right home cache path", func() {
						cfg, err := NewConfig("version")
						So(err, ShouldBeNil)
						So(cfg.mongoUrl, ShouldEqual, mongoUrl)
						So(cfg.cachePath, ShouldEqual, userHome+"/.cache/dp-mongodb-in-memory/"+filename+"/mongod")
					})
				})
				Convey("And running on Windows", func() {
					goOS = "win32"
					Convey("Then NewConfig errors", func() {
						cfg, err := NewConfig("version")
						So(cfg, ShouldBeNil)
						So(err, ShouldBeError)
						expectedError := &UnsupportedSystemError{msg: "OS 'win32'"}
						So(err, ShouldResemble, expectedError)
						So(err, ShouldHaveSameTypeAs, expectedError)
					})
				})
				Reset(func() {
					goOS = originalGoOs
				})
			})
		})

		Convey("And the url is invalid", func() {
			getDownloadUrl = func(version string) (string, error) {
				return ":invalid", nil
			}
			Convey("Then NewConfig errors", func() {
				cfg, err := NewConfig("version")
				So(err, ShouldBeError)
				So(cfg, ShouldBeNil)
			})
		})

		Reset(func() {
			getDownloadUrl = originalGetDownloadUrl
			getEnv = originalGetEnv
		})

	})

	Convey("Given an error occurs when determining the download url", t, func() {
		expectedError := errors.New("unsupported system")
		getDownloadUrl = func(version string) (string, error) {
			return "", expectedError
		}

		Convey("Then NewConfig errors", func() {
			cfg, err := NewConfig("version")
			So(cfg, ShouldBeNil)
			So(err, ShouldEqual, expectedError)
		})

		Reset(func() {
			getDownloadUrl = originalGetDownloadUrl
		})
	})
}
