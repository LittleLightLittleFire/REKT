package gotwi

import "time"

func String(s string) *string {
	return &s
}

func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func Bool(b bool) *bool {
	return &b
}

func BoolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func Int(i int) *int {
	return &i
}

func IntValue(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func Float64(f float64) *float64 {
	return &f
}

func Float64Value(f *float64) float64 {
	if f == nil {
		return float64(0)
	}
	return *f
}

func Time(t time.Time) *time.Time {
	return &t
}

func TimeValue(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
