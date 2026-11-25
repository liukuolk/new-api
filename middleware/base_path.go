/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/QuantumNous/new-api/common"
)

// BasePathHandler strips common.BasePath from the request URL before the Gin
// engine resolves a route.
//
// This has to wrap the engine instead of being registered with server.Use:
// Gin selects the handler chain from Request.URL.Path at the start of
// ServeHTTP, before any middleware runs, so a middleware rewrite would only
// affect handlers that were already chosen (NoRoute) and every registered
// route under the prefix would 404.
//
// Requests that do not carry the prefix are forwarded unchanged so a gateway
// that already removes it keeps working. When common.BasePath is empty the
// engine is returned as-is.
func BasePathHandler(next http.Handler) http.Handler {
	base := common.BasePath
	if base == "" {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path != base && !strings.HasPrefix(path, base+"/") {
			next.ServeHTTP(w, r)
			return
		}

		stripped := strings.TrimPrefix(path, base)
		if stripped == "" {
			stripped = "/"
		}

		cloned := new(http.Request)
		*cloned = *r
		cloned.URL = new(url.URL)
		*cloned.URL = *r.URL
		cloned.URL.Path = stripped
		if cloned.URL.RawPath != "" {
			cloned.URL.RawPath = strings.TrimPrefix(cloned.URL.RawPath, base)
			if cloned.URL.RawPath == "" {
				cloned.URL.RawPath = "/"
			}
		}
		cloned.RequestURI = cloned.URL.RequestURI()

		next.ServeHTTP(w, cloned)
	})
}
