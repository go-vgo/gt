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
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"io/ioutil"
	"path/filepath"
)

// Read read file and return string
func Read(fileName string) (string, error) {
	fin, err := os.Open(fileName)
	if err != nil {
		log.Println("os.Open: ", fileName, err)
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
		log.Println("Write file "+fileName, err)
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
		log.Println("File open failed. err: " + err.Error())
		return err
	}

	n, _ := f.Seek(0, os.SEEK_END)
	_, err = f.WriteAt([]byte(content), n)

	f.Close()
	return err
}

// Empty empty the file
func Empty(fileName string, args ...int64) {
	var size int64
	if len(args) > 0 {
		size = args[0]
	}

	os.Truncate(fileName, size)
}

// List list file
func List(dirPth, suffix string) (files []string, err error) {
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

// ListDir list dir
func ListDir(dirPth, suffix string) (files []string, err error) {
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

// Walk walk file
func Walk(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)
	err = filepath.Walk(
		dirPth, func(filename string, fi os.FileInfo, err error) error {
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

// WalkDir walk dir
func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)
	err = filepath.Walk(
		dirPth, func(filename string, fi os.FileInfo, err error) error {
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

// Copy copies file from source to target path.
func Copy(src, dst string) error {
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
		return os.Symlink(target, dst)
	}

	sr, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sr.Close()

	dw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dw.Close()

	if _, err = io.Copy(dw, sr); err != nil {
		return err
	}

	// Set back file information.
	if err = os.Chtimes(dst, si.ModTime(), si.ModTime()); err != nil {
		return err
	}
	return os.Chmod(dst, si.Mode())
}

// CopyFile copies file from source to target path.
func CopyFile(src, dst string) (int64, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	// if Exist(dst) != true {
	if !Exist(dst) {
		Write(dst, "")
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		// fmt.Println(err.Error())
		return 0, err
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
