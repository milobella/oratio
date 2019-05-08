package logrushttp

// LogrusResponse is used to store durably some useful information from the HTTP response.
type LogrusResponse struct {
	Status int
	Size   int
}
