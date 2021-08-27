package download

import (
	"path"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/spf13/afero"
)

func TestGetMongoDB(t *testing.T) {
	// Use a memory backed filesystem (no persistence)
	afs = afero.Afero{Fs: afero.NewMemMapFs()}

	tmpCache, _ := afs.TempDir("", "")

	var fileModificationTime time.Time

	Convey("Given a config object", t, func() {
		cfg := new(Config)
		cfg.mongoUrl = "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-amazon2-5.0.2.tgz"
		cfg.cachePath = path.Join(tmpCache, "mongod")

		Convey("When we call GetMongoDB", func() {
			binPath, err := GetMongoDB(*cfg)
			Convey("Then it downloads the file when there is not one in cache", func() {
				So(err, ShouldBeNil)
				So(binPath, ShouldEqual, cfg.cachePath)

				stat, err := afs.Stat(binPath)
				So(err, ShouldBeNil)
				So(stat.Size(), ShouldBeGreaterThan, 50000000)
				So(stat.Mode()&0100, ShouldNotBeZeroValue)
				fileModificationTime = stat.ModTime()
			})
			Convey("And it finds the file in cache if present", func() {
				So(err, ShouldBeNil)
				So(binPath, ShouldEqual, cfg.cachePath)

				stat, err := afs.Stat(binPath)
				So(err, ShouldBeNil)
				So(stat.Size(), ShouldBeGreaterThan, 50000000)
				So(stat.Mode()&0100, ShouldNotBeZeroValue)
				// It should not have changed since when it was first downloaded
				So(stat.ModTime(), ShouldEqual, fileModificationTime)
			})
		})
	})
}
