package gotwi

import (
	"context"
	"net/http"

	"github.com/michimani/gotwi/internal/util"
	"github.com/michimani/gotwi/resources"
)

type TypedClient[T util.Response] struct {
	Client               *http.Client
	accessToken          string
	authenticationMethod AuthenticationMethod
	oauthToken           string
	oauthConsumerKey     string
	signingKey           string
}

func NewTypedClient[T util.Response](c *Client) *TypedClient[T] {
	if c == nil {
		return nil
	}

	return &TypedClient[T]{
		Client:               c.Client,
		accessToken:          c.AccessToken(),
		authenticationMethod: c.AuthenticationMethod(),
		oauthToken:           c.OAuthToken(),
		oauthConsumerKey:     c.OAuthConsumerKey(),
		signingKey:           c.SigningKey(),
	}
}

func (c *TypedClient[T]) IsReady() bool {
	if c == nil {
		return false
	}

	if !c.AuthenticationMethod().Valid() {
		return false
	}

	switch c.AuthenticationMethod() {
	case AuthenMethodOAuth1UserContext:
		if c.OAuthToken() == "" || c.SigningKey() == "" {
			return false
		}
	case AuthenMethodOAuth2BearerToken:
		if c.AccessToken() == "" {
			return false
		}
	}

	return true
}

func (c *TypedClient[T]) Exec(req *http.Request, i util.Response) (*resources.Non2XXError, error) {
	// only satisfy IClient interface
	return nil, nil
}

func (c *TypedClient[T]) AccessToken() string {
	return c.accessToken
}

func (c *TypedClient[T]) AuthenticationMethod() AuthenticationMethod {
	return c.authenticationMethod
}

func (c *TypedClient[T]) OAuthToken() string {
	return c.oauthToken
}
func (c *TypedClient[T]) OAuthConsumerKey() string {
	return c.oauthConsumerKey
}
func (c *TypedClient[T]) SigningKey() string {
	return c.signingKey
}

func (c *TypedClient[T]) CallStreamAPI(ctx context.Context, endpoint, method string, p util.Parameters) (*StreamClient[T], error) {
	req, err := prepare(ctx, endpoint, method, p, c)
	if err != nil {
		return nil, wrapErr(err)
	}

	res, non200err, err := c.ExecStream(req)
	if err != nil {
		return nil, wrapErr(err)
	}

	if non200err != nil {
		return nil, wrapWithAPIErr(non200err)
	}

	s, err := newStreamClient[T](res)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (c *TypedClient[T]) ExecStream(req *http.Request) (*http.Response, *resources.Non2XXError, error) {
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	if _, ok := okCodes[res.StatusCode]; !ok {
		non200err, err := resolveNon2XXResponse(res)
		if err != nil {
			return nil, nil, err
		}
		return nil, non200err, nil
	}

	return res, nil, nil
}
