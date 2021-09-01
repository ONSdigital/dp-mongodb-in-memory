package download

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/spf13/afero"
)

func TestGetMongoDB(t *testing.T) {
	const (
		validMongodTarball   = "/mongodb-test.tgz"
		invalidMongodTarball = "/random.tgz"
		notTarball           = "/test"
	)

	// Use a memory backed filesystem (no persistence)
	afs = afero.Afero{Fs: afero.NewMemMapFs()}

	tmpCache, _ := afs.TempDir("", "")

	Convey("Having set up a mocked server", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case invalidMongodTarball:
				fallthrough
			case validMongodTarball:
				f, err := os.Open("testdata" + r.URL.Path)
				if err != nil {
					http.NotFound(w, r)
					return
				}
				defer f.Close()
				io.Copy(w, f)
			case notTarball:
				w.Write([]byte("Test data"))
			default:
				http.NotFound(w, r)
			}
		}))

		cfg := new(Config)
		cfg.cachePath = path.Join(tmpCache, "mongod")

		Convey("When the mongod exec file is not in cache", func() {
			afs.Remove(cfg.cachePath)
			Convey("And the requested url exists", func() {
				cfg.mongoUrl = ts.URL + validMongodTarball
				Convey("Then it downloads the tarball and stores the exec file in cache", func() {
					startTime := time.Now()

					binPath, err := GetMongoDB(*cfg)
					So(err, ShouldBeNil)
					So(binPath, ShouldEqual, cfg.cachePath)

					stat, err := afs.Stat(binPath)
					So(err, ShouldBeNil)
					So(stat.Size(), ShouldBeGreaterThan, 0)
					So(stat.Mode()&0100, ShouldNotBeZeroValue)
					So(stat.ModTime(), ShouldHappenBetween, startTime, time.Now())
				})
			})
			Convey("And the requested url can not be found", func() {
				cfg.mongoUrl = ts.URL + "/invalid"
				Convey("Then an error is returned", func() {

					binPath, err := GetMongoDB(*cfg)
					So(err, ShouldBeError)
					So(binPath, ShouldBeBlank)
				})
			})
			Convey("And the requested url is not a tarball", func() {
				cfg.mongoUrl = ts.URL + notTarball
				Convey("Then an error is returned", func() {

					binPath, err := GetMongoDB(*cfg)
					So(err, ShouldBeError)
					So(binPath, ShouldBeBlank)
				})
			})
			Convey("And the requested url is a tarball not containing a mongod file", func() {
				cfg.mongoUrl = ts.URL + invalidMongodTarball
				Convey("Then an error is returned", func() {

					binPath, err := GetMongoDB(*cfg)
					So(err, ShouldBeError)
					So(binPath, ShouldBeBlank)
				})
			})
		})

		Convey("When the mongod exec file is found in cache", func() {
			afs.Create(cfg.cachePath)

			Convey("Then it uses the file in cache and it does not download it again", func() {
				cfg.mongoUrl = ts.URL + "/should-not-be-called"

				binPath, err := GetMongoDB(*cfg)
				So(err, ShouldBeNil)
				So(binPath, ShouldEqual, cfg.cachePath)
			})
		})

		Reset(func() {
			ts.Close()
			afs.Remove(cfg.cachePath)
		})

	})
}
