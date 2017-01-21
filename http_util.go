package cqrs

import "net/http"

func absoluteUrl(r *http.Request) string {
	return r.URL.Scheme + "://" + r.Host + r.URL.Path
}
