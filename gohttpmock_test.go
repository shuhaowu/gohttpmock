package gohttpmock

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
)

const (
	exampleCom = "http://example.com/"
	fakeDomain = "http://lkqjwerklniofqklrkjaakljwqrnlkjaoijfqklnclkajfoqer.none"
)

func readBody(body io.Reader) string {
	buf := &bytes.Buffer{}
	buf.ReadFrom(body)
	return buf.String()
}

func TestRespond(t *testing.T) {
	record := StartTestHTTPCall()
	defer EndTestHTTPCall()
	record.When("GET", exampleCom).Respond(200, "body", "text/plain")
	resp, err := http.Get(exampleCom)

	if err != nil {
		t.Fatalf("err should be nil but is %q", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("resp.StatusCode is not '200' but is %q", resp.StatusCode)
	}

	respBody := readBody(resp.Body)
	if respBody != "body" {
		t.Fatalf("resp.Body is not 'body' but is %q", respBody)
	}

	if len(record.Requests) != 1 {
		t.Fatalf("record.Requests should only have 1 item but has %d", len(record.Requests))
	}

	if record.Requests[0].Method != "GET" {
		t.Fatalf("record.Requests[0].Method is not GET but is %q", record.Requests[0].Method)
	}

	if record.Requests[0].URL.String() != exampleCom {
		t.Fatalf("record.Requests[0].URL.String() is not example.com but is %q", record.Requests[0].URL.String())
	}

	if record.RequestBody(0) != "" {
		t.Fatalf("record.RequestBody(0) is not '' but is %q", record.RequestBody(0))
	}

	resp, err = http.Post(exampleCom, "text/plain", bytes.NewBufferString("reqbody"))

	if err != nil {
		t.Fatalf("err should be nil but is %q", err)
	}

	if resp.StatusCode != 404 {
		t.Fatalf("resp.StatusCode should be 404 but is %q", resp.StatusCode)
	}

	if len(record.Requests) != 2 {
		t.Fatalf("record.Requests should only have 2 item but has %d", len(record.Requests))
	}

	if record.Requests[1].Method != "POST" {
		t.Fatalf("record.Requests[1].Method is not POST but is %q", record.Requests[0].Method)
	}

	if record.Requests[1].URL.String() != exampleCom {
		t.Fatalf("record.Requests[1].URL.String() is not example.com but is %q", record.Requests[0].URL.String())
	}

	if record.RequestBody(1) != "reqbody" {
		t.Fatalf("record.RequestBody(1) is not 'reqbody' but is %q", record.RequestBody(0))
	}

	record.When("POST", exampleCom).Respond(200, "postedbody", "text/plain")

	resp, err = http.Post(exampleCom, "text/plain", bytes.NewBufferString("reqbody"))

	if err != nil {
		t.Fatalf("err should be nil but is %q", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("resp.StatusCode should be 200 but is %q", resp.StatusCode)
	}

	respBody = readBody(resp.Body)
	if respBody != "postedbody" {
		t.Fatalf("resp.Body is not 'postedbody' but is %q", respBody)
	}

	if len(record.Requests) != 3 {
		t.Fatalf("record.Requests should only have 3 item but has %d", len(record.Requests))
	}

	if record.Requests[2].Method != "POST" {
		t.Fatalf("record.Requests[2].Method is not POST but is %q", record.Requests[0].Method)
	}

	if record.Requests[2].URL.String() != exampleCom {
		t.Fatalf("record.Requests[2].URL.String() is not example.com but is %q", record.Requests[0].URL.String())
	}

	if record.RequestBody(2) != "reqbody" {
		t.Fatalf("record.RequestBody(2) is not 'reqbody' but is %q", record.RequestBody(0))
	}
}

func TestRespondFunc(t *testing.T) {
	record := StartTestHTTPCall()
	defer EndTestHTTPCall()

	i := 0
	f := func(req *http.Request) *TestResponse {
		i++
		return NewTestResponse(200, fmt.Sprintf("body %d", i), "text/plain")
	}

	record.When("GET", exampleCom).RespondFunc(f)

	resp, err := http.Get(exampleCom)

	if err != nil {
		t.Fatalf("err should be nil but is %q", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("resp.StatusCode is not '200' but is %q", resp.StatusCode)
	}

	respBody := readBody(resp.Body)
	if respBody != "body 1" {
		t.Fatalf("resp.Body is not 'body 1' but is %q", respBody)
	}

	if len(record.Requests) != 1 {
		t.Fatalf("record.Requests should only have 1 item but has %d", len(record.Requests))
	}

	if record.Requests[0].Method != "GET" {
		t.Fatalf("record.Requests[0].Method is not GET but is %q", record.Requests[0].Method)
	}

	if record.Requests[0].URL.String() != exampleCom {
		t.Fatalf("record.Requests[0].URL.String() is not example.com but is %q", record.Requests[0].URL.String())
	}

	if record.RequestBody(0) != "" {
		t.Fatalf("record.RequestBody(0) is not '' but is %q", record.RequestBody(0))
	}

	resp, err = http.Get(exampleCom)

	if err != nil {
		t.Fatalf("err should be nil but is %q", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("resp.StatusCode is not '200' but is %q", resp.StatusCode)
	}

	respBody = readBody(resp.Body)
	if respBody != "body 2" {
		t.Fatalf("resp.Body is not 'body 2' but is %q", respBody)
	}

	if len(record.Requests) != 2 {
		t.Fatalf("record.Requests should only have 2 item but has %d", len(record.Requests))
	}

	if record.Requests[1].Method != "GET" {
		t.Fatalf("record.Requests[1].Method is not GET but is %q", record.Requests[0].Method)
	}

	if record.Requests[1].URL.String() != exampleCom {
		t.Fatalf("record.Requests[1].URL.String() is not example.com but is %q", record.Requests[0].URL.String())
	}

	if record.RequestBody(1) != "" {
		t.Fatalf("record.RequestBody(0) is not '' but is %q", record.RequestBody(0))
	}
}

func TestEndTestHTTPCall(t *testing.T) {
	_ = StartTestHTTPCall()
	EndTestHTTPCall()

	_, err := http.Get(fakeDomain)

	if err == nil {
		t.Fatalf("Either %q is a real place now, or this test just failed", fakeDomain)
	}
}
