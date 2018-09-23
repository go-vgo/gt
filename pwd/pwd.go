// Copyright 2017 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/gt/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package pwd

import (
	"fmt"
	"strconv"
	"time"

	// "crypto/bcrypt"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// PwGenSha1 generate the password
func PwGenSha1(pass string) string {
	salt := strconv.FormatInt(time.Now().UnixNano()%9000+1000, 10)
	return Base64Encode(Sha1(Md5(pass)+salt) + salt)
}

// PwGen generate the password
func PwGen(pass string) string {
	salt := strconv.FormatInt(time.Now().UnixNano()%9000+1000, 10)
	return Base64Encode(Sha256(Md5(pass)+salt) + salt)
}

// Base64Encode base64 encode
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode base64 decode
func Base64Decode(str string) string {
	res, _ := base64.StdEncoding.DecodeString(str)
	return string(res)
}

// Bcrypt bcrypt.Sum
func Bcrypt(str string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return fmt.Sprintf("%x", hash)
}

// BcryptCheck check bcrypt hash
func BcryptCheck(encodePW, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encodePW), []byte(pwd))
	if err != nil {
		return false
	}

	return true
}

// Sha1 sha1.Sum
func Sha1(str string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(str)))
}

// Sha256 sha256.Sum
func Sha256(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}

// Sha512 sha512.Sum
func Sha512(str string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(str)))
}

// Md5 md5.Sum
func Md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

// PwCheckSha1 password check
func PwCheckSha1(pwd, saved string) bool {
	saved = Base64Decode(saved)
	if len(saved) < 4 {
		return false
	}

	salt := saved[len(saved)-4:]
	return Sha1(Md5(pwd)+salt)+salt == saved
}

// PwCheck password check
func PwCheck(pwd, saved string) bool {
	saved = Base64Decode(saved)
	if len(saved) < 4 {
		return false
	}

	salt := saved[len(saved)-4:]
	return Sha256(Md5(pwd)+salt)+salt == saved
}
