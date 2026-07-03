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

	// HTTPClient, when non-nil, is used as-is and takes precedence over the
	// transport-shaping options above (Timeout, SSLVerify, FollowRedirects),
	// which the caller is then responsible for configuring on their client.
	// This is the escape hatch for custom transports, retry round-trippers,
	// instrumentation, etc.
	HTTPClient *http.Client

	// UserAgent is sent as the User-Agent header. Empty means the library
	// default is used.
	UserAgent string
}

// DefaultClientConfig returns a ClientConfig with sensible defaults.
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Timeout:         DefaultConfigTimeout,
		SSLVerify:       true,
		FollowRedirects: true,
		UserAgent:       defaultUserAgent,
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

// WithHTTPClient supplies a custom *http.Client to use as-is. When set, it
// takes precedence over WithConfigTimeout/WithSSLVerify/WithFollowRedirects.
func WithHTTPClient(client *http.Client) ConfigOption {
	return func(c *ClientConfig) {
		c.HTTPClient = client
	}
}

// WithUserAgent sets the User-Agent header sent with requests.
func WithUserAgent(userAgent string) ConfigOption {
	return func(c *ClientConfig) {
		c.UserAgent = userAgent
	}
}

// ToHTTPClient converts the ClientConfig to a configured http.Client.
// If a custom HTTPClient was supplied, it is returned unchanged.
func (c *ClientConfig) ToHTTPClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}

	// Clone the standard default transport so we keep its proxy support,
	// HTTP/2, and connection-pool/timeout defaults, then override only TLS
	// verification. A bare &http.Transport{} would drop all of those.
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{}
	}
	transport.TLSClientConfig.InsecureSkipVerify = !c.SSLVerify

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

// Clone creates a copy of the ClientConfig. The HTTPClient pointer is shared,
// not deep-copied, since an *http.Client is intended to be reused.
func (c *ClientConfig) Clone() *ClientConfig {
	return &ClientConfig{
		Timeout:         c.Timeout,
		SSLVerify:       c.SSLVerify,
		FollowRedirects: c.FollowRedirects,
		HTTPClient:      c.HTTPClient,
		UserAgent:       c.UserAgent,
	}
}
