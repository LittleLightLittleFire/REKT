package util

import (
	"io"
	"net/url"
	"strings"
)

type Parameters interface {
	SetAccessToken(token string)
	AccessToken() string
	ResolveEndpoint(endpointBase string) string
	Body() (io.Reader, error)
	ParameterMap() map[string]string
}

func QueryValue(params []string) string {
	if len(params) == 0 {
		return ""
	}

	return strings.Join(params, ",")
}

func QueryString(paramsMap map[string]string, includes map[string]struct{}) string {
	q := url.Values{}
	for k, v := range paramsMap {
		if _, ok := includes[k]; ok {
			q.Add(k, v)
		}
	}

	return q.Encode()
}
