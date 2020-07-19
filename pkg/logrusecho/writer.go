package logrusecho

import "net/http"

// LogrusWriter is a wrapper around the http.ResponseWriter.
// It is designed to catch all incoming writings to pick useful informations and store it durably into the logrus response object.
type LogrusWriter struct {
	writer   http.ResponseWriter
	Response *LogrusResponse
}

// NewLogrusWriter : ctor
func NewLogrusWriter(writer http.ResponseWriter) *LogrusWriter {
	return &LogrusWriter{
		writer:   writer,
		Response: &LogrusResponse{Status: http.StatusOK, Size: 0},
	}
}

// Write is a wrapper for the "real" ResponseWriter.Write
func (lw *LogrusWriter) Write(b []byte) (int, error) {
	size, err := lw.writer.Write(b)
	lw.Response.Size += size
	return size, err
}

// WriteHeader is a wrapper around ResponseWriter.WriteHeader
func (lw *LogrusWriter) WriteHeader(s int) {
	lw.writer.WriteHeader(s)
	lw.Response.Status = s
}

// Header is a wrapper around ResponseWriter.Header
func (lw *LogrusWriter) Header() http.Header {
	return lw.writer.Header()
}
