package go_http_proxy

type ProxyType string

const (
	WS ProxyType = "ws"

	WEB ProxyType = "web"
)

func (this ProxyType) String() string {
	switch this {
	case WS:
		return "ws"
	case WEB:
		return "web"
	default:
		return ""
	}
}
