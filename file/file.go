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
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"path/filepath"
)

// Read read file and return string
func Read(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var str string
	buf := make([]byte, 1024)
	for {
		n, _ := f.Read(buf)
		if 0 == n {
			break
		}
		// os.Stdout.Write(buf[:n])
		strBuf := string(buf[:n])
		str += strBuf
	}

	return str, nil
}

// WriteFile write []byte data to a file by filename.
func WriteFile(fileName string, data []byte) error {
	err := os.MkdirAll(path.Dir(fileName), os.ModePerm)
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, data, 0655)
}

// Write write string data to a file by filename.
func Write(fileName, writeStr string) error {
	err := os.MkdirAll(path.Dir(fileName), os.ModePerm)
	if err != nil {
		return err
	}

	fout, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer fout.Close()

	_, err = fout.WriteString(writeStr)
	return err
}

// ReadIo read file return io.Reader
func ReadIo(fileName string) (io.Reader, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(f)
	return reader, err
}

// WriteIo wite file with io.Reader
func WriteIo(fileName string, fio io.Reader) error {
	data, err := io.ReadAll(fio)
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, data, 0655)
}

// ReadFromIO read io.Reader return string and error
func ReadFromIO(fio io.Reader) (string, error) {
	b, err := io.ReadAll(fio)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// AppendTo append to file
func AppendTo(fileName, content string) error {
	// write only
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	n, _ := f.Seek(0, io.SeekEnd)
	_, err = f.WriteAt([]byte(content), n)

	f.Close()
	return err
}

// Empty empty the file
func Empty(fileName string, args ...int64) error {
	var size int64
	if len(args) > 0 {
		size = args[0]
	}

	return os.Truncate(fileName, size)
}

// Move move file to new path
func Move(file, move string) error {
	return os.Rename(file, move)
}

// Remove remove the file by file name
func Remove(file string) {
	os.RemoveAll(file)
}

// List list the file
func List(dir, suffix string, isDir ...bool) (files []string, err error) {
	files = make([]string, 0, 10)
	dirIo, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	pthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix)
	for _, fi := range dirIo {
		if len(isDir) > 0 {
			if !fi.IsDir() {
				continue
			}
		} else {
			if fi.IsDir() {
				continue
			}
		}

		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, dir+pthSep+fi.Name())
		}
	}

	return files, nil
}

// ListDir list the dir
func ListDir(dir, suffix string) (files []string, err error) {
	return List(dir, suffix, false)
}

// Walk walk the file
func Walk(dir, suffix string, isDir ...bool) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)
	fn := func(filename string, fi os.FileInfo, err error) error {
		if len(isDir) > 0 {
			if !fi.IsDir() {
				return nil
			}
		} else {
			if fi.IsDir() {
				return nil
			}
		}

		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	}

	err = filepath.Walk(dir, fn)
	return files, err
}

// WalkDir walk the dir
func WalkDir(dir, suffix string) (files []string, err error) {
	return Walk(dir, suffix, false)
}

// Copy copies file from source to target path.
func Copy(src, dst string) error {
	// Get file information to set back later
	si, err := os.Lstat(src)
	if err != nil {
		return err
	}

	// Handle symbolic link
	if si.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(src)
		if err != nil {
			return err
		}
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

	// Set back file information
	err = os.Chtimes(dst, si.ModTime(), si.ModTime())
	if err != nil {
		return err
	}

	return os.Chmod(dst, si.Mode())
}

// CopyFile copies file from source to target path
func CopyFile(src, dst string) (int64, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	if !Exist(dst) {
		err := Write(dst, "")
		if err != nil {
			return 0, err
		}
	}

	dstFile, err := os.Create(dst)
	if err != nil {
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
