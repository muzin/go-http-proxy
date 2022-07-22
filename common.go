package go_http_proxy

import "regexp"

var (
	upgradeHeader      = "(^|,)\\s*upgrade\\s*($|,)"
	upgradeHeaderRegex = regexp.MustCompile(upgradeHeader)

	isSSL      = "^https|wss"
	isSSLRegex = regexp.MustCompile(isSSL)
)

//
// Copies the right headers from `options` and `req` to
// `outgoing` which is then used to fire the proxied
// request.
//
// Examples:
//
//    common.setupOutgoing(outgoing, options, req)
//    // => { host: ..., hostname: ...}
//
// @param {Object} Outgoing Base object to be filled with required properties
//
// @param {Object} Options Config object passed to the proxy
//
// @param {ClientRequest} Req Request Object
//
// @param {String} Forward String to select forward or target
//
// @return {Object} Outgoing Object with all required properties set
//
// @api private
//
func SetupOutgoing() {

}
