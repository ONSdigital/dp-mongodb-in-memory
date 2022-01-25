package download

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/spf13/afero"
	"golang.org/x/crypto/openpgp"
)

var afs = afero.Afero{Fs: afero.NewOsFs()}

// GetMongoDB ensures there is a mongodb binary in the cache path
// It will download one if not already present in the cache
func GetMongoDB(ctx context.Context, cfg Config) error {
	// Check the cache
	existsInCache, existsErr := afs.Exists(cfg.cachePath)
	if existsErr != nil {
		log.Error(ctx, "error checking cache", existsErr)
		return existsErr
	}
	if existsInCache {
		log.Info(ctx, "File found in cache", log.Data{"filename": cfg.cachePath})
		return nil
	} else {
		return downloadMongoDB(ctx, cfg)
	}
}

// downloadMongoDB will download a mongodb tarball and
// store the mongod exec file in the cache path.
// It returns the path to the saved file
func downloadMongoDB(ctx context.Context, cfg Config) error {

	downloadStartTime := time.Now()

	downloadedFile, downloadErr := downloadFile(ctx, cfg.mongoUrl)
	if downloadErr != nil {
		log.Error(ctx, "error downloading file", downloadErr, log.Data{"url": cfg.mongoUrl})
		return downloadErr
	}

	defer func() {
		_ = downloadedFile.Close()
		_ = afs.Remove(downloadedFile.Name())
	}()

	validErr := verify(ctx, cfg, downloadedFile.Name())
	if validErr != nil {
		log.Error(ctx, "error verifying integrity of MongoDB package", validErr, log.Data{"url": cfg.mongoUrl})
		return validErr
	}

	mongodTmpFile, mongoTmpErr := extractMongoBin(ctx, downloadedFile)
	if mongoTmpErr != nil {
		return mongoTmpErr
	}

	mkdirErr := afs.MkdirAll(path.Dir(cfg.cachePath), 0755)
	if mkdirErr != nil {
		log.Error(ctx, "error creating cache directory", mkdirErr, log.Data{"dir": path.Dir(cfg.cachePath)})
		return mkdirErr
	}

	renameErr := afs.Rename(mongodTmpFile, cfg.cachePath)
	if renameErr != nil {
		log.Error(ctx, "error copying mongod binary", renameErr, log.Data{"filename-from": mongodTmpFile, "filename-to": cfg.cachePath})
		return renameErr
	}

	log.Info(ctx, "mongod downloaded and stored in cache", log.Data{"filename": cfg.cachePath, "ellapsed": time.Since(downloadStartTime).String()})

	return nil
}

// downloadFile downloads the file from the given url and stores it in a temporary file.
// It returns the temporary file where it has been downloaded
func downloadFile(ctx context.Context, urlStr string) (afero.File, error) {
	log.Info(ctx, "Downloading file", log.Data{"url": urlStr})

	resp, httpGetErr := http.Get(urlStr)
	if httpGetErr != nil {
		return nil, httpGetErr
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	tgzTempFile, tmpFileErr := afs.TempFile("", "")
	if tmpFileErr != nil {
		return nil, tmpFileErr
	}

	_, copyErr := io.Copy(tgzTempFile, resp.Body)
	if copyErr != nil {
		_ = tgzTempFile.Close()
		_ = afs.Remove(tgzTempFile.Name())
		return nil, copyErr
	}

	tgzTempFile, err := afs.Open(tgzTempFile.Name())
	if err != nil {
		_ = tgzTempFile.Close()
		_ = afs.Remove(tgzTempFile.Name())
		return nil, err
	}

	log.Info(ctx, "Downloaded to temp file", log.Data{"file": tgzTempFile.Name(), "url": urlStr})
	return tgzTempFile, nil
}

// extractMongoBin extracts the mongod executable file
// from the given tarball to a temporary file.
// It returns the path to the extracted file
func extractMongoBin(ctx context.Context, tgzTempFile afero.File) (string, error) {
	_, seekErr := tgzTempFile.Seek(0, 0)
	if seekErr != nil {
		log.Error(ctx, "error seeking back to start of file", seekErr)
		return "", seekErr
	}

	gzReader, gzErr := gzip.NewReader(tgzTempFile)
	if gzErr != nil {
		log.Error(ctx, "error intializing gzip reader", gzErr, log.Data{"file": tgzTempFile.Name()})
		return "", gzErr
	}

	tarReader := tar.NewReader(gzReader)

	for {
		nextFile, tarErr := tarReader.Next()
		if tarErr == io.EOF {
			return "", fmt.Errorf("did not find a mongod binary in the tar file")
		}
		if tarErr != nil {
			log.Error(ctx, "error reading from tar file", tarErr, log.Data{"file": tgzTempFile.Name()})
			return "", tarErr
		}

		if strings.HasSuffix(nextFile.Name, "bin/mongod") {
			break
		}
	}

	// Extract to a temp file first, then copy to the destination, so we get
	// atomic behavior if there's multiple parallel downloaders
	mongodTmpFile, tmpFileErr := afs.TempFile("", "")
	if tmpFileErr != nil {
		log.Error(ctx, "error creating temp file for mongod", tmpFileErr)
		return "", tmpFileErr
	}
	defer func() {
		_ = mongodTmpFile.Close()
	}()

	_, writeErr := io.Copy(mongodTmpFile, tarReader)
	if writeErr != nil {
		log.Error(ctx, "error writing mongod binary", writeErr, log.Data{"filename": mongodTmpFile.Name()})
		return "", writeErr
	}

	_ = mongodTmpFile.Close()

	chmodErr := afs.Chmod(mongodTmpFile.Name(), 0755)
	if chmodErr != nil {
		log.Error(ctx, "error chmod-ing mongod binary", chmodErr, log.Data{"filename": mongodTmpFile.Name()})
		return "", chmodErr
	}
	return mongodTmpFile.Name(), nil
}

// verify checks the integrity of the mongoFile.
// It uses the config file to download the checksum and signature files
// and compares their value against the actual mongoFile checksum and GPG signature
func verify(ctx context.Context, cfg Config, mongoFile string) error {
	if err := verifyChecksum(ctx, cfg, mongoFile); err != nil {
		return err
	}
	log.Info(ctx, "checksum verified successfully", log.Data{"url": cfg.mongoChecksumUrl()})

	if err := verifySignature(ctx, cfg, mongoFile); err != nil {
		return err
	}
	log.Info(ctx, "signature verified successfully", log.Data{"url": cfg.mongoSignatureUrl()})

	return nil
}

func verifyChecksum(ctx context.Context, cfg Config, mongoFile string) error {

	checksumFile, downloadErr := downloadFile(ctx, cfg.mongoChecksumUrl())
	if downloadErr != nil {
		log.Error(ctx, "error downloading checksum file", downloadErr, log.Data{"url": cfg.mongoChecksumUrl()})
		return downloadErr
	}

	defer func() {
		_ = checksumFile.Close()
		_ = afs.Remove(checksumFile.Name())
	}()

	content, err := afs.ReadFile(checksumFile.Name())
	if err != nil {
		log.Error(ctx, "error reading checksum file", err)
		return err
	}
	s := strings.Split(string(content), " ")
	checksum := s[0]

	mongoChecksum, err := sha256Sum(mongoFile)
	if err != nil {
		log.Error(ctx, "error calculating SHA256 sum", err)
		return err
	}

	if checksum != mongoChecksum {
		return fmt.Errorf("checksum verification failed")
	}
	return nil
}

func verifySignature(ctx context.Context, cfg Config, mongoFilename string) error {
	// Get public key
	keyFile, err := getMongoPublicKey(ctx, cfg.mongoVersion)
	if err != nil {
		return err
	}

	defer func() {
		_ = keyFile.Close()
		_ = afs.Remove(keyFile.Name())
	}()

	keyring, err := openpgp.ReadArmoredKeyRing(keyFile)
	if err != nil {
		log.Error(ctx, "error reading keyring file", err)
		return err
	}

	// Get signature
	signatureFile, err := downloadFile(ctx, cfg.mongoSignatureUrl())
	if err != nil {
		log.Error(ctx, "error downloading signature file", err, log.Data{"url": cfg.mongoSignatureUrl()})
		return err
	}

	defer func() {
		_ = signatureFile.Close()
		_ = afs.Remove(signatureFile.Name())
	}()

	// Get file to verify
	mongoFile, err := afs.Open(mongoFilename)
	if err != nil {
		return err
	}

	defer func() {
		_ = mongoFile.Close()
	}()

	// Verify signature
	_, err = openpgp.CheckArmoredDetachedSignature(keyring, mongoFile, signatureFile)
	if err != nil {
		return fmt.Errorf("signature verification failed")
	}

	return nil
}

var getMongoPublicKey = func(ctx context.Context, version Version) (afero.File, error) {
	keyUrl := fmt.Sprintf("https://www.mongodb.org/static/pgp/server-%d.%d.asc", version.Major, version.Minor)

	keyFile, err := downloadFile(ctx, keyUrl)
	if err != nil {
		log.Error(ctx, "error downloading Mongo public key", err, log.Data{"url": keyUrl})
		return nil, err
	}
	return keyFile, nil
}

// sha256Sum returns the SHA256 checksum of the file
func sha256Sum(filename string) (string, error) {
	file, err := afs.Open(filename)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
