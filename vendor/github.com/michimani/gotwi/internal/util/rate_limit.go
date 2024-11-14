package util

import (
	"net/http"
	"strconv"
	"time"
)

const (
	RATE_LIMIT_LIMIT_HEADER_KEY     = "X-Rate-Limit-Limit"
	RATE_LIMIT_REMAINING_HEADER_KEY = "X-Rate-Limit-Remaining"
	RATE_LIMIT_RESET_HEADER_KEY     = "X-Rate-Limit-Reset"
)

type RateLimitInformation struct {
	Limit     int
	Remaining int
	ResetAt   *time.Time
}

func GetRateLimitInformation(res *http.Response) (*RateLimitInformation, error) {
	i := RateLimitInformation{}
	limit, err := rateLimitLimit(res.Header)
	if err != nil {
		return nil, err
	}
	i.Limit = limit

	remaining, err := rateLimitRemaining(res.Header)
	if err != nil {
		return nil, err
	}
	i.Remaining = remaining

	reset, err := rateLimitResetAt(res.Header)
	if err != nil {
		return nil, err
	}
	i.ResetAt = reset

	return &i, nil
}

func rateLimitLimit(h http.Header) (int, error) {
	values := HeaderValues(RATE_LIMIT_LIMIT_HEADER_KEY, h)
	if len(values) == 0 {
		return 0, nil
	}
	limit, err := strconv.Atoi(values[0])
	if err != nil {
		return 0, err
	}

	return limit, nil
}

func rateLimitRemaining(h http.Header) (int, error) {
	values := HeaderValues(RATE_LIMIT_REMAINING_HEADER_KEY, h)
	if len(values) == 0 {
		return 0, nil
	}
	remaining, err := strconv.Atoi(values[0])
	if err != nil {
		return 0, err
	}

	return remaining, nil
}

func rateLimitResetAt(h http.Header) (*time.Time, error) {
	values := HeaderValues(RATE_LIMIT_RESET_HEADER_KEY, h)
	if len(values) == 0 {
		return nil, nil
	}
	resetInt, err := strconv.Atoi(values[0])
	if err != nil {
		return nil, err
	}

	reset := time.Unix(int64(resetInt), 0)
	return &reset, nil
}
