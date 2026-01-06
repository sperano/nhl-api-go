package nhl

import (
	"crypto/tls"
	"net/http"
	"time"
)

const (
	// DefaultConfigTimeout is the default HTTP client timeout.
	DefaultConfigTimeout = 10 * time.Second
)

// ClientConfig holds configuration options for the NHL API client.
type ClientConfig struct {
	// Timeout is the maximum duration for HTTP requests.
	Timeout time.Duration

	// SSLVerify controls whether SSL certificates are verified.
	SSLVerify bool

	// FollowRedirects controls whether HTTP redirects are followed.
	FollowRedirects bool
}

// DefaultClientConfig returns a ClientConfig with sensible defaults.
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Timeout:         DefaultConfigTimeout,
		SSLVerify:       true,
		FollowRedirects: true,
	}
}

// ConfigOption is a functional option for configuring ClientConfig.
type ConfigOption func(*ClientConfig)

// NewClientConfig creates a new ClientConfig with the provided options.
// If no options are provided, defaults are used.
func NewClientConfig(opts ...ConfigOption) *ClientConfig {
	cfg := DefaultClientConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

// WithConfigTimeout sets the HTTP client timeout.
func WithConfigTimeout(timeout time.Duration) ConfigOption {
	return func(c *ClientConfig) {
		c.Timeout = timeout
	}
}

// WithSSLVerify sets whether SSL certificates should be verified.
func WithSSLVerify(verify bool) ConfigOption {
	return func(c *ClientConfig) {
		c.SSLVerify = verify
	}
}

// WithFollowRedirects sets whether HTTP redirects should be followed.
func WithFollowRedirects(follow bool) ConfigOption {
	return func(c *ClientConfig) {
		c.FollowRedirects = follow
	}
}

// ToHTTPClient converts the ClientConfig to a configured http.Client.
func (c *ClientConfig) ToHTTPClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !c.SSLVerify,
		},
	}

	client := &http.Client{
		Timeout:   c.Timeout,
		Transport: transport,
	}

	if !c.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return client
}

// Clone creates a deep copy of the ClientConfig.
func (c *ClientConfig) Clone() *ClientConfig {
	return &ClientConfig{
		Timeout:         c.Timeout,
		SSLVerify:       c.SSLVerify,
		FollowRedirects: c.FollowRedirects,
	}
}
