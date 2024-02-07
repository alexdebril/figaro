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

func WithContentTypeJson() Option {
	return WithHeader(HeaderContentType, TypeApplicationJson)
}

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

func WithBody(body []byte) Option {
	return BodyOption{body: body}
}

type BodyOption struct {
	body []byte
}

func (p BodyOption) Apply(response *Response) {
	response.Body = p.body
}

func WithJsonBody(v any) Option {
	return JsonBodyOption{
		data:   v,
		Option: WithContentTypeJson(),
	}
}

type JsonBodyOption struct {
	data any
	Option
}

type Response struct {
	Status  int
	Headers map[string]string
	Body    []byte
}

func (jp JsonBodyOption) Apply(response *Response) {
	jp.Option.Apply(response)
	j, err := json.Marshal(jp.data)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Body = []byte(fmt.Sprintf("{\"error\": \"%s\"}", err.Error()))
		return
	}
	response.Body = j
}

func NewResponse(opts ...Option) *Response {
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

func (r *Response) Write(w http.ResponseWriter) {
	for name, value := range r.Headers {
		w.Header().Set(name, value)
	}
	w.WriteHeader(r.Status)
	_, _ = w.Write(r.Body)
}
