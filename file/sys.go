package file

import (
	"errors"
	"fmt"
	"hash"
	"io"
	"os"

	"crypto/sha1"
	"crypto/sha256"
	"path/filepath"
)

// Exist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func Exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

// PathExists path is exist
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Is returns true if path is a file,
// or returns false when it's a directory or not exist.
func Is(filePath string) bool {
	f, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return !f.IsDir()
}

// IsDir returns true if path is a directory,
// or returns false when it's a file or not exist.
func IsDir(filePath string) bool {
	f, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return f.IsDir()
}

// Mode returns file mode if file is a exist.
func Mode(filePath string) os.FileMode {
	f, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return f.Mode()
}

// Search Search a file in paths.
func Search(fileName string, paths ...string) (fullpath string, err error) {
	for _, path := range paths {
		if fullpath = filepath.Join(path, fileName); Exist(fullpath) {
			return
		}
	}
	err = errors.New(fullpath + " not found in paths")
	return
}

// Size returns file size in bytes and possible error.
func Size(file string) (int64, error) {
	f, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return f.Size(), nil
}

// MTime returns file modified time and possible error.
func MTime(file string) (int64, error) {
	f, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return f.ModTime().Unix(), nil
}

// Sha open file return sha
func Sha(filePath string, args ...string) (sha string, err error) {
	file, fsErr := os.Open(filePath)
	if fsErr != nil {
		return "", fsErr
	}
	defer file.Close()

	if len(args) > 0 {
		sha, err = IoSha(file, args[0])
		return
	}

	sha, err = IoSha(file)
	return
}

// IoSha file sha
func IoSha(fileIO *os.File, args ...string) (string, error) {
	var h hash.Hash

	if len(args) > 0 {
		h = sha256.New()
	} else {
		h = sha1.New()
	}

	_, err := io.Copy(h, fileIO)
	if err != nil {
		return "", err
	}

	sha := fmt.Sprintf("%x", h.Sum(nil))

	return sha, nil
}
