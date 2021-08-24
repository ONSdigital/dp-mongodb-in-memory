package download

import (
	"context"
	"net/url"
	"os"
	"path"

	"github.com/ONSdigital/log.go/v2/log"
)

const folderName = "dp-mongodb-in-memory"

// Get the mongodb download url
// Keeping this function as a variable for unit testing
var getDownloadUrl = func(version string) (string, error) {
	spec, err := MakeDownloadSpec(version)
	if err != nil {
		return "", err
	}

	return spec.GetDownloadURL()
}

// Get an environment variable
var getEnv = func(key string) string {
	return os.Getenv(key)
}

type Config struct {
	//The URL where the required mongodb tarball can be downloaded from
	mongoUrl string
	// The path where the mongod executable can be found if previously downloaded
	cachePath string
}

// Create the config values for the given version
// It will detect the OS system and
func NewConfig(version string) (*Config, error) {
	downloadUrl, err := getDownloadUrl(version)
	if err != nil {
		return nil, err
	}

	cachePath, err := buildBinCachePath(downloadUrl)
	if err != nil {
		return nil, err
	}

	return &Config{
		mongoUrl:  downloadUrl,
		cachePath: cachePath,
	}, nil
}

func buildBinCachePath(downloadUrl string) (string, error) {
	cacheHome, err := defaultBaseCachePath()
	if err != nil {
		log.Error(context.Background(), "cache directory not found", err)
		return "", err
	}

	urlParsed, err := url.Parse(downloadUrl)
	if err != nil {
		log.Error(context.Background(), "error parsing url", err, log.Data{"url": downloadUrl})
		return "", err
	}

	dirname := path.Base(urlParsed.Path)

	return path.Join(cacheHome, folderName, dirname, "mongod"), nil
}

func defaultBaseCachePath() (string, error) {
	var cacheHome = getEnv("XDG_CACHE_HOME")

	if cacheHome == "" {
		switch goOS {
		case "darwin":
			cacheHome = path.Join(getEnv("HOME"), "Library", "Caches")
		case "linux":
			cacheHome = path.Join(getEnv("HOME"), ".cache")
		default:
			return "", &UnsupportedSystemError{msg: "OS '" + goOS + "'"}
		}
	}
	return cacheHome, nil
}

// Get the url for the public signature file
func (cfg *Config) mongoSignatureUrl() string {
	return cfg.mongoUrl + ".sig"
}

// Get the url for the SHA256 file
func (cfg *Config) mongoChecksumUrl() string {
	return cfg.mongoUrl + ".sha256"
}
