package gotwierrors

const (
	ErrorClientNotReady string = "Twitter client is not ready."
	ErrorParametersNil  string = "Parameter for %s is nil."
	ErrorNon2XXStatus   string = "Twitter API returned a status other than 200. Status: %s."
	ErrorUndefined      string = "Undefined error."
)
