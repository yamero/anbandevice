package http

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

// 执行post请求
func HttpPostForm(url string, values url.Values) string {
	resp, err := http.PostForm(url, values)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

// 执行get请求
// u = "http://www.xx.com"
func HttpGet(u string) string {
	resp, err := http.Get(u)
	if err != nil {
		return "error"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}