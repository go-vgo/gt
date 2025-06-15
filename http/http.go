// Copyright 2016 The go-ego Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-ego/ego/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

package http

import (
	"bytes"
	"io"
	"os"
	"time"

	"mime/multipart"
	"net/http"
	"net/url"
)

// Map a [string]interface{} map
type Map map[string]interface{}

// Get http get
func Get(api string, args ...url.Values) ([]byte, error) {
	var params url.Values
	if len(args) > 0 {
		params = args[0]
	}

	u, err := url.Parse(api)
	if err != nil {
		return nil, err
	}

	// URLEncode
	u.RawQuery = params.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// Post http post, params is url.Values type
func Post(api string, args ...interface{}) ([]byte, error) {
	var params url.Values
	if len(args) > 0 {
		params = args[0].(url.Values)
	}

	out := 1000
	if len(args) > 1 {
		out = args[1].(int)
	}

	timeOut := time.Duration(out) * time.Millisecond
	c := &http.Client{
		Timeout: timeOut,
	}

	resp, err := c.PostForm(api, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// Api http api
func Api(api string, args ...interface{}) (rs []byte, err error) {
	paramMap := Map{}
	if len(args) > 0 {
		paramMap = args[0].(Map)
	}

	param := url.Values{}
	for k, v := range paramMap {
		param.Set(k, v.(string))
	}

	apiMethod := "post"
	if len(args) > 1 {
		apiMethod = args[1].(string)
	}

	if apiMethod == "get" {
		rs, err = Get(api, param)
		return
	}

	rs, err = Post(api, param)
	return
}

// Do http.Do
func Do(url, method string, out int, args ...[]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", GetRandomUserAgent(args...))

	// client := http.DefaultClient
	timeOut := time.Duration(out) * time.Millisecond
	client := &http.Client{
		Timeout: timeOut,
	}

	res, e := client.Do(req)
	if e != nil {
		return res, e
	}

	return res, nil
}

// DoPost http.Do post
func DoPost(url string, args ...interface{}) (*http.Response, error) {
	var out int
	if len(args) > 0 {
		out = args[0].(int)
	}

	if len(args) > 1 {
		res, err := Do(url, "POST", out, args[1].([]string))
		return res, err
	}

	res, err := Do(url, "POST", out)
	return res, err
}

// DoGet http.Do get
func DoGet(url string, args ...interface{}) (*http.Response, error) {
	var out int
	if len(args) > 0 {
		out = args[0].(int)
	}

	if len(args) > 1 {
		res, err := Do(url, "GET", out, args[1].([]string))
		return res, err
	}

	res, err := Do(url, "GET", out)
	return res, err
}

// PostFile post file
func PostFile(filename, targetUrl, upParam string) (string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// uploadfile
	fileWriter, err := bodyWriter.CreateFormFile(upParam, filename)
	if err != nil {
		return "", err
	}

	// openfile
	fh, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	// iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return "", err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
