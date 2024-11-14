package gotwi

import (
	"bufio"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/michimani/gotwi/internal/util"
)

type StreamClient[T util.Response] struct {
	response *http.Response
	stream   *bufio.Scanner
}

func newStreamClient[T util.Response](httpRes *http.Response) (*StreamClient[T], error) {
	if httpRes == nil {
		return nil, errors.New("HTTP Response is nil.")
	}

	if httpRes.Close {
		return nil, errors.New("HTTP Response body has already closed.")
	}

	s := bufio.NewScanner(httpRes.Body)

	return &StreamClient[T]{
		response: httpRes,
		stream:   s,
	}, nil
}

func (s *StreamClient[T]) Receive() bool {
	if s == nil {
		return false
	}
	return s.stream.Scan()
}

func (s *StreamClient[T]) Stop() {
	if s == nil {
		return
	}
	s.response.Body.Close()
}

func safeUnmarshal(input []byte, target interface{}) error {
	if len(input) == 0 {
		return nil
	}
	return json.Unmarshal(input, target)
}

func (s *StreamClient[T]) Read() (T, error) {
	var n T
	if s == nil {
		return n, errors.New("StreamClient is nil.")
	}

	t := s.stream.Text()
	out := new(T)
	if err := safeUnmarshal([]byte(t), out); err != nil {
		return n, err
	}

	return *out, nil
}
