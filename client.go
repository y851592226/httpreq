package httpreq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	*http.Client
	ops []Option
}

type CheckRedirect func(*http.Request, []*http.Request) error
type DialContext func(ctx context.Context, network, addr string) (net.Conn, error)
type Dial func(network, addr string) (net.Conn, error)
type DialTLS func(network, addr string) (net.Conn, error)

func (c *Client) SetTransport(transport http.RoundTripper) {
	c.Client.Transport = transport
}

func (c *Client) SetCheckRedirect(checkRedirect CheckRedirect) {
	c.Client.CheckRedirect = checkRedirect
}

func (c *Client) SetCookieJar(jar http.CookieJar) {
	c.Client.Jar = jar
}

func (c *Client) SetProxyURL(URL string) error {
	proxyURL, err := url.Parse(URL)
	if err != nil {
		return err
	}
	proxy := http.ProxyURL(proxyURL)
	return c.SetProxy(proxy)
}

func (c *Client) SetProxy(proxy Proxy) error {
	t, ok := c.Transport.(*http.Transport)
	if !ok {
		return fmt.Errorf("unsupport set poroxy Transport:%T", c.Transport)
	}
	t.Proxy = proxy
	return nil
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.Client.Timeout = timeout
}

func (c *Client) SetDial(dial Dial) error {
	t, ok := c.Transport.(*http.Transport)
	if !ok {
		return fmt.Errorf("unsupport set Dial Transport:%T", c.Transport)
	}
	t.Dial = dial //nolint
	return nil
}

func (c *Client) SetDialContext(dialContext DialContext) error {
	t, ok := c.Transport.(*http.Transport)
	if !ok {
		return fmt.Errorf("unsupport set DialContext Transport:%T", c.Transport)
	}
	t.DialContext = dialContext
	return nil
}

func (c *Client) SetDialTLS(dialTLS DialTLS) error {
	t, ok := c.Transport.(*http.Transport)
	if !ok {
		return fmt.Errorf("unsupport set DialTLS Transport:%T", c.Transport)
	}
	t.DialTLS = dialTLS
	return nil
}

func NewClient(client *http.Client, ops ...Option) *Client {
	return &Client{client, ops}
}

var DefaultClient = NewClient(&http.Client{Transport: http.DefaultTransport})

func Get(url string, ops ...Option) (*Response, error) {
	return DefaultClient.Get(url, ops...)
}

func (c *Client) Get(url string, ops ...Option) (*Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, ops...)
}

func Post(url, contentType string, body io.Reader, ops ...Option) (*Response, error) {
	return DefaultClient.Post(url, contentType, body, ops...)
}

func (c *Client) Post(url, contentType string, body io.Reader, ops ...Option) (*Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req, ops...)
}

func PostForm(url string, data url.Values, ops ...Option) (*Response, error) {
	return DefaultClient.PostForm(url, data, ops...)
}

func (c *Client) PostForm(url string, data url.Values, ops ...Option) (*Response, error) {
	return c.Post(url, "application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()), ops...)
}

func Head(url string, ops ...Option) (*Response, error) {
	return DefaultClient.Head(url, ops...)
}

func (c *Client) Head(url string, ops ...Option) (*Response, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, ops...)
}

func (c *Client) Do(req *Request, ops ...Option) (*Response, error) {
	// set request options
	opt := NewOptions()
	for _, op := range c.ops {
		op(opt)
	}
	for _, op := range ops {
		op(opt)
	}

	// set request context
	ctx := req.Context()

	for _, kv := range opt.KVs {
		ctx = context.WithValue(ctx, kv.Key, kv.Value)
	}
	if opt.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opt.Timeout)
		defer cancel()
	}
	req = req.WithContext(ctx)

	// set request header
	if opt.Header != nil {
		for k, v := range opt.Header {
			req.Header.Set(k, v[0])
		}
	}

	// set request form
	if opt.Form != nil {
		opt.Body = strings.NewReader(opt.Form.Encode())
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	// set request body
	if opt.Body != nil {
		switch v := opt.Body.(type) {
		case *bytes.Buffer:
			req.ContentLength = int64(v.Len())
			buf := v.Bytes()
			req.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return ioutil.NopCloser(r), nil
			}
		case *bytes.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		case *strings.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		case string:
			req.ContentLength = int64(len(v))
			req.GetBody = func() (io.ReadCloser, error) {
				return ioutil.NopCloser(strings.NewReader(v)), nil
			}
		case []byte:
			req.ContentLength = int64(len(v))
			req.GetBody = func() (io.ReadCloser, error) {
				return ioutil.NopCloser(bytes.NewReader(v)), nil
			}
		default:
			data, err := json.Marshal(opt.Body)
			if err != nil {
				return nil, err
			}
			req.ContentLength = int64(len(data))
			req.Header.Set("Content-Type", "application/json")
			req.GetBody = func() (io.ReadCloser, error) {
				return ioutil.NopCloser(bytes.NewReader(data)), nil
			}
		}
		if req.GetBody != nil && req.ContentLength == 0 {
			req.Body = http.NoBody
			req.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
		} else {
			var err error
			req.Body, err = req.GetBody()
			if err != nil {
				return nil, err
			}
		}
	}

	// set request query
	if opt.Query != nil {
		req.URL.RawQuery = opt.Query.Encode()
	}

	// set request cookie
	if opt.Cookies != nil {
		for _, cookie := range opt.Cookies {
			req.AddCookie(cookie)
		}
	}

	if opt.Debug {
		opt.MiddleWares = append(opt.MiddleWares, DebugMW(true))
	}
	endpoint := Chain(opt.MiddleWares)(c.do)

	// do request
	var resp *Response
	var err error
	for i := 0; i <= opt.RetryTimes; i++ {
		resp, err = endpoint(req)
		if err != nil {
			continue
		}
		return resp, err
	}
	return nil, err
}

func (c *Client) do(req *Request) (*Response, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return newResponse(req, resp)
}
