package api

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

type ResponseOption interface {
	Apply(response *Response)
}

func WithStatus(code int) ResponseOption {
	return StatusOption{code: code}
}

type StatusOption struct {
	code int
}

func (s StatusOption) Apply(response *Response) {
	response.Status = s.code
}

// WithContentTypeJson tells the Response its content-Type will be "application/json"
func WithContentTypeJson() ResponseOption {
	return WithHeader(HeaderContentType, TypeApplicationJson)
}

// WithHeader adds a header to the Response. Invoke WithHeader as many times as you have headers to set.
func WithHeader(name, value string) ResponseOption {
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

// WithBody sets the response's body.
func WithBody(body []byte) ResponseOption {
	return BodyOption{body: body}
}

type BodyOption struct {
	body []byte
}

func (p BodyOption) Apply(response *Response) {
	response.Body = p.body
}

// WithJsonBody takes an object meant to be translated to JSON inside the response's body and sets the content-type to "application/json"
func WithJsonBody(v any) ResponseOption {
	return JsonBodyOption{
		data:           v,
		ResponseOption: WithContentTypeJson(),
	}
}

type JsonBodyOption struct {
	data any
	ResponseOption
}

type Response struct {
	Status  int
	Headers map[string]string
	Body    []byte
}

func (jp JsonBodyOption) Apply(response *Response) {
	jp.ResponseOption.Apply(response)
	j, err := json.Marshal(jp.data)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Body = []byte(fmt.Sprintf("{\"error\": \"%s\"}", err.Error()))
		return
	}
	response.Body = j
}

// NewResponse creates a Response with default attributes:
// Status=200
// Body=empty byte slice
// Headers=empty hashmap
func NewResponse(opts ...ResponseOption) *Response {
	response := &Response{
		Status:  DefaultHttpStatus,
		Body:    make([]byte, 0),
		Headers: make(map[string]string),
	}
	for _, opt := range opts {
		opt.Apply(response)
	}
	return response
}

func NewJsonResponse(v any, opts ...ResponseOption) *Response {
	return NewResponse(append(opts, WithJsonBody(v))...)
}

// Write sends the Response's content to the http.ResponseWriter.
func (r *Response) Write(w http.ResponseWriter) {
	for name, value := range r.Headers {
		w.Header().Set(name, value)
	}
	w.WriteHeader(r.Status)
	_, _ = w.Write(r.Body)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{
		Error: err.Error(),
	}
}
