package gohttpmock

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

var defaultTransport = http.DefaultTransport

type TestResponse struct {
	// The status code of the response
	StatusCode int

	// The body of the response in a string. It will be turned into a io.ReadCloser
	Body string

	// The headers of the response
	Header http.Header
}

// Shortcut function to create a new TestResponse
func NewTestResponse(status int, body, contentType string) *TestResponse {
	return &TestResponse{
		StatusCode: status,
		Body:       body,
		Header:     http.Header{"Content-Type": {contentType}},
	}
}

type HandlerFunc func(*http.Request) *TestResponse

// The transport that replaces http.DefaultTransport
type RecordingTransport struct {
	// A list of *http.Request that is in the order of which they are requested.
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

// Gets the _i_th request's body as a string. A shortcut to reading it
// via a buffer or something.
func (t *RecordingTransport) RequestBody(i int) string {
	if t.Requests[i].Body == nil {
		return ""
	}

	buf := bytes.Buffer{}
	buf.ReadFrom(t.Requests[i].Body)
	return buf.String()
}

// Creates an expectation. Usually used in shorthand like the following:
//
//   record.When(method, url).Respond(status, body, contentType)
//
// Using it on its own is not really that useful
func (t *RecordingTransport) When(method, url string) *RequestHandler {
	return &RequestHandler{method: method, url: url, recordingTransport: t}
}

// A temporary struct for a prettier syntax.
type RequestHandler struct {
	method             string
	url                string
	recordingTransport *RecordingTransport
}

// Define the response for a particular method/url.
func (t *RequestHandler) Respond(status int, body string, contentType string) {
	t.recordingTransport.responses[requestKey(t.method, t.url)] = NewTestResponse(status, body, contentType)
}

// Define the response for a particular method/url but with a function.
func (t *RequestHandler) RespondFunc(f HandlerFunc) {
	t.recordingTransport.responses[requestKey(t.method, t.url)] = f
}

// Starts the test http call process. Disallow any real http call and record
// all of them.
//
// Calling this should follow with a defer gohttpmock.EndTestHTTPCall()
func StartTestHTTPCall() *RecordingTransport {
	t := &RecordingTransport{}
	t.responses = make(map[string]interface{})
	http.DefaultTransport = t
	return t
}

// Restores the http.DefaultTransport, allowing real http calls again.
func EndTestHTTPCall() {
	http.DefaultTransport = defaultTransport
}
