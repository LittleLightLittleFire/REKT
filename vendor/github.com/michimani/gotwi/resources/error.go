package resources

import (
	"github.com/michimani/gotwi/internal/util"
)

type Non2XXError struct {
	APIErrors     []ErrorInformation         `json:"errors"`
	Title         string                     `json:"title,omitempty"`
	Detail        string                     `json:"detail,omitempty"`
	Type          string                     `json:"type,omitempty"`
	Status        string                     `json:"-"`
	StatusCode    int                        `json:"-"`
	RateLimitInfo *util.RateLimitInformation `json:"-"`
}

type ErrorInformation struct {
	Message    string              `json:"message"`
	Code       ErrorCode           `json:"code,omitempty"`
	Label      string              `json:"label,omitempty"`
	Parameters map[string][]string `json:"parameters,omitempty"`
}

type PartialError struct {
	ResourceType *string `json:"resource_type"`
	Field        *string `json:"field"`
	Parameter    *string `json:"parameter"`
	ResourceID   *string `json:"resource_id"`
	Title        *string `json:"title"`
	Section      *string `json:"section"`
	Detail       *string `json:"detail"`
	Value        *string `json:"value"`
	Type         *string `json:"type"`
}
