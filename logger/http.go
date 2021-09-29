package logger

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	protoKey     = "Proto"
	methodKey    = "Method"
	fromKey      = "From"
	urlKey       = "URL"
	bodyKey      = "Body"
	statusKey    = "Status"
	headerPrefix = "Header-"
)

func HeaderFields(h http.Header) logrus.Fields {
	m := make(map[string]interface{}, len(h)+10)
	for k, v := range h {
		if len(v) == 0 {
			continue
		}
		m[headerPrefix+k] = strings.Join(v, ";")
	}
	return m
}

func RequestFieldsBody(req *http.Request, body []byte) logrus.Fields {
	m := HeaderFields(req.Header)
	m[protoKey] = req.Proto
	m[methodKey] = req.Method
	m[fromKey] = req.RemoteAddr
	m[urlKey] = req.RequestURI
	m[bodyKey] = string(body)
	return m
}

func RequestFields(req *http.Request) logrus.Fields {
	body, _ := io.ReadAll(req.Body)
	req.Body = io.NopCloser(bytes.NewReader(body))
	return RequestFieldsBody(req, body)
}

func ResponseFieldsBody(resp *http.Response, body []byte) logrus.Fields {
	m := HeaderFields(resp.Header)
	m[protoKey] = resp.Proto
	m[statusKey] = resp.Status
	m[bodyKey] = string(body)
	return m
}

func ResponseFields(resp *http.Response) logrus.Fields {
	body, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewReader(body))
	return ResponseFieldsBody(resp, body)
}
