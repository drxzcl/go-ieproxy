// Package ieproxy is a utility to retrieve the proxy parameters (especially of Internet Explorer on windows)
//
// On windows, it gathers the parameters from the registry (regedit), while it uses env variable on other platforms
package ieproxy

import (
	"net/http"
	"net/url"
	"os"
	"strings"
)

// ProxyConf gathers the configuration for proxy
type ProxyConf struct {
	Static    StaticProxyConf    // static configuration
	Automatic AutomaticProxyConf // automatic configuration
}

// ProxyFromEnvironment is a drop-in replacement for the stdlib
// net/http.ProxyFromEnvironment function
func ProxyFromEnvironment(r *http.Request) (*url.URL, error) {
	config := GetConf()
	switch {
	case config.Static.Active:
		return config.Static.FindProxyForRequest(r)
	case config.Automatic.Active:
		return config.Automatic.FindProxyForRequest(r)
	default:
		return nil, nil
	}
}

// StaticProxyConf contains the configuration for static proxy
type StaticProxyConf struct {
	// Is the proxy active?
	Active bool
	// Proxy address for each scheme (http, https)
	// "" (empty string) is the fallback proxy
	Protocols map[string]string
	// Addresses not to be browsed via the proxy (comma-separated, linux-like)
	NoProxy string
}

// FindProxyForRequest computes the proxy for a given URL according to the static proxy rules
func (spc *StaticProxyConf) FindProxyForRequest(r *http.Request) (*url.URL, error) {
	if !spc.Active {
		return nil, nil
	}
	for _, s := range strings.Split(spc.NoProxy, ",") {
		if r.URL.String() == s {
			return nil, nil
		}
	}
	if proxyURL, ok := spc.Protocols[r.URL.Scheme]; ok {
		return url.Parse(proxyURL)
	}
	if proxyURL, ok := spc.Protocols[""]; ok {
		return url.Parse(proxyURL)
	}
	return nil, nil
}

// AutomaticProxyConf contains the configuration for automatic proxy
type AutomaticProxyConf struct {
	// Is the proxy active?
	Active bool
	// URL of the .pac file
	URL string
}

// GetConf retrieves the proxy configuration from the Windows Regedit
func GetConf() ProxyConf {
	return getConf()
}

// OverrideEnvWithStaticProxy writes new values to the
// `http_proxy`, `https_proxy` and `no_proxy` environment variables.
// The values are taken from the Windows Regedit (should be called in `init()` function - see example)
func OverrideEnvWithStaticProxy() {
	overrideEnvWithStaticProxy(GetConf(), os.Setenv)
}

// FindProxyForRequest computes the proxy for a given URL according to the pac file
func (apc *AutomaticProxyConf) FindProxyForRequest(r *http.Request) (*url.URL, error) {
	return apc.findProxyForURL(r.URL)
}

type envSetter func(string, string) error
