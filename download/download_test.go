package download

import (
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestGetMongoDB(t *testing.T) {
	// Use a memory backed filesystem (no persistence)
	afs = afero.Afero{Fs: afero.NewMemMapFs()}

	tmpCache, err := afs.TempDir("", "")
	require.NoError(t, err)

	cfg := new(Config)
	cfg.mongoUrl = "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-amazon2-5.0.2.tgz"
	cfg.cachePath = path.Join(tmpCache, "mongod")

	// First call should download the file
	binPath, err := GetMongoDB(*cfg)
	require.NoError(t, err)

	require.Equal(t, cfg.cachePath, binPath)

	stat, err := afs.Stat(binPath)
	require.NoError(t, err)

	require.True(t, stat.Size() > 50000000)
	require.True(t, stat.Mode()&0100 != 0)

	// Second call should use the cached file
	binPath2, err := GetMongoDB(*cfg)
	require.NoError(t, err)

	require.Equal(t, binPath, binPath2)

	stat2, err := afs.Stat(binPath2)
	require.NoError(t, err)

	require.True(t, stat.Size() > 50000000)
	require.True(t, stat.Mode()&0100 != 0)
	require.Equal(t, stat.ModTime(), stat2.ModTime())
}
