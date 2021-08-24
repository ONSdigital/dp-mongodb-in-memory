package download

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	filename := "mongodb-linux-x86_64-ubuntu2004-5.0.2.tgz"
	mongoUrl := "https://fastdl.mongodb.org/linux/" + filename

	getDownloadUrl = func(version string) (string, error) {
		return mongoUrl, nil
	}

	tests := map[string]struct {
		os            string
		cacheHome     string
		userHome      string
		expectedCache string
		expectedErr   string
	}{
		"XDG_CACHE_HOME set": {
			cacheHome:     "/cache/home",
			expectedCache: "/cache/home/dp-mongodb-in-memory/" + filename + "/mongod",
		},
		"OSX default cache": {
			os:            "darwin",
			cacheHome:     "",
			userHome:      "/usr/test",
			expectedCache: "/usr/test/Library/Caches/dp-mongodb-in-memory/" + filename + "/mongod",
		},
		"Linux default cache": {
			os:            "linux",
			cacheHome:     "",
			userHome:      "/home/test",
			expectedCache: "/home/test/.cache/dp-mongodb-in-memory/" + filename + "/mongod",
		},
		"Unsupported system": {
			os:          "win32",
			cacheHome:   "",
			userHome:    "/home/test",
			expectedErr: "unsupported system: OS 'win32'",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(tt *testing.T) {

			goOS = tc.os

			getEnv = func(key string) string {
				switch key {
				case "XDG_CACHE_HOME":
					return tc.cacheHome
				case "HOME":
					return tc.userHome
				default:
					return ""
				}
			}

			cfg, err := NewConfig("version")

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, mongoUrl, cfg.mongoUrl)
				require.Equal(t, tc.expectedCache, cfg.cachePath)
			}
		})
	}
}
