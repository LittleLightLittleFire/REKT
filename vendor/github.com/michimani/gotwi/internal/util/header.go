package util

import "net/http"

func HeaderValues(key string, h http.Header) []string {
	if hv, ok := h[key]; ok {
		return hv
	}

	return []string{}
}
