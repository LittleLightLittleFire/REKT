package util

type Response interface {
	HasPartialError() bool
}
