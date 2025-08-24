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

var ERROR_MALFORMED_START_LINE = fmt.Errorf("malformed start line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported HTTP Version")
var SEPARATOR = "\r\n"

func parseRequestLine(data string) (*RequestLine, string, error) {
	startLineIdx := strings.Index(data, SEPARATOR)
	if startLineIdx == -1 {
		return nil, data, nil
	}

	startLine := data[:startLineIdx]
	restOfMsg := data[startLineIdx+len(SEPARATOR):]

	parts := strings.Split(startLine, " ")
	if len(parts) < 3 {
		return nil, restOfMsg, ERROR_MALFORMED_START_LINE
	}

	httpParts := strings.Split(parts[2], "/")
	if len(parts) < 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, restOfMsg, ERROR_UNSUPPORTED_HTTP_VERSION
	}

	requestLine := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}

	return requestLine, restOfMsg, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("unable to read all"),
			err,
		)
	}

	requestLine, _, err := parseRequestLine(string(data))
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, err
}
