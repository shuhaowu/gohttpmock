package gohttpmock

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

var defaultTransport = http.DefaultTransport

type TestResponse struct {
	StatusCode int
	Body       string
	Header     http.Header
}

func NewTestResponse(status int, body, contentType string) *TestResponse {
	return &TestResponse{
		StatusCode: status,
		Body:       body,
		Header:     http.Header{"Content-Type": {contentType}},
	}
}

type HandlerFunc func(*http.Request) *TestResponse

type RecordingTransport struct {
	Requests []*http.Request

	responses map[string]interface{}
}

func requestKey(method, url string) string {
	return method + " " + url
}

func (t *RecordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.Requests = append(t.Requests, req)
	s, ok := t.responses[requestKey(req.Method, req.URL.String())]
	var res *TestResponse
	if !ok {
		res = NewTestResponse(404, "Not Found", "text/plain")
	} else {
		switch s.(type) {
		case *TestResponse:
			res = s.(*TestResponse)
		case HandlerFunc:
			res = (s.(HandlerFunc))(req)
		}
	}

	return &http.Response{
		StatusCode: res.StatusCode,
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     res.Header,
		Body:       ioutil.NopCloser(strings.NewReader(res.Body)),
		Request:    req,
	}, nil
}

func (t *RecordingTransport) RequestBody(index int) string {
	if t.Requests[index].Body == nil {
		return ""
	}

	buf := bytes.Buffer{}
	buf.ReadFrom(t.Requests[index].Body)
	return buf.String()
}

func (t *RecordingTransport) When(method, url string) *RequestHandler {
	return &RequestHandler{method: method, url: url, recordingTransport: t}
}

type RequestHandler struct {
	method             string
	url                string
	recordingTransport *RecordingTransport
}

func (t *RequestHandler) Respond(status int, body string, contentType string) {
	t.recordingTransport.responses[requestKey(t.method, t.url)] = NewTestResponse(status, body, contentType)
}

func (t *RequestHandler) RespondFunc(f HandlerFunc) {
	t.recordingTransport.responses[requestKey(t.method, t.url)] = f
}

func StartTestHTTPCall() *RecordingTransport {
	t := &RecordingTransport{}
	t.responses = make(map[string]interface{})
	http.DefaultTransport = t
	return t
}

func EndTestHTTPCall() {
	http.DefaultTransport = defaultTransport
}
