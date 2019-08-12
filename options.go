package httpreq

import (
	"net/http"
	"net/url"
	"time"
)

type Proxy func(*http.Request) (*url.URL, error)

type KV struct {
	Key   interface{}
	Value interface{}
}

type Options struct {
	Header      http.Header
	Body        interface{}
	Query       url.Values
	Form        url.Values
	Cookies     []*http.Cookie
	RetryTimes  int
	Timeout     time.Duration
	MiddleWares []MiddleWare
	Debug       bool
	KVs         []*KV
}

func NewOptions() *Options {
	return &Options{}
}

type Option func(*Options)

func WithHeader(key, value string) Option {
	return func(opt *Options) {
		if opt.Header == nil {
			opt.Header = http.Header{}
		}
		opt.Header.Add(key, value)
	}
}

func WithBody(body interface{}) Option {
	return func(opt *Options) {
		opt.Body = body
	}
}

func WithQuery(key, value string) Option {
	return func(opt *Options) {
		if opt.Query == nil {
			opt.Query = url.Values{}
		}
		opt.Query.Add(key, value)
	}
}

func WithQueryValues(values url.Values) Option {
	return func(opt *Options) {
		if opt.Query == nil {
			opt.Query = url.Values{}
		}
		for key, values := range values {
			for _, value := range values {
				opt.Query.Add(key, value)
			}
		}
	}
}

func WithForm(key, value string) Option {
	return func(opt *Options) {
		if opt.Form == nil {
			opt.Form = url.Values{}
		}
		opt.Form.Add(key, value)
	}
}

func WithFormValues(values url.Values) Option {
	return func(opt *Options) {
		if opt.Form == nil {
			opt.Form = url.Values{}
		}
		for key, values := range values {
			for _, value := range values {
				opt.Form.Add(key, value)
			}
		}
	}
}

func WithCookie(name, value string) Option {
	return func(opt *Options) {
		cookie := &http.Cookie{
			Name:  name,
			Value: url.QueryEscape(value),
		}
		opt.Cookies = append(opt.Cookies, cookie)
	}
}

func WithRetryTime(retryTimes int) Option {
	return func(opt *Options) {
		opt.RetryTimes = retryTimes
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(opt *Options) {
		opt.Timeout = timeout
	}
}

func WithMiddleWares(mw MiddleWare) Option {
	return func(opt *Options) {
		opt.MiddleWares = append(opt.MiddleWares, mw)
	}

}

func WithDebug(debug bool) Option {
	return func(opt *Options) {
		opt.Debug = debug
	}
}

func WithBasicAuth(username, password string) Option {
	return WithHeader("Authorization", "Basic "+BasicAuth(username, password))
}

func WithKV(key, value string) Option {
	return func(opt *Options) {
		opt.KVs = append(opt.KVs, &KV{key, value})
	}
}
