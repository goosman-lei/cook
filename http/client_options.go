package http

import (
	"net/http"
	"net/url"
	"time"
)

type HttpClientOptions struct {
	Proxy                 func(req *http.Request) (*url.URL, error)
	ConnTimeout           time.Duration
	KeepAlive             time.Duration
	MaxIdleConn           int
	IdleConnTimeout       time.Duration
	TlsHandshakeTimeout   time.Duration
	ExpectContinueTimeout time.Duration
	RequestTimeout        time.Duration
}

var (
	DefaultHttpClientOptions = HttpClientOptions{
		Proxy:                 http.ProxyFromEnvironment,
		ConnTimeout:           1 * time.Second,
		KeepAlive:             30 * time.Second,
		MaxIdleConn:           8,
		IdleConnTimeout:       90 * time.Second,
		TlsHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		RequestTimeout:        5 * time.Second,
	}
)

func MergeOptions(options ...interface{}) *HttpClientOptions {
	opts := DefaultHttpClientOptions

	// make compatible with old version
	if len(options) == 1 {
		if options[0] == nil {
			return &opts
		} else if opts_mapping, ok := options[0].(map[string]interface{}); ok {
			if v, ok := opts_mapping["proxy"]; ok {
				opts.Proxy = v.(func(req *http.Request) (*url.URL, error))
			}
			if v, ok := opts_mapping["conn_timeout"]; ok {
				opts.ConnTimeout = v.(time.Duration)
			}
			if v, ok := opts_mapping["keep_alive"]; ok {
				opts.KeepAlive = v.(time.Duration)
			}
			if v, ok := opts_mapping["max_idle_conn"]; ok {
				opts.MaxIdleConn = v.(int)
			}
			if v, ok := opts_mapping["idle_conn_timeout"]; ok {
				opts.IdleConnTimeout = v.(time.Duration)
			}
			if v, ok := opts_mapping["tls_handshake_timeout"]; ok {
				opts.TlsHandshakeTimeout = v.(time.Duration)
			}
			if v, ok := opts_mapping["expect_continue_timeout"]; ok {
				opts.ExpectContinueTimeout = v.(time.Duration)
			}
			if v, ok := opts_mapping["request_timeout"]; ok {
				opts.RequestTimeout = v.(time.Duration)
			}
			return &opts
		}
	}

	for _, opt := range options {
		opt.(func(o *HttpClientOptions))(&opts)
	}

	return &opts
}

func Proxy(proxy func(req *http.Request) (*url.URL, error)) func(o *HttpClientOptions) {
	return func(o *HttpClientOptions) {
		o.Proxy = proxy
	}
}

func ConnTimeout(conn_timeout time.Duration) func(o *HttpClientOptions) {
	return func(o *HttpClientOptions) {
		o.ConnTimeout = conn_timeout
	}
}

func KeepAlive(keep_alive time.Duration) func(o *HttpClientOptions) {
	return func(o *HttpClientOptions) {
		o.KeepAlive = keep_alive
	}
}

func MaxIdleConn(max_idle_conn int) func(o *HttpClientOptions) {
	return func(o *HttpClientOptions) {
		o.MaxIdleConn = max_idle_conn
	}
}

func IdleConnTimeout(idle_conn_timeout time.Duration) func(o *HttpClientOptions) {
	return func(o *HttpClientOptions) {
		o.IdleConnTimeout = idle_conn_timeout
	}
}

func TlsHandshakeTimeout(tls_handshake_timeout time.Duration) func(o *HttpClientOptions) {
	return func(o *HttpClientOptions) {
		o.TlsHandshakeTimeout = tls_handshake_timeout
	}
}

func ExpectContinueTimeout(expect_continue_timeout time.Duration) func(o *HttpClientOptions) {
	return func(o *HttpClientOptions) {
		o.ExpectContinueTimeout = expect_continue_timeout
	}
}

func RequestTimeout(request_timeout time.Duration) func(o *HttpClientOptions) {
	return func(o *HttpClientOptions) {
		o.RequestTimeout = request_timeout
	}
}
