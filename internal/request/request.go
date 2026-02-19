package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

var MALFORMED_REQUEST_LINE = fmt.Errorf("MALFORMED_REQUEST_LINE")
var UNSUPPORTED_FORMAT = fmt.Errorf("UNSUPPORTED_FORMAT")
var SEPARATOR = "\r\n"

func (r *RequestLine) validHttp() bool {
	return r.HttpVersion == "HTTP/1.1"
}

func parseRequestLine(b string) (*RequestLine, string, error) {
	i := strings.Index(b, SEPARATOR)
	if i == -1 {
		return nil, "", nil
	}
	startLine := b[:i]
	restOfMsg := b[i+len(SEPARATOR):]

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, restOfMsg, MALFORMED_REQUEST_LINE
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2],
	}

	if !rl.validHttp() {
		return nil, restOfMsg, UNSUPPORTED_FORMAT
	}
	return rl, restOfMsg, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("error reading request body: %w", err), err)
	}

	str := string(data)
	rl, _, err := parseRequestLine(str)
	if err != nil {
		return nil, err
	}

	return &Request{RequestLine: *rl}, nil
}
