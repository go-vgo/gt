// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-ego/ego/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package file

// package gt

import (
	"errors"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"crypto/sha1"
	"crypto/sha256"
	"io/ioutil"
	"path/filepath"
)

// Exist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func Exist(filename string) bool {
	_, err := os.Stat(filename)
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

// Search Search a file in paths.
func Search(filename string, paths ...string) (fullpath string, err error) {
	for _, path := range paths {
		if fullpath = filepath.Join(path, filename); Exist(fullpath) {
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

// Copy copies file from source to target path.
func Copy(src, dest string) error {
	// Gather file information to set back later.
	si, err := os.Lstat(src)
	if err != nil {
		return err
	}

	// Handle symbolic link.
	if si.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(src)
		if err != nil {
			return err
		}
		// NOTE: os.Chmod and os.Chtimes don't recoganize symbolic link,
		// which will lead "no such file or directory" error.
		return os.Symlink(target, dest)
	}

	sr, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sr.Close()

	dw, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dw.Close()

	if _, err = io.Copy(dw, sr); err != nil {
		return err
	}

	// Set back file information.
	if err = os.Chtimes(dest, si.ModTime(), si.ModTime()); err != nil {
		return err
	}
	return os.Chmod(dest, si.Mode())
}

// CopyFile copies file from source to target path.
func CopyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()

	// if Exist(dst) != true {
	if !Exist(dst) {
		Write("", dst)
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()
	return io.Copy(dstFile, srcFile)
}

// OpenCopy open and copy file
func OpenCopy(srcName, dstName string) (int64, error) {
	src, err := os.Open(srcName)
	if err != nil {
		return 0, err
	}
	defer src.Close()

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

// Read read file and return string
func Read(userFile string) (string, error) {
	// userFile := fname
	fin, err := os.Open(userFile)
	if err != nil {
		log.Println("os.Open: ", userFile, err)
		return "", err
	}
	defer fin.Close()

	var str string
	buf := make([]byte, 1024)
	for {
		n, _ := fin.Read(buf)
		if 0 == n {
			break
		}
		// os.Stdout.Write(buf[:n])
		strBuf := string(buf[:n])
		str += strBuf
	}

	return str, nil
}

// WriteFile writes data to a file named by filename.
// If the file does not exist, WriteFile creates it
// and its upper level paths.
func WriteFile(fileName string, data []byte) error {
	os.MkdirAll(path.Dir(fileName), os.ModePerm)
	return ioutil.WriteFile(fileName, data, 0655)
}

// Write writes data to a file named by filename.
// If the file does not exist, WriteFile creates it
// and its upper level paths.
func Write(fileName, writeStr string) {
	os.MkdirAll(path.Dir(fileName), os.ModePerm)

	fout, err := os.Create(fileName)
	if err != nil {
		log.Println("write file "+fileName, err)
		return
	}
	defer fout.Close()

	fout.WriteString(writeStr)
}

// AppendTo append to file
func AppendTo(fileName, content string) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		log.Println("file create failed. err: " + err.Error())
		return err
	}

	n, _ := f.Seek(0, os.SEEK_END)
	_, err = f.WriteAt([]byte(content), n)

	defer f.Close()
	return err
}

// List file list
func List(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix)
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, nil
}

// ListDir dir list
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix)
	for _, fi := range dir {
		if !fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, nil
}

// Walk file walk
func Walk(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}

		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

// WalkDir dir walk
func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {

		if !fi.IsDir() {
			return nil
		}

		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}
