package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	HeaderContentType   = "Content-Type"
	TypeApplicationJson = "application/json"
	DefaultHttpStatus   = http.StatusOK
)

type Option interface {
	Apply(response *Response)
}

func WithStatus(code int) Option {
	return StatusOption{code: code}
}

type StatusOption struct {
	code int
}

func (s StatusOption) Apply(response *Response) {
	response.Status = s.code
}

// WithContentTypeJson tells the Response its content-Type will be "application/json"
func WithContentTypeJson() Option {
	return WithHeader(HeaderContentType, TypeApplicationJson)
}

// WithHeader adds a header to the Response. Invoke WithHeader as many times as you have headers to set.
func WithHeader(name, value string) Option {
	return HeaderOption{
		name:  name,
		value: value,
	}
}

type HeaderOption struct {
	name  string
	value string
}

func (h HeaderOption) Apply(response *Response) {
	response.Headers[h.name] = h.value
}

// WithPayload sets the response's body.
func WithPayload(payload []byte) Option {
	return PayloadOption{payload: payload}
}

type PayloadOption struct {
	payload []byte
}

func (p PayloadOption) Apply(response *Response) {
	response.Payload = p.payload
}

// WithJsonPayload takes an object meant to be translated to JSON inside the response's body and sets the content-type to "application/json"
func WithJsonPayload(v any) Option {
	return JsonPayloadOption{
		data:   v,
		Option: WithContentTypeJson(),
	}
}

type JsonPayloadOption struct {
	data any
	Option
}

type Response struct {
	Status  int
	Headers map[string]string
	Payload []byte
}

func (jp JsonPayloadOption) Apply(response *Response) {
	jp.Option.Apply(response)
	j, err := json.Marshal(jp.data)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Payload = []byte(fmt.Sprintf("{\"error\": \"%s\"}", err.Error()))
		return
	}
	response.Payload = j
}

// NewResponse creates a Response with default attributes:
// Status=200
// Body=empty byte slice
// Headers=empty hashmap
func NewResponse(opts ...Option) *Response {
	response := &Response{
		Status:  DefaultHttpStatus,
		Payload: make([]byte, 0),
		Headers: make(map[string]string),
	}
	for _, opt := range opts {
		opt.Apply(response)
	}
	return response
}

// Write sends the Response's content to the http.ResponseWriter.
func (r *Response) Write(w http.ResponseWriter) {
	for name, value := range r.Headers {
		w.Header().Set(name, value)
	}
	w.WriteHeader(r.Status)
	_, _ = w.Write(r.Payload)
}
