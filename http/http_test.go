package http

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestCase_Get(t *testing.T) {
	c := NewClient(nil)
	body, resp, err := c.Get("http://www.baidu.com")
	if err != nil {
		t.Fatalf("unexcept error: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("not 200 ok")
	}

	if len(body) < 10 {
		t.Fatalf("content length too small")
	}
}

func TestCase_PostForm(t *testing.T) {
	c := NewClient(nil)
	data := url.Values{}
	data.Set("workPlace", "0/4/7/9")
	data.Set("recruitType", "2")
	data.Set("pageSize", "10")
	data.Set("curPage", "1")
	data.Set("keyWord", "php")
	data.Set("_", strconv.FormatInt(time.Now().UnixNano()/1000000, 10))
	body, resp, err := c.PostForm("http://talent.baidu.com/baidu/web/httpservice/getPostList", data)
	if err != nil {
		t.Fatalf("unexcept error: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("not 200 ok")
	}

	if len(body) < 10 {
		t.Fatalf("content length too small")
	}
}
