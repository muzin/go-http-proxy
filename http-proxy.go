package go_http_proxy

import (
	"crypto/tls"
	"github.com/gorilla/websocket"
	"github.com/muzin/go_rt/events"
	"github.com/muzin/go_rt/try"
	"net"
	"net/http"
	"net/http/httputil"
	urlpkg "net/url"
	"strconv"
	"strings"
)

type ProxyServer struct {
	events.EventEmitter

	web func(response http.ResponseWriter, request *http.Request)

	ws func(request *http.Request, conn *websocket.Conn)

	proxyReqFunc func(proxyReq *http.Request, request *http.Request, response http.ResponseWriter, options *ProxyServerOptions)

	proxyRespFunc func(proxyResp *http.Response, request *http.Request, response http.ResponseWriter, options *ProxyServerOptions)

	proxyErrorFunc func(resp http.ResponseWriter, req *http.Request, err error)

	options *ProxyServerOptions
}

func (this *ProxyServer) onError(args ...interface{}) {
	throwable := args[0].(try.Throwable)
	try.Throw(throwable)
}

//
// web proxy
//
// @param response http.ResponseWriter
//
// @param request *http.Request
//
// @param extOptions *ProxyServerOptions
func (this *ProxyServer) Web(response http.ResponseWriter, request *http.Request, args ...interface{}) {

	var extOptions *ProxyServerOptions = nil
	if len(args) > 0 {
		extOptions = args[0].(*ProxyServerOptions)
	} else {
		extOptions = this.options
	}

	target := extOptions.Target
	u, _ := urlpkg.Parse(target)

	requestUrl := request.URL

	targetUrl := &urlpkg.URL{
		Scheme: u.Scheme,
		Opaque: requestUrl.Opaque,
		User:   requestUrl.User,
		Host:   u.Host,
		//Path:        u.Path + requestUrl.Path,
		Path:        u.Path,
		RawPath:     requestUrl.RawPath,
		ForceQuery:  requestUrl.ForceQuery,
		RawQuery:    requestUrl.RawQuery,
		Fragment:    requestUrl.Fragment,
		RawFragment: requestUrl.RawFragment,
	}
	request.Host = u.Host

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	var tlsConfig = &tls.Config{
		InsecureSkipVerify: !extOptions.Secure, // 忽略证书验证
	}
	var transport http.RoundTripper = &http.Transport{
		Proxy: nil, // 不使用代理，如果想使用系统代理，请使用 http.ProxyFromEnvironment
		DialContext: (&net.Dialer{
			Timeout:   extOptions.ProxyTimeout,
			KeepAlive: extOptions.KeepAlive,
		}).DialContext,
		MaxIdleConns:          extOptions.MaxIdleConns,
		IdleConnTimeout:       extOptions.IdleConnTimeout,
		TLSHandshakeTimeout:   extOptions.TLSHandshakeTimeout,
		ExpectContinueTimeout: extOptions.ExpectContinueTimeout,
		TLSClientConfig:       tlsConfig,
		DisableCompression:    extOptions.DisableCompression,
	}
	proxy.Transport = transport
	proxy.ModifyResponse = func(r *http.Response) error {
		r.Header.Del("X-Frame-Options") // 重点：代理时移除 X-Frame-Options 头
		if !extOptions.Xfwd {
			r.Header.Del("X-Forwarded-For") // Xfwd为false时, 移除 X-Forwarded-For 头
		}

		// 如果 存在 附加 Header
		if extOptions.Headers != nil && len(extOptions.Headers) > 0 {
			for key, value := range extOptions.Headers {
				r.Header.Add(key, value)
			}
		}

		if this.proxyRespFunc != nil {
			this.proxyRespFunc(r, request, response, extOptions)
		}
		return nil
	}
	// Proxy Error Handler
	if this.proxyErrorFunc != nil {
		proxy.ErrorHandler = this.proxyErrorFunc
	}

	proxy.ServeHTTP(response, request)

}

func (this *ProxyServer) Listen(port int, hostname string) error {

	options := this.options

	webClosure := func(response http.ResponseWriter, request *http.Request) {
		this.Web(response, request)
	}

	//fmt.Println(closure)

	addr := hostname + ":" + strconv.Itoa(port)

	handler := &DefaultHttpProxyHandler{}
	handler.SetWebHandler(webClosure)

	server := &http.Server{Addr: addr, Handler: handler}

	var err error = nil
	if options.Ssl != nil {
		err = http.ListenAndServeTLS(addr, "nil", "nil", handler)
	} else {
		return server.ListenAndServe()
	}

	return err
}

//
// OnProxyReq
//
// @param listener func(listener func(proxyReq *http.Request, request *http.Request, response http.ResponseWriter, options *ProxyServerOptions)
//
func (this *ProxyServer) OnProxyReq(listener func(proxyReq *http.Request,
	request *http.Request,
	response http.ResponseWriter,
	options *ProxyServerOptions,
)) {
	this.proxyReqFunc = listener
}

//
// OnProxyReqWs
//
// @param listener func(listener func(proxyReq *http.Request, request *http.Request, response http.ResponseWriter, options *ProxyServerOptions)
//
func (this *ProxyServer) OnProxyReqWs(listener func(proxyReq *http.Request,
	request *http.Request,
	response http.ResponseWriter,
	options *ProxyServerOptions,
)) {
	this.proxyReqFunc = listener
}

//
// OnProxyResp
//
// @param listener func(listener func(proxyReq *http.Request, request *http.Request, response http.ResponseWriter, options *ProxyServerOptions)
//
func (this *ProxyServer) OnProxyResp(listener func(proxyResp *http.Response,
	request *http.Request,
	response http.ResponseWriter,
	options *ProxyServerOptions,
)) {
	this.proxyRespFunc = listener
}

func (this *ProxyServer) OnProxyError(listener func(resp http.ResponseWriter,
	req *http.Request,
	err error)) {
	this.proxyErrorFunc = listener
}

//
// OnError
//
// @param listener func(throwable try.Throwable)
//
func (this *ProxyServer) OnError(listener func(...interface{})) {
	this.AddListener("error", listener)
}

//
// OnClose
//
// @param listener func(throwable try.Throwable)
//
func (this *ProxyServer) OnClose(listener func(...interface{})) {
	this.AddListener("close", listener)
}

func (this *ProxyServer) createProxyRequestByRequest(request *http.Request) *http.Request {

	method := request.Method
	url := request.URL
	body := request.Body

	newRequest, _ := http.NewRequest(method, url.RequestURI(), body)
	newRequest.Form = request.Form
	newRequest.ContentLength = request.ContentLength
	newRequest.Header = request.Header
	newRequest.MultipartForm = request.MultipartForm
	newRequest.PostForm = request.PostForm
	//newRequest.RemoteAddr = request.RemoteAddr
	newRequest.Proto = request.Proto
	newRequest.ProtoMajor = request.ProtoMajor
	newRequest.ProtoMinor = request.ProtoMinor
	newRequest.TLS = request.TLS
	newRequest.Trailer = request.Trailer
	newRequest.TransferEncoding = request.TransferEncoding

	return newRequest
}

type DefaultHttpProxyHandler struct {
	webHandler func(response http.ResponseWriter, request *http.Request)

	upgrader websocket.Upgrader
}

func (this *DefaultHttpProxyHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if this.webHandler != nil {
		this.webHandler(response, request)
	}
}

func (this *DefaultHttpProxyHandler) SetWebHandler(f func(resp http.ResponseWriter, request *http.Request)) {
	this.webHandler = f
}

func NewProxyServer(options *ProxyServerOptions) *ProxyServer {

	var proxyServer = &ProxyServer{
		options: options,
	}

	proxyServer.On("error", proxyServer.onError)

	return proxyServer
}

// Given a string of the form "host", "host:port", or "[ipv6::address]:port",
// return true if the string includes a port.
func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

// removeEmptyPort strips the empty port in ":port" to ""
// as mandated by RFC 3986 Section 6.2.3.
func removeEmptyPort(host string) string {
	if hasPort(host) {
		return strings.TrimSuffix(host, ":")
	}
	return host
}
