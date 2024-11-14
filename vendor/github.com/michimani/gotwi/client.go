package gotwi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/michimani/gotwi/internal/gotwierrors"
	"github.com/michimani/gotwi/internal/util"
	"github.com/michimani/gotwi/resources"
)

const (
	APIKeyEnvName       = "GOTWI_API_KEY"
	APIKeySecretEnvName = "GOTWI_API_KEY_SECRET"
)

type AuthenticationMethod string

const (
	AuthenMethodOAuth1UserContext = "OAuth 1.0a User context"
	AuthenMethodOAuth2BearerToken = "OAuth 2.0 Bearer token"
)

func (a AuthenticationMethod) Valid() bool {
	return a == AuthenMethodOAuth1UserContext || a == AuthenMethodOAuth2BearerToken
}

type NewClientInput struct {
	HTTPClient           *http.Client
	AuthenticationMethod AuthenticationMethod
	OAuthToken           string
	OAuthTokenSecret     string
	APIKey               string
	APIKeySecret         string
	Debug                bool
}

type NewClientWithAccessTokenInput struct {
	HTTPClient  *http.Client
	AccessToken string
}

type IClient interface {
	Exec(req *http.Request, i util.Response) (*resources.Non2XXError, error)
	IsReady() bool
	AccessToken() string
	AuthenticationMethod() AuthenticationMethod
	OAuthToken() string
	OAuthConsumerKey() string
	SigningKey() string
}

type Client struct {
	Client               *http.Client
	authenticationMethod AuthenticationMethod
	accessToken          string
	oauthToken           string
	oauthConsumerKey     string
	signingKey           string
	apiKeyOverride       string
	apiKeySecretOverride string
	debug                bool
}

type ClientResponse struct {
	StatusCode int
	Status     string
	Error      *resources.Non2XXError
	Body       []byte
	Response   util.Response
}

var defaultHTTPClient = &http.Client{
	Timeout: time.Duration(30) * time.Second,
}

func NewClient(in *NewClientInput) (*Client, error) {
	if in == nil {
		return nil, fmt.Errorf("NewClientInput is nil.")
	}

	if !in.AuthenticationMethod.Valid() {
		return nil, fmt.Errorf("AuthenticationMethod is invalid.")
	}

	c := Client{
		Client:               defaultHTTPClient,
		authenticationMethod: in.AuthenticationMethod,
		apiKeyOverride:       in.APIKey,
		apiKeySecretOverride: in.APIKeySecret,
		debug:                in.Debug,
	}

	if in.HTTPClient != nil {
		c.Client = in.HTTPClient
	}

	if err := c.authorize(in.OAuthToken, in.OAuthTokenSecret); err != nil {
		return nil, err
	}

	return &c, nil
}

func NewClientWithAccessToken(in *NewClientWithAccessTokenInput) (*Client, error) {
	if in == nil {
		return nil, fmt.Errorf("NewClientWithAccessTokenInput is nil.")
	}

	if in.AccessToken == "" {
		return nil, fmt.Errorf("AccessToken is empty.")
	}

	c := Client{
		Client:               defaultHTTPClient,
		authenticationMethod: AuthenMethodOAuth2BearerToken,
		accessToken:          in.AccessToken,
	}

	if in.HTTPClient != nil {
		c.Client = in.HTTPClient
	}

	return &c, nil
}

func (c *Client) authorize(oauthToken, oauthTokenSecret string) error {
	apiKey := c.APIKey()
	apiKeySecret := c.APIKeySecret()
	if apiKey == "" || apiKeySecret == "" {
		return fmt.Errorf("env '%s' and '%s' is required.", APIKeyEnvName, APIKeySecretEnvName)
	}
	c.oauthConsumerKey = apiKey

	switch c.AuthenticationMethod() {
	case AuthenMethodOAuth1UserContext:
		if oauthToken == "" || oauthTokenSecret == "" {
			return fmt.Errorf("OAuthToken and OAuthTokenSecret is required for using %s.", AuthenMethodOAuth1UserContext)
		}

		c.oauthToken = oauthToken
		c.signingKey = fmt.Sprintf("%s&%s",
			url.QueryEscape(apiKeySecret),
			url.QueryEscape(oauthTokenSecret))
	case AuthenMethodOAuth2BearerToken:
		accessToken, err := GenerateBearerToken(c, apiKey, apiKeySecret)
		if err != nil {
			return err
		}

		c.accessToken = accessToken
	}

	return nil
}

func (c *Client) IsReady() bool {
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

func (c *Client) AccessToken() string {
	return c.accessToken
}

func (c *Client) AuthenticationMethod() AuthenticationMethod {
	return c.authenticationMethod
}

func (c *Client) APIKey() string {
	if c.apiKeyOverride != "" {
		return c.apiKeyOverride
	}
	return os.Getenv(APIKeyEnvName)
}

func (c *Client) APIKeySecret() string {
	if c.apiKeySecretOverride != "" {
		return c.apiKeySecretOverride
	}
	return os.Getenv(APIKeySecretEnvName)
}

func (c *Client) OAuthToken() string {
	return c.oauthToken
}
func (c *Client) OAuthConsumerKey() string {
	return c.oauthConsumerKey
}
func (c *Client) SigningKey() string {
	return c.signingKey
}

func (c *Client) SetAccessToken(v string) {
	c.accessToken = v
}

func (c *Client) SetAuthenticationMethod(v AuthenticationMethod) {
	c.authenticationMethod = v
}

func (c *Client) SetOAuthToken(v string) {
	c.oauthToken = v
}
func (c *Client) SetOAuthConsumerKey(v string) {
	c.oauthConsumerKey = v
}
func (c *Client) SetSigningKey(v string) {
	c.signingKey = v
}

func (c *Client) CallAPI(ctx context.Context, endpoint, method string, p util.Parameters, i util.Response) error {
	req, err := prepare(ctx, endpoint, method, p, c)
	if err != nil {
		return wrapErr(err)
	}

	non200err, err := c.Exec(req, i)
	if err != nil {
		return wrapErr(err)
	}

	if non200err != nil {
		return wrapWithAPIErr(non200err)
	}

	return nil
}

var okCodes map[int]struct{} = map[int]struct{}{
	http.StatusOK:      {},
	http.StatusCreated: {},
}

func (c *Client) Exec(req *http.Request, i util.Response) (*resources.Non2XXError, error) {
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if _, ok := okCodes[res.StatusCode]; !ok {
		non200err, err := resolveNon2XXResponse(res)
		if err != nil {
			return nil, err
		}
		return non200err, nil
	}

	var tr io.Reader
	debugBuf := new(bytes.Buffer)
	if c.debug {
		tr = io.TeeReader(res.Body, debugBuf)
	} else {
		tr = res.Body
	}

	jerr := json.NewDecoder(tr).Decode(i)
	if c.debug {
		fmt.Printf("------DEBUG------\n[request url]\n%v\n[response header]\n%v\n[response body]\n%s\n------DEBUG END------\n", req.URL, res.Header, debugBuf.String())
	}
	if jerr != nil && jerr != io.EOF {
		return nil, jerr
	}

	return nil, nil
}

func prepare(ctx context.Context, endpointBase, method string, p util.Parameters, c IClient) (*http.Request, error) {
	if p == nil {
		return nil, fmt.Errorf(gotwierrors.ErrorParametersNil, endpointBase)
	}

	if !c.IsReady() {
		return nil, fmt.Errorf(gotwierrors.ErrorClientNotReady)
	}

	endpoint := p.ResolveEndpoint(endpointBase)
	p.SetAccessToken(c.AccessToken())
	req, err := newRequest(ctx, endpoint, method, p)
	if err != nil {
		return nil, err
	}

	switch c.AuthenticationMethod() {
	case AuthenMethodOAuth1UserContext:
		pm := p.ParameterMap()
		req, err = setOAuth1Header(req, pm, c)
		if err != nil {
			return nil, err
		}
	case AuthenMethodOAuth2BearerToken:
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.AccessToken()))
	}

	return req, nil
}

const oauth1header = `OAuth oauth_consumer_key="%s",oauth_nonce="%s",oauth_signature="%s",oauth_signature_method="%s",oauth_timestamp="%s",oauth_token="%s",oauth_version="%s"`

// setOAuth1Header returns http.Request with the header information required for OAuth1.0a authentication.
func setOAuth1Header(r *http.Request, paramsMap map[string]string, c IClient) (*http.Request, error) {
	in := &CreateOAuthSignatureInput{
		HTTPMethod:       r.Method,
		RawEndpoint:      r.URL.String(),
		OAuthConsumerKey: c.OAuthConsumerKey(),
		OAuthToken:       c.OAuthToken(),
		SigningKey:       c.SigningKey(),
		ParameterMap:     paramsMap,
	}

	out, err := CreateOAuthSignature(in)
	if err != nil {
		return nil, err
	}

	r.Header.Add("Authorization", fmt.Sprintf(oauth1header,
		url.QueryEscape(c.OAuthConsumerKey()),
		url.QueryEscape(out.OAuthNonce),
		url.QueryEscape(out.OAuthSignature),
		url.QueryEscape(out.OAuthSignatureMethod),
		url.QueryEscape(out.OAuthTimestamp),
		url.QueryEscape(c.OAuthToken()),
		url.QueryEscape(out.OAuthVersion),
	))

	return r, nil
}

func newRequest(ctx context.Context, endpoint, method string, p util.Parameters) (*http.Request, error) {
	body, err := p.Body()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	return req, nil
}

func resolveNon2XXResponse(res *http.Response) (*resources.Non2XXError, error) {
	non200err := &resources.Non2XXError{
		Status:     res.Status,
		StatusCode: res.StatusCode,
	}

	cts := util.HeaderValues("Content-Type", res.Header)
	if len(cts) == 0 {
		non200err.APIErrors = []resources.ErrorInformation{
			{Message: "Content-Type is undefined."},
		}
		return non200err, nil
	}

	if !strings.Contains(cts[0], "application/json") {
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		non200err.APIErrors = []resources.ErrorInformation{
			{Message: strings.TrimRight(string(bytes), "\n")},
		}
	} else {
		if err := json.NewDecoder(res.Body).Decode(non200err); err != nil {
			return nil, err
		}
	}

	// additional information for Rate Limit
	if res.StatusCode == http.StatusTooManyRequests {
		rri, err := util.GetRateLimitInformation(res)
		if err != nil {
			return nil, err
		}

		non200err.RateLimitInfo = rri
	}

	return non200err, nil
}
