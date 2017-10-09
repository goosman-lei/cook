package http

import (
	"bufio"
	"errors"
	cook_io "gitlab.niceprivate.com/golang/cook/io"
	cook_os "gitlab.niceprivate.com/golang/cook/os"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Client struct {
	hclient *http.Client
}

type ResponseInfo struct {
	*http.Response
	UseTime time.Duration
}

var (
	ErrNot200                           = errors.New("not 200 ok")
	default_opts map[string]interface{} = map[string]interface{}{
		"proxy":                   http.ProxyFromEnvironment,
		"conn_timeout":            30 * time.Second,
		"keep_alive":              30 * time.Second,
		"max_idle_conn":           100,
		"idle_conn_timeout":       90 * time.Second,
		"tls_handshake_timeout":   10 * time.Second,
		"expect_continue_timeout": 1 * time.Second,
		"request_timeout":         0 * time.Second,
	}
)

func NewClient(opts map[string]interface{}) *Client {
	opts = merge_opts(opts)

	c := &Client{}
	c.hclient = &http.Client{
		Transport: &http.Transport{
			Proxy: opts["proxy"].(func(*http.Request) (*url.URL, error)),
			DialContext: (&net.Dialer{
				Timeout:   opts["conn_timeout"].(time.Duration),
				KeepAlive: opts["keep_alive"].(time.Duration),
				DualStack: false,
			}).DialContext,
			MaxIdleConns:          opts["max_idle_conn"].(int),
			IdleConnTimeout:       opts["idle_conn_timeout"].(time.Duration),
			TLSHandshakeTimeout:   opts["tls_handshake_timeout"].(time.Duration),
			ExpectContinueTimeout: opts["expect_continue_timeout"].(time.Duration),
		},
		Timeout: opts["request_timeout"].(time.Duration),
	}

	return c
}

func merge_opts(opts map[string]interface{}) map[string]interface{} {
	if opts == nil {
		return default_opts
	}
	for opt_k, opt_v := range default_opts {
		if _, exists := opts[opt_k]; !exists {
			opts[opt_k] = opt_v
		}
	}
	return opts
}

func (c *Client) Get_file(url string, fname string) (*os.File, *ResponseInfo, error) {
	var (
		body string
		resp *ResponseInfo
		fp   *os.File
		w    *bufio.Writer
		err  error
	)

	if fp, err = cook_os.OpenFileWithMkdir(fname, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755); err != nil {
		return nil, nil, err
	}

	body, resp, err = c.shuck_response_do(func() (*http.Response, error) {
		return c.hclient.Get(url)
	})

	// have no error, write to file
	if err == nil {
		w = bufio.NewWriter(fp)
		_, err = w.WriteString(body)
	}

	return fp, resp, err
}

func (c *Client) Get(url string) (string, *ResponseInfo, error) {
	return c.shuck_response_do(func() (*http.Response, error) {
		return c.hclient.Get(url)
	})

}

func (c *Client) PostForm(url string, data url.Values) (string, *ResponseInfo, error) {
	return c.Post(url, "application/x-www-form-urlencoded", data.Encode())
}

func (c *Client) Post(url string, contentType string, data string) (string, *ResponseInfo, error) {
	return c.shuck_response_do(func() (*http.Response, error) {
		return c.hclient.Post(url, contentType, strings.NewReader(data))
	})
}

func (c *Client) Get_with_header(url string, headers http.Header) (string, *ResponseInfo, error) {
	return c.shuck_request_do(func() (http.Header, *http.Request, error) {
		if req, err := http.NewRequest(http.MethodGet, url, nil); err != nil {
			return nil, nil, err
		} else {
			return headers, req, nil
		}
	})

}

func (c *Client) PostForm_with_header(url string, data url.Values, headers http.Header) (string, *ResponseInfo, error) {
	if len(headers.Get("Content-Type")) <= 0 {
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return c.Post_with_header(url, data.Encode(), headers)
}

func (c *Client) Post_with_header(url string, data string, headers http.Header) (string, *ResponseInfo, error) {
	return c.shuck_request_do(func() (http.Header, *http.Request, error) {
		if req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data)); err != nil {
			return nil, nil, err
		} else {
			return headers, req, nil
		}
	})
}

func (c *Client) shuck_response_do(do func() (*http.Response, error)) (string, *ResponseInfo, error) {
	var (
		begin time.Time     = time.Now()
		resp  *ResponseInfo = &ResponseInfo{}
		body  string
		err   error
	)
	defer func() { resp.UseTime = time.Now().Sub(begin) }()

	if resp.Response, err = do(); err != nil {
		return body, resp, err
	}
	defer resp.Body.Close()

	body, err = cook_io.ReadAll_string(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err = ErrNot200
	}

	return body, resp, err
}

func (c *Client) shuck_request_do(do func() (http.Header, *http.Request, error)) (string, *ResponseInfo, error) {
	return c.shuck_response_do(func() (*http.Response, error) {
		headers, req, err := do()
		if err != nil {
			return nil, err
		}
		for header_name, header_values := range headers {
			for _, header_value := range header_values {
				req.Header.Add(header_name, header_value)
			}
		}
		return c.hclient.Do(req)
	})
}
