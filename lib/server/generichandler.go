package server

import (
	"github.com/dustin/go-humanize"
	"github.com/rs/cors"
	"net/http"
)

// HandlerFunc - useful to chain different middleware http.Handler
type HandlerFunc func(http.Handler) http.Handler

func RegisterHandlers(h http.Handler, handlerFns ...HandlerFunc) http.Handler {
	for _, hFn := range handlerFns {
		h = hFn(h)
	}
	return h
}

// Adds limiting body size middleware

// Maximum allowed form data field values. 64MiB is a guessed practical value
// which is more than enough to accommodate any form data fields and headers.
const requestFormDataSize = 64 * humanize.MiByte

// For any HTTP request, request body should be not more than 16GiB + requestFormDataSize
// where, 16GiB is the maximum allowed object size for object upload.
const requestMaxBodySize = globalMaxFileSize + requestFormDataSize

type requestSizeLimitHandler struct {
	handler     http.Handler
	maxBodySize int64
}

func setRequestSizeLimitHandler(h http.Handler) http.Handler {
	return requestSizeLimitHandler{handler: h, maxBodySize: requestMaxBodySize}
}

func (h requestSizeLimitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Restricting read data to a given maximum length
	r.Body = http.MaxBytesReader(w, r.Body, h.maxBodySize)
	h.handler.ServeHTTP(w, r)
}

const (
	// Maximum size for http headers - See: https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingMetadata.html
	maxHeaderSize = 8 * 1024
	// Maximum size for user-defined metadata - See: https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingMetadata.html
	maxUserDataSize = 2 * 1024
)

type requestHeaderSizeLimitHandler struct {
	http.Handler
}

func setRequestHeaderSizeLimitHandler(h http.Handler) http.Handler {
	return requestHeaderSizeLimitHandler{h}
}

type resourceHandler struct {
	handler http.Handler
}

// List of default allowable HTTP methods.
var defaultAllowableHTTPMethods = []string{
	http.MethodGet,
	http.MethodPut,
	http.MethodHead,
	http.MethodPost,
	http.MethodDelete,
	http.MethodOptions,
}

// setCorsHandler handler for CORS (Cross Origin Resource Sharing)
func setCorsHandler(h http.Handler) http.Handler {
	commonS3Headers := []string{
		"Date",
		"ETag",
		"Server",
		"Connection",
		"Accept-Ranges",
		"Content-Range",
		"Content-Encoding",
		"Content-Length",
		"Content-Type",
		"x-amz-request-id",
	}
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   defaultAllowableHTTPMethods,
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   commonS3Headers,
		AllowCredentials: true,
	})
	return c.Handler(h)
}

// List of not implemented bucket queries
var notimplementedBucketResourceNames = map[string]bool{
	"acl":            true,
	"cors":           true,
	"lifecycle":      true,
	"logging":        true,
	"replication":    true,
	"tagging":        true,
	"versions":       true,
	"requestPayment": true,
	"versioning":     true,
	"website":        true,
	"inventory":      true,
	"metrics":        true,
	"accelerate":     true,
}

type securityHeaderHandler struct {
	handler http.Handler
}

func addSecurityHeaders(h http.Handler) http.Handler {
	return securityHeaderHandler{handler: h}
}

func (s securityHeaderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("X-XSS-Protection", "1; mode=block")                  // Prevents against XSS attacks
	header.Set("Content-Security-Policy", "block-all-mixed-content") // prevent mixed (HTTP / HTTPS content)
	s.handler.ServeHTTP(w, r)
}

// List of some generic handlers which are applied for all incoming requests.
var GlobalHandlers = []HandlerFunc{
	// set HTTP security headers such as Content-Security-Policy.
	addSecurityHeaders,

	// Limits all requests size to a maximum fixed limit
	setRequestSizeLimitHandler,
	// Limits all header sizes to a maximum fixed limit
	setRequestHeaderSizeLimitHandler,

	setCorsHandler,
}
