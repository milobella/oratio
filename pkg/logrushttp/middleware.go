package logrushttp

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// LogrusMiddleware : Logrus HTTP middleware used to log automatically a request and response with logrus.
type LogrusMiddleware interface {
	Handle(next http.Handler) http.Handler
}

// LogrusMiddlewareBuilder : Builder of Logrus HTTP middleware. All fields are optional and can be ommited.
type LogrusMiddlewareBuilder interface {
	// Activated request data from default ones (all default ones if ommited).
	ActivatedRequestData([]string) LogrusMiddlewareBuilder

	// Custom request data (will replace a default one if same name).
	CustomRequestData([]RequestData) LogrusMiddlewareBuilder

	// Custom request message (if you don't like the awesome default one).
	RequestMessage(string) LogrusMiddlewareBuilder

	// Activated response data from default ones (all default ones if ommited)
	ActivatedResponseData([]string) LogrusMiddlewareBuilder

	// Custom response data (will replace a default one if same name)
	CustomResponseData([]ResponseData) LogrusMiddlewareBuilder

	// Custom response message (if you don't like the awesome default one).
	ResponseMessage(string) LogrusMiddlewareBuilder

	Build() LogrusMiddleware
}

type logrusMiddleware struct {
	requestData    map[string]requestAccessor
	requestMessage string

	responseData    map[string]responseAccessor
	responseMessage string
}

type logrusMiddlewareBuilder struct {
	activatedRequestData []string
	customRequestData    []RequestData
	requestMessage       string

	activatedResponseData []string
	customResponseData    []ResponseData
	responseMessage       string
}

type requestAccessor func(*http.Request) interface{}
type responseAccessor func(*LogrusResponse) interface{}

// RequestData : Name of the field that we'll create in the request logrus log + mean to access it from request
type RequestData struct {
	name     string
	accessor requestAccessor
}

// ResponseData : Name of the field that we'll create in the response logrus log + mean to access it from response
type ResponseData struct {
	name     string
	accessor responseAccessor
}

// Handle : Handle function than can be passed to router or to HandleFunc chain.
func (lrm *logrusMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Use the accessor of each of our request data objects to build the logrus fields from request
		fields := make(map[string]interface{})
		for name, accessor := range lrm.requestData {
			fields[name] = accessor(r)
		}

		// Log the request
		logrus.WithFields(fields).Print(lrm.requestMessage)

		// Wrap the real writer with the logrus one
		logrusWriter := NewLogrusWriter(w)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(logrusWriter, r)

		// Retrieve the parsed response from our writer wrapper
		response := logrusWriter.Response

		// Use the accessor of each of our response data objects to build the logrus fields from response
		fields = make(map[string]interface{})
		for name, accessor := range lrm.responseData {
			fields[name] = accessor(response)
		}

		// Log the response
		logrus.WithFields(fields).Print(lrm.responseMessage)

	})
}

// NewLogrusMiddlewareBuilder : ctor
func NewLogrusMiddlewareBuilder() LogrusMiddlewareBuilder {
	return &logrusMiddlewareBuilder{}
}

func (lrmb *logrusMiddlewareBuilder) ActivatedRequestData(activatedRequestData []string) LogrusMiddlewareBuilder {
	lrmb.activatedRequestData = activatedRequestData
	return lrmb
}

func (lrmb *logrusMiddlewareBuilder) CustomRequestData(customRequestData []RequestData) LogrusMiddlewareBuilder {
	lrmb.customRequestData = customRequestData
	return lrmb
}

func (lrmb *logrusMiddlewareBuilder) RequestMessage(requestMessage string) LogrusMiddlewareBuilder {
	lrmb.requestMessage = requestMessage
	return lrmb
}

func (lrmb *logrusMiddlewareBuilder) ActivatedResponseData(activatedResponseData []string) LogrusMiddlewareBuilder {
	lrmb.activatedResponseData = activatedResponseData
	return lrmb
}

func (lrmb *logrusMiddlewareBuilder) CustomResponseData(customResponseData []ResponseData) LogrusMiddlewareBuilder {
	lrmb.customResponseData = customResponseData
	return lrmb
}

func (lrmb *logrusMiddlewareBuilder) ResponseMessage(responseMessage string) LogrusMiddlewareBuilder {
	lrmb.responseMessage = responseMessage
	return lrmb
}

func (lrmb *logrusMiddlewareBuilder) Build() LogrusMiddleware {
	// The default request data, will be overriden by the custom ones
	defaultRequestData := map[string]requestAccessor{
		"method":     func(r *http.Request) interface{} { return r.Method },
		"request":    func(r *http.Request) interface{} { return r.RequestURI },
		"remote":     func(r *http.Request) interface{} { return r.RemoteAddr },
		"referer":    func(r *http.Request) interface{} { return r.Referer() },
		"user-agent": func(r *http.Request) interface{} { return r.UserAgent() },
	}

	// The default response data, will be overriden by the custom ones
	defaultResponseData := map[string]responseAccessor{
		"status": func(r *LogrusResponse) interface{} { return r.Status },
		"size":   func(r *LogrusResponse) interface{} { return r.Size },
	}

	// Initialize all values that are not been by user
	lrmb.setDefaults(defaultRequestData, defaultResponseData)

	return &logrusMiddleware{
		// Build the request data from default ones activated by user and custom ones
		requestData:    lrmb.buildRequestData(defaultRequestData),
		requestMessage: lrmb.requestMessage,
		// Build the response data from default ones activated by user and custom ones
		responseData:    lrmb.buildResponseData(defaultResponseData),
		responseMessage: lrmb.responseMessage,
	}
}

// setDefaults : Initialize all values that are not been by user
func (lrmb *logrusMiddlewareBuilder) setDefaults(defaultRequestAccessors map[string]requestAccessor,
	defaultResponseAccessors map[string]responseAccessor) {

	// Fill the activated request fields to all if not set
	if lrmb.activatedRequestData == nil || len(lrmb.activatedRequestData) == 0 {
		for name := range defaultRequestAccessors {
			lrmb.activatedRequestData = append(lrmb.activatedRequestData, name)
		}
	}

	// Fill the activated response fields to all if not set
	if lrmb.activatedResponseData == nil || len(lrmb.activatedResponseData) == 0 {
		for name := range defaultResponseAccessors {
			lrmb.activatedResponseData = append(lrmb.activatedResponseData, name)
		}
	}

	// Fill the custom request fields to empty if not set
	if lrmb.customRequestData == nil {
		lrmb.customRequestData = []RequestData{}
	}

	// Fill the custom response fields to empty if not set
	if lrmb.customResponseData == nil {
		lrmb.customResponseData = []ResponseData{}
	}

	// Fill the request message to default one if not set
	if lrmb.requestMessage == "" {
		lrmb.requestMessage = "Request received !"
	}

	// Fill the response message to default one if not set
	if lrmb.responseMessage == "" {
		lrmb.responseMessage = "Response sent !"
	}

	return
}

func (lrmb *logrusMiddlewareBuilder) buildRequestData(accessorsMap map[string]requestAccessor) (res map[string]requestAccessor) {
	res = make(map[string]requestAccessor)

	// Search the activated request data in the defaults
	for _, field := range lrmb.activatedRequestData {
		if val, ok := accessorsMap[field]; ok {
			res[field] = val
		} else {
			logrus.Warnf("Field %s cannot be activated because is not available.", field)
		}
	}

	// We consider that the custom request data are always activated and even override the default ones
	for _, data := range lrmb.customRequestData {
		res[data.name] = data.accessor
	}

	return
}

func (lrmb *logrusMiddlewareBuilder) buildResponseData(accessorsMap map[string]responseAccessor) (res map[string]responseAccessor) {
	res = make(map[string]responseAccessor)

	// Search the activated response data in the defaults
	for _, field := range lrmb.activatedResponseData {
		if val, ok := accessorsMap[field]; ok {
			res[field] = val
		} else {
			logrus.Warnf("Field %s cannot be activated because is not available.", field)
		}
	}

	// We consider that the custom request data are always activated and even override the default ones
	for _, data := range lrmb.customResponseData {
		res[data.name] = data.accessor
	}
	return
}
