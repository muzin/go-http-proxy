package go_http_proxy

import (
	"time"
)

// ProxyServerOptions
type ProxyServerOptions struct {

	// url string to be parsed with the url module
	Target string

	// object to be passed to http(s).request
	//Agent *http.Request

	// object to be passed to https.createServer()
	Ssl *ProxyServerOptions

	// true/false, if you want to proxy websockets
	//Ws bool

	// true/false, Default: true adds x-forward headers
	Xfwd bool

	// true/false, Default: false  verify SSL certificate
	Secure bool

	// true/false, explicitly specify if we are proxying to another proxy
	//ToProxy bool

	// true/false, Default: true - specify whether you want to prepend the target's path to the proxy path
	//PrependPath bool

	// true/false, Default: false - specify whether you want to ignore the proxy path of the incoming request
	//IgnorePath bool

	// Local interface string to bind for outgoing connections
	//LocalAddress string

	// true/false, Default: false - changes the origin of the host header to the target URL
	//ChangeOrigin bool

	// true/false, Default: false - specify whether you want to keep letter case of response header key
	//PreserveHeaderKeyCase bool

	// Basic authentication i.e. 'user:password' to compute an Authorization header
	Auth string

	// rewrites the location hostname on (201/301/302/307/308) redirects, Default: ""
	//HostRewrite string

	// rewrites the location host/port on (201/301/302/307/308) redirects based on requested host/port. Default: false
	//AutoRewrite bool

	// rewrites the location protocol on (201/301/302/307/308) redirects to 'http' or 'https'. Default: ""
	//ProtocolRewrite string

	// rewrites domain of `set-cookie` headers
	CookieDomainRewrite map[string]string

	// rewrites path of `set-cookie` headers
	CookiePathRewrite map[string]string

	// object with extra headers to be added to target requests
	Headers map[string]string

	// ProxyTimeout Default: 60s
	ProxyTimeout time.Duration

	// KeepAlive Default: 60s
	KeepAlive time.Duration

	// Maximum idle connections default: 100
	MaxIdleConns int

	// Idle Conn Timeout default: 90s
	IdleConnTimeout time.Duration

	// Idle Conn Timeout default: 10s
	TLSHandshakeTimeout time.Duration

	// Idle Conn Timeout default: 3s
	ExpectContinueTimeout time.Duration

	// Disable Compression default: true
	DisableCompression bool

	// true/false, Default: true - specify whether you want to follow redirects
	FollowRedirects bool
}

func NewProxyServerOptions() *ProxyServerOptions {
	proxyServerOptions := &ProxyServerOptions{
		Xfwd:                  true,
		Secure:                false,
		ProxyTimeout:          60 * time.Second,
		KeepAlive:             60 * time.Second,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 3 * time.Second,
		DisableCompression:    true,
		CookieDomainRewrite:   make(map[string]string),
		CookiePathRewrite:     make(map[string]string),
		Headers:               make(map[string]string),
		FollowRedirects:       true,
	}
	return proxyServerOptions
}
