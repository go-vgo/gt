// Copyright 2016 The go-ego Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-ego/ego/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package http

import (
	"bytes"
	"io"
	"log"
	"os"
	"time"

	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
)

// Map a [string]interface{} map
type Map map[string]interface{}

// Get http get
func Get(apiUrl string, params url.Values) ([]byte, error) {
	// var Url *url.URL
	u, err := url.Parse(apiUrl)
	if err != nil {
		log.Printf("analytic url error: \r\n %v", err)
		return nil, err
	}

	// URLEncode
	u.RawQuery = params.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		log.Println("http get error: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// Post http post, params is url.Values type
func Post(apiUrl string, params url.Values, args ...int) ([]byte, error) {
	out := 1000
	if len(args) > 0 {
		out = args[0]
	}

	timeOut := time.Duration(out) * time.Millisecond

	c := &http.Client{
		Timeout: timeOut,
	}

	resp, err := c.PostForm(apiUrl, params)
	if err != nil {
		return nil, err
	}
	// fmt.Println("http:", resp)
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// API http api
func API(httpUrl string, paramMap Map, method ...string) (rs []byte, err error) {
	param := url.Values{}
	for k, v := range paramMap {
		param.Set(k, v.(string))
	}

	apiMethod := "post"
	if len(method) > 0 {
		apiMethod = method[0]
	}

	if apiMethod == "get" {
		rs, err = Get(httpUrl, param)
		return
	}

	rs, err = Post(httpUrl, param)
	return
}

// PostFile post file
func PostFile(filename, targetUrl, upParam string) (string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// uploadfile
	fileWriter, err := bodyWriter.CreateFormFile(upParam, filename)
	if err != nil {
		log.Println("error writing to buffer.")
		return "", err
	}

	// openfile
	fh, err := os.Open(filename)
	if err != nil {
		log.Println("error opening file.")
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

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Println("resp.Status is: ", resp.Status)
	// fmt.Println(string(respBody))

	return string(respBody), nil
}

var (
	userAgent = [...]string{
		"Mozilla/4.0 (compatible, MSIE 8.0, Windows NT 6.0, Trident/4.0)",
		"Mozilla/5.0 (compatible, MSIE 9.0, Windows NT 6.1, Trident/5.0,",
		"Mozilla/5.0 (Macintosh, Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
		"Mozilla/5.0 (Macintosh, U, Intel Mac OS X 10_6_8, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	}

	r = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// GetRandomUserAgent get random UserAgent
func GetRandomUserAgent(args ...[]string) string {
	if len(args) > 0 {
		return userAgent[r.Intn(len(args[0]))]
	}
	return userAgent[r.Intn(len(userAgent))]
}

// Do http.Do
func Do(url, method string, out int, args ...[]string) (*http.Response, error) {
	// POST
	// var doMethod = "GET"
	// if len(method) > 0 {
	// 	doMethod = method[0]
	// }

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println("http.NewRequest error...: ", err)
	}
	// fmt.Println("req", req)

	if len(args) > 0 {
		req.Header.Set("User-Agent", GetRandomUserAgent(args[0]))
	} else {
		req.Header.Set("User-Agent", GetRandomUserAgent())
	}

	// client := http.DefaultClient
	timeOut := time.Duration(out) * time.Millisecond
	client := &http.Client{
		Timeout: timeOut,
	}

	res, e := client.Do(req)
	if e != nil {
		log.Printf("Get request %s returned error: %s", url, e)
		return res, err
	}

	// log.Println("res... ", res)
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
