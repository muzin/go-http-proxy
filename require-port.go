package go_http_proxy

import "strings"

//
// Check if we're required to add a port number.
//
// @param int port Port number we need to check
//
// @param string protocol Protocol we need to check against.
//
// @returns bool Is it a default port for the given protocol
//
// @api private
//
func required(port int, protocol string) bool {
	protocol = strings.Split(protocol, ":")[0]
	port = +port

	if port == 0 || port > 65535 {
		return false
	}

	switch protocol {
	case "http", "ws":
		return port != 80
	case "https", "wss":
		return port != 443
	case "ftp":
		return port != 21
	case "gopher":
		return port != 70
	case "file":
		return false
	}
	return port != 0
}
