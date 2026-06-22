// (C) Copyright 2026 Hewlett Packard Enterprise Development LP
package auth

import "strings"

// NormalizeBaseURL trims trailing slashes from the Morpheus appliance base URL.
//
// The generated SDK builds request paths by concatenating the configured base
// URL with an absolute path, e.g. baseURL + "/oauth/token". If the configured
// URL carries a trailing slash (e.g. "https://morpheus.example.com/") the
// result is a double slash ("https://morpheus.example.com//oauth/token").
//
// Morpheus 9.0 serves its OAuth2 token endpoint from a Spring Authorization
// Server mounted at the servlet path "/oauth/*". A double-slash request such as
// "//oauth/token" does not match that mapping, so the request bypasses the
// authorization server and Morpheus returns a non-JSON error page. The SDK then
// fails to decode the response ("undefined response type"), which surfaces as an
// authentication failure. The same applies to "//api/..." requests once the
// token is obtained. Trimming trailing slashes keeps every request path
// single-slashed and routable.
func NormalizeBaseURL(url string) string {
	return strings.TrimRight(url, "/")
}
