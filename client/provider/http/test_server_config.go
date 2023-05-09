package http

import (
	"time"
)

type serverTestConfig struct {
	returnHeaders map[string]string
	returnStatus  int
	returnBody    string
	sleep         time.Duration

	reqMethod  string
	reqHeaders map[string]string
	reqBody    string
	reqURL     string
}
type serverTestParam func(*serverTestConfig)

// New inits new config whith params
func (c *serverTestConfig) New(opts ...serverTestParam) serverTestConfig {
	res := *c
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(&res)
		}
	}

	return res
}

func withReturnHeaders(add map[string]string) serverTestParam {
	return func(h *serverTestConfig) {
		for k, v := range add {
			h.returnHeaders[k] = v
		}
	}
}

func withReturnStatus(v int) serverTestParam {
	return func(h *serverTestConfig) {
		h.returnStatus = v
	}
}

func withReturnBody(v string) serverTestParam {
	return func(h *serverTestConfig) {
		h.returnBody = v
	}
}

func withSleep(v time.Duration) serverTestParam {
	return func(h *serverTestConfig) {
		h.sleep = v
	}
}

func withReqHeaders(add map[string]string) serverTestParam {
	return func(h *serverTestConfig) {
		for k, v := range add {
			h.reqHeaders[k] = v
		}
	}

}

func withReqURL(url string) serverTestParam {
	return func(h *serverTestConfig) {
		h.reqURL = url
	}
}

func withReqBody(body string) serverTestParam {
	return func(h *serverTestConfig) {
		h.reqBody = body
	}
}

func withReqMethod(method string) serverTestParam {
	return func(h *serverTestConfig) {
		h.reqMethod = method
	}
}
