package httpreq

import (
	"fmt"
	"net/http/httputil"
)

type EndPoint func(*Request) (*Response, error)

type MiddleWare func(EndPoint) EndPoint

func Chain(mws []MiddleWare) MiddleWare {
	if len(mws) == 0 {
		return EmptyMiddleware
	}
	return func(next EndPoint) EndPoint {
		return mws[0](Chain(mws[1:])(next))
	}
}

func EmptyMiddleware(next EndPoint) EndPoint {
	return next
}

func DebugMW(body bool) MiddleWare {
	return func(next EndPoint) EndPoint {
		return func(req *Request) (*Response, error) {
			dump, err1 := httputil.DumpRequest(req, body)
			if err1 != nil {
				fmt.Println("Error:", err1)
			} else {
				fmt.Printf("[http] HTTP Request: %s\n", string(dump))
			}

			resp, err := next(req)

			dump, err2 := httputil.DumpResponse(resp.RawResponse, false)
			if err2 != nil {
				fmt.Println("Error:", err2)
			} else {
				fmt.Printf("[http] HTTP Response: %s%s\n",
					string(dump), resp.String())
			}
			return resp, err
		}
	}
}

func ExpectStatusCodeMW(code int) MiddleWare {
	return func(next EndPoint) EndPoint {
		return func(req *Request) (*Response, error) {
			resp, err := next(req)
			if err != nil {
				return nil, err
			}
			if resp.StatusCode() != code {
				return nil, fmt.Errorf("unexpect StatusCode:%s\n    Body:%s", resp.Status(), resp.Body())
			}
			return resp, nil
		}
	}
}
