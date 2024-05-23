package util

import (
	"fmt"
	"net/http"
)

// Mimetypes tells computers how the media should be interpreted.
// A practical example is formatting HTTP body in the JSON format.
//
// Example:
// w := http.ResponseWriter
// w.Header().Add(util.CONTENT_TYPE, util.MIMETYPE_JSON)
const (
	MIMETYPE_JSON           = "application/json"
	MIMETYPE_PLAINTEXT      = "text/plain"
	MIMETYPE_PLAINTEXT_UTF8 = "text/plain; charset=utf-8"
)

// List of commonly used headers.
//
// Example:
// w := http.ResponseWriter
// w.Header().Add(util.CONTENT_TYPE, util.MIMETYPE_JSON)
const (
	CONTENT_TYPE          = "content-type"
	X_CONTENT_TYPE_OPTION = "X-Content-Type-Options"
)

// HttpError is a drop-in replacement for http.Error.
//
// # Description
//
// The problem with http.Error is it is generic and not specific to all use-cases. HttpError builds on top of http.Error
// by replacing the mimetype with 'application/json' and formatting the returning message to as a JSON object.
//
// # Example
//
// util.HttpError(w, "this is a test", http.StatusBadRequest)
//
// Output:
//
//	{
//	    "message": "this is a test"
//	}
func HttpError(w http.ResponseWriter, error string, code int) {
	w.Header().Add(CONTENT_TYPE, MIMETYPE_JSON)
	w.Header().Set(X_CONTENT_TYPE_OPTION, "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, "{\"message\": \""+error+"\"}")
}
