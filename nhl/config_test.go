package nhl

import (
	"net/http"
	"testing"
	"time"
)

func TestDefaultClientConfig(t *testing.T) {
	cfg := DefaultClientConfig()

	if cfg.Timeout != DefaultConfigTimeout {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, DefaultConfigTimeout)
	}

	if !cfg.SSLVerify {
		t.Error("SSLVerify should be true by default")
	}

	if !cfg.FollowRedirects {
		t.Error("FollowRedirects should be true by default")
	}
}

func TestNewClientConfig(t *testing.T) {
	t.Run("with no options", func(t *testing.T) {
		cfg := NewClientConfig()

		if cfg.Timeout != DefaultConfigTimeout {
			t.Errorf("Timeout = %v, want %v", cfg.Timeout, DefaultConfigTimeout)
		}

		if !cfg.SSLVerify {
			t.Error("SSLVerify should be true by default")
		}

		if !cfg.FollowRedirects {
			t.Error("FollowRedirects should be true by default")
		}
	})

	t.Run("with timeout option", func(t *testing.T) {
		customTimeout := 30 * time.Second
		cfg := NewClientConfig(WithConfigTimeout(customTimeout))

		if cfg.Timeout != customTimeout {
			t.Errorf("Timeout = %v, want %v", cfg.Timeout, customTimeout)
		}

		// Other fields should have defaults
		if !cfg.SSLVerify {
			t.Error("SSLVerify should be true by default")
		}

		if !cfg.FollowRedirects {
			t.Error("FollowRedirects should be true by default")
		}
	})

	t.Run("with SSL verify disabled", func(t *testing.T) {
		cfg := NewClientConfig(WithSSLVerify(false))

		if cfg.SSLVerify {
			t.Error("SSLVerify should be false")
		}

		// Other fields should have defaults
		if cfg.Timeout != DefaultConfigTimeout {
			t.Errorf("Timeout = %v, want %v", cfg.Timeout, DefaultConfigTimeout)
		}

		if !cfg.FollowRedirects {
			t.Error("FollowRedirects should be true by default")
		}
	})

	t.Run("with redirects disabled", func(t *testing.T) {
		cfg := NewClientConfig(WithFollowRedirects(false))

		if cfg.FollowRedirects {
			t.Error("FollowRedirects should be false")
		}

		// Other fields should have defaults
		if cfg.Timeout != DefaultConfigTimeout {
			t.Errorf("Timeout = %v, want %v", cfg.Timeout, DefaultConfigTimeout)
		}

		if !cfg.SSLVerify {
			t.Error("SSLVerify should be true by default")
		}
	})

	t.Run("with multiple options", func(t *testing.T) {
		customTimeout := 5 * time.Second
		cfg := NewClientConfig(
			WithConfigTimeout(customTimeout),
			WithSSLVerify(false),
			WithFollowRedirects(false),
		)

		if cfg.Timeout != customTimeout {
			t.Errorf("Timeout = %v, want %v", cfg.Timeout, customTimeout)
		}

		if cfg.SSLVerify {
			t.Error("SSLVerify should be false")
		}

		if cfg.FollowRedirects {
			t.Error("FollowRedirects should be false")
		}
	})
}

func TestClientConfig_ToHTTPClient(t *testing.T) {
	t.Run("default configuration", func(t *testing.T) {
		cfg := DefaultClientConfig()
		client := cfg.ToHTTPClient()

		if client.Timeout != cfg.Timeout {
			t.Errorf("client.Timeout = %v, want %v", client.Timeout, cfg.Timeout)
		}

		if client.Transport == nil {
			t.Fatal("client.Transport should not be nil")
		}

		transport, ok := client.Transport.(*http.Transport)
		if !ok {
			t.Fatal("client.Transport should be *http.Transport")
		}

		if transport.TLSClientConfig == nil {
			t.Fatal("TLSClientConfig should not be nil")
		}

		if transport.TLSClientConfig.InsecureSkipVerify {
			t.Error("InsecureSkipVerify should be false when SSLVerify is true")
		}

		// Default should follow redirects (CheckRedirect should be nil)
		if client.CheckRedirect != nil {
			t.Error("CheckRedirect should be nil when FollowRedirects is true")
		}
	})

	t.Run("SSL verification disabled", func(t *testing.T) {
		cfg := NewClientConfig(WithSSLVerify(false))
		client := cfg.ToHTTPClient()

		transport, ok := client.Transport.(*http.Transport)
		if !ok {
			t.Fatal("client.Transport should be *http.Transport")
		}

		if !transport.TLSClientConfig.InsecureSkipVerify {
			t.Error("InsecureSkipVerify should be true when SSLVerify is false")
		}
	})

	t.Run("redirects disabled", func(t *testing.T) {
		cfg := NewClientConfig(WithFollowRedirects(false))
		client := cfg.ToHTTPClient()

		if client.CheckRedirect == nil {
			t.Error("CheckRedirect should not be nil when FollowRedirects is false")
		}

		// Test that CheckRedirect returns the expected error
		err := client.CheckRedirect(nil, nil)
		if err != http.ErrUseLastResponse {
			t.Errorf("CheckRedirect should return http.ErrUseLastResponse, got %v", err)
		}
	})

	t.Run("custom timeout", func(t *testing.T) {
		customTimeout := 15 * time.Second
		cfg := NewClientConfig(WithConfigTimeout(customTimeout))
		client := cfg.ToHTTPClient()

		if client.Timeout != customTimeout {
			t.Errorf("client.Timeout = %v, want %v", client.Timeout, customTimeout)
		}
	})
}

func TestClientConfig_Clone(t *testing.T) {
	original := NewClientConfig(
		WithConfigTimeout(15*time.Second),
		WithSSLVerify(false),
		WithFollowRedirects(false),
	)

	cloned := original.Clone()

	// Verify all fields match
	if cloned.Timeout != original.Timeout {
		t.Errorf("cloned.Timeout = %v, want %v", cloned.Timeout, original.Timeout)
	}

	if cloned.SSLVerify != original.SSLVerify {
		t.Errorf("cloned.SSLVerify = %v, want %v", cloned.SSLVerify, original.SSLVerify)
	}

	if cloned.FollowRedirects != original.FollowRedirects {
		t.Errorf("cloned.FollowRedirects = %v, want %v", cloned.FollowRedirects, original.FollowRedirects)
	}

	// Verify it's a different instance
	if cloned == original {
		t.Error("cloned config should be a different instance than original")
	}

	// Verify modifying clone doesn't affect original
	cloned.Timeout = 20 * time.Second
	if original.Timeout == cloned.Timeout {
		t.Error("modifying clone should not affect original")
	}
}

func TestConfigOptions(t *testing.T) {
	t.Run("WithConfigTimeout", func(t *testing.T) {
		cfg := &ClientConfig{}
		timeout := 25 * time.Second

		opt := WithConfigTimeout(timeout)
		opt(cfg)

		if cfg.Timeout != timeout {
			t.Errorf("Timeout = %v, want %v", cfg.Timeout, timeout)
		}
	})

	t.Run("WithSSLVerify true", func(t *testing.T) {
		cfg := &ClientConfig{}

		opt := WithSSLVerify(true)
		opt(cfg)

		if !cfg.SSLVerify {
			t.Error("SSLVerify should be true")
		}
	})

	t.Run("WithSSLVerify false", func(t *testing.T) {
		cfg := &ClientConfig{}

		opt := WithSSLVerify(false)
		opt(cfg)

		if cfg.SSLVerify {
			t.Error("SSLVerify should be false")
		}
	})

	t.Run("WithFollowRedirects true", func(t *testing.T) {
		cfg := &ClientConfig{}

		opt := WithFollowRedirects(true)
		opt(cfg)

		if !cfg.FollowRedirects {
			t.Error("FollowRedirects should be true")
		}
	})

	t.Run("WithFollowRedirects false", func(t *testing.T) {
		cfg := &ClientConfig{}

		opt := WithFollowRedirects(false)
		opt(cfg)

		if cfg.FollowRedirects {
			t.Error("FollowRedirects should be false")
		}
	})
}
