package gotwi

import (
	"errors"
	"fmt"
	"strings"

	"github.com/michimani/gotwi/internal/gotwierrors"
	"github.com/michimani/gotwi/resources"
)

type GotwiError struct {
	err   error
	OnAPI bool
	resources.Non2XXError
}

func wrapErr(e error) *GotwiError {
	if e == nil {
		return nil
	}

	if w, ok := e.(*GotwiError); ok {
		return w
	}

	return &GotwiError{err: e}
}

func wrapWithAPIErr(n2xxerr *resources.Non2XXError) *GotwiError {
	if n2xxerr == nil {
		return nil
	}
	return &GotwiError{
		err:         errors.New(non2XXErrorSummary(n2xxerr)),
		OnAPI:       true,
		Non2XXError: *n2xxerr,
	}
}

func (e *GotwiError) Error() string {
	if e == nil {
		return ""
	}

	if e.err != nil {
		return e.err.Error()
	}

	return gotwierrors.ErrorUndefined
}

func (e *GotwiError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func non2XXErrorSummary(e *resources.Non2XXError) string {
	if e == nil {
		return ""
	}

	summary := []string{"The Twitter API returned a Response with a status other than 2XX series."}
	if e.Status != "" {
		summary = append(summary, fmt.Sprintf("httpStatus=\"%s\"", e.Status))
	}
	if e.StatusCode > 0 {
		summary = append(summary, fmt.Sprintf("httpStatusCode=%d", e.StatusCode))
	}
	if e.Title != "" {
		summary = append(summary, fmt.Sprintf("title=\"%s\"", e.Title))
	}
	if e.Detail != "" {
		summary = append(summary, fmt.Sprintf("detail=\"%s\"", e.Detail))
	}

	ercnt := 1
	for _, er := range e.APIErrors {
		if er.Message != "" && er.Code > 0 {
			detail := er.Code.Detail()
			summary = append(summary, fmt.Sprintf("errorCode%d=%d errorText%d=\"%s\" errorDescription%d=\"%s\"", ercnt, er.Code, ercnt, detail.Text, ercnt, detail.Description))
			ercnt++
		}
	}
	if e.RateLimitInfo != nil {
		summary = append(summary, fmt.Sprintf("rateLimit=%d rateLimitRemaining=%d rateLimitReset=\"%s\"", e.RateLimitInfo.Limit, e.RateLimitInfo.Remaining, e.RateLimitInfo.ResetAt))
	}

	return strings.Join(summary, " ")
}
