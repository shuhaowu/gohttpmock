gohttpmock
==========

Mocking HTTP calls in Go. 

[![Travis Status](https://travis-ci.org/shuhaowu/gohttpmock.svg)][travis]

[travis]: https://travis-ci.org/shuhaowu/gohttpmock

Examples
--------

Here's a simple use case, with test written with [gocheck][gc].

[gc]: http://labix.org/gocheck

    func (*Suite) TestRequest(c *C) {
        // Overrides http.DefaultTransport with a recording transport.
        record := gohttpmock.StartTestHTTPCall()
        // Makes sure that we clean up and don't pollute state with other
        // tests. Also restores http.DefaultTransport
        defer gohttpmock.EndTestHTTPCall()

        // Defining an URL.
        record.When("POST", "http://example.com/").Respond(200, "body", "text/plain")

        resp, err := http.Post("http://example.com/", "request_body", "text/plain")

        // no error
        c.Assert(err, IsNil)

        // Read the response body
        buf := &bytes.Buffer{}
        buf.ReadFrom(resp.Body)
        respBody := buf.String()

        // resp.Body reads to "body" as we defined.
        c.Assert(respBody, Equals, "body")

        // record.Requests contains a list of *http.Request in the order of
        // which they are recorded. It currently has a length of 1.
        c.Assert(record.Requests, HasLen, 1)

        // All the information from the original request is obtained.
        c.Assert(record.Requests[0].URL.String(), Equals, "http://example.com/")
        c.Assert(record.Requests[0].Method, Equals, "POST")

        // func RequestBody(i int) string is just a convenience function for 
        // getting the request body to be a string, at index i
        c.Assert(record.RequestBody(0), Equals, "request_body")
    }

You can also generate a response on demand:
    
    i := 0
    f := func(req *http.Request) *gohttpmock.TestResponse {
        i++
        return gohttpmock.NewTestResponse(200, fmt.Sprintf("body %d", i), "text/plain")
    }
    record.When("POST", "http://example.com/").RespondFunc(f)

    // body 1
    resp, err = http.Post("http://example.com/", "request_body", "text/plain")
    // body 2
    resp, err = http.Post("http://example.com/", "request_body", "text/plain")

If you do not define a route, any request while `StartTestHTTPCall` is called
and `EndTestHTTPCall` is not will result in a 404.

For more information, see the [documentation][doc].

[doc]: http://godoc.org/github.com/shuhaowu/gohttpmock