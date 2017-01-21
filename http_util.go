package cqrs

import "net/http"

func absoluteUrl(r *http.Request) string {
	return r.URL.Scheme + "://" + r.URL.Host + r.URL.Path
}
