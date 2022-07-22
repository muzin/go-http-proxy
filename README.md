# go-http-proxy
go-http-proxy

## example
```go

import (
	"fmt"
	go_http_proxy "github.com/muzin/go-http-proxy"
	"net/http"
	"strings"
)

type MyHttpHandler2 struct {
	proxyServer *go_http_proxy.ProxyServer
}
func (this *MyHttpHandler2) ServeHTTP(resp http.ResponseWriter, req *http.Request){
	this.proxyServer.Web(resp, req)
}

func main() {

	options := go_http_proxy.NewProxyServerOptions()
	options.Target = "http://localhost:19001"

	options.Headers["X-Clover-Proxy"] = "true"

	proxyServer := go_http_proxy.NewProxyServer(options)

	proxyServer.OnProxyError(func(resp http.ResponseWriter, req *http.Request, err error) {
		// 如果 连接被拒绝 ， 跳转到
		if strings.HasSuffix(err.Error(), "connect: connection refused") {
			resp.Header().Add("Location", "http://baidu.com")
			resp.WriteHeader(http.StatusFound)
		}
	})

	proxyServer.Listen(19000, "0.0.0.0")

	fmt.Println("1")

}

```