package gotwi

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	OAuthVersion10               = "1.0"
	OAuthSignatureMethodHMACSHA1 = "HMAC-SHA1"
)

type Endpoint string

type EndpointInfo struct {
	Raw                      string
	Base                     string
	EncodedQueryParameterMap map[string]string
}

type CreateOAuthSignatureInput struct {
	HTTPMethod       string
	RawEndpoint      string
	OAuthConsumerKey string
	OAuthToken       string
	SigningKey       string
	ParameterMap     map[string]string
}

type CreateOAuthSignatureOutput struct {
	OAuthNonce           string
	OAuthSignatureMethod string
	OAuthTimestamp       string
	OAuthVersion         string
	OAuthSignature       string
}

func CreateOAuthSignature(in *CreateOAuthSignatureInput) (*CreateOAuthSignatureOutput, error) {
	out := CreateOAuthSignatureOutput{
		OAuthSignatureMethod: OAuthSignatureMethodHMACSHA1,
		OAuthVersion:         OAuthVersion10,
	}
	nonce, err := generateOAthNonce()
	if err != nil {
		return nil, err
	}
	out.OAuthNonce = nonce

	ts := fmt.Sprintf("%d", time.Now().Unix())
	out.OAuthTimestamp = ts
	endpointBase := endpointBase(in.RawEndpoint)

	parameterString := createParameterString(nonce, ts, in)
	sigBase := createSignatureBase(in.HTTPMethod, endpointBase, parameterString)
	sig, err := calculateSignature(sigBase, in.SigningKey)
	if err != nil {
		return nil, err
	}
	out.OAuthSignature = sig

	return &out, nil
}

func generateOAthNonce() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	nonce := base64.StdEncoding.EncodeToString(key)
	symbols := []string{"+", "/", "="}
	for _, s := range symbols {
		nonce = strings.Replace(nonce, s, "", -1)
	}
	return nonce, nil
}

func endpointBase(e string) string {
	queryIdx := strings.Index(e, "?")
	if queryIdx < 0 {
		return e
	}

	return e[:queryIdx]
}

func (e Endpoint) String() string {
	return string(e)
}

func (e Endpoint) Detail() (*EndpointInfo, error) {
	d := EndpointInfo{
		Raw:                      e.String(),
		EncodedQueryParameterMap: map[string]string{},
	}

	queryIdx := strings.Index(e.String(), "?")
	if queryIdx < 0 {
		d.Base = string(e)
		return &d, nil
	}

	d.Base = e.String()[:queryIdx]
	queryPart := e.String()[queryIdx+1:]
	paramsPairs := strings.Split(queryPart, "&")
	for _, pp := range paramsPairs {
		keyValue := strings.Split(pp, "=")
		var err error
		v := ""
		if len(keyValue) == 2 {
			v, err = url.QueryUnescape(keyValue[1])
			if err != nil {
				return nil, err
			}
		}
		d.EncodedQueryParameterMap[keyValue[0]] = v
	}

	return &d, nil
}

func createParameterString(nonce, ts string, in *CreateOAuthSignatureInput) string {
	qv := url.Values{}
	for k, v := range in.ParameterMap {
		qv.Add(k, v)
	}

	qv.Add("oauth_consumer_key", in.OAuthConsumerKey)
	qv.Add("oauth_nonce", nonce)
	qv.Add("oauth_signature_method", OAuthSignatureMethodHMACSHA1)
	qv.Add("oauth_timestamp", ts)
	qv.Add("oauth_token", in.OAuthToken)
	qv.Add("oauth_version", OAuthVersion10)

	encoded := qv.Encode()
	encoded = regexp.MustCompile(`([^%])(\+)`).ReplaceAllString(encoded, "$1%20")
	return encoded
}

func createSignatureBase(method, endpointBase, parameterString string) string {
	return fmt.Sprintf(
		"%s&%s&%s",
		url.QueryEscape(strings.ToUpper(method)),
		url.QueryEscape(endpointBase),
		url.QueryEscape(parameterString),
	)
}

func calculateSignature(base, key string) (string, error) {
	b := []byte(key)
	h := hmac.New(sha1.New, b)
	_, err := io.WriteString(h, base)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
