package api

import (
	"net/http"
	"testing"
)

const (
	errorOnNonEmptyBody = "default body should be empty"
	errorOnStatusCode   = "wrong http status code: %+v"
	errorOnBody         = "unexpected data: %+v"
	errorOnHeader       = "wrong header: %+v"
	expectedJsonBody    = `{"data":"hello world"}`
)

func TestNewResponse(t *testing.T) {
	resp := NewResponse()
	if resp.Status != http.StatusOK {
		t.Errorf(errorOnStatusCode, resp.Status)
	}
	if len(resp.Body) > 0 {
		t.Error(errorOnNonEmptyBody)
	}
	if len(resp.Headers) > 0 {
		t.Error(errorOnNonEmptyBody)
	}
}

func TestResponseWrite(t *testing.T) {
	resp := NewResponse()
	w := newMockResponseWriter()
	resp.Write(w)
	if w.statusCode != http.StatusOK {
		t.Errorf(errorOnStatusCode, w.statusCode)
	}
}

func TestWithStatus(t *testing.T) {
	resp := NewResponse(WithStatus(http.StatusCreated))
	w := newMockResponseWriter()
	resp.Write(w)
	if w.statusCode != http.StatusCreated {
		t.Errorf(errorOnStatusCode, w.statusCode)
	}
}

func TestWithBody(t *testing.T) {
	html := "<html><body><p>>Hello World</p></body></html>"
	resp := NewResponse(WithBody([]byte(html)))
	w := newMockResponseWriter()
	resp.Write(w)
	if w.statusCode != http.StatusOK {
		t.Errorf(errorOnStatusCode, w.statusCode)
	}
	if string(w.data) != html {
		t.Errorf(errorOnBody, w.data)
	}
}

func TestWithHeader(t *testing.T) {
	contentType := "multipart/form-data"
	resp := NewResponse(WithHeader(HeaderContentType, contentType))
	w := newMockResponseWriter()
	resp.Write(w)
	if w.statusCode != http.StatusOK {
		t.Errorf(errorOnStatusCode, w.statusCode)
	}
	if w.header.Get(HeaderContentType) != contentType {
		t.Errorf(errorOnHeader, w.header.Get(HeaderContentType))
	}
}

func TestWithContentTypeJson(t *testing.T) {
	resp := NewResponse(WithContentTypeJson())
	w := newMockResponseWriter()
	resp.Write(w)
	if w.statusCode != http.StatusOK {
		t.Errorf(errorOnStatusCode, w.statusCode)
	}
	if w.header.Get(HeaderContentType) != TypeApplicationJson {
		t.Errorf(errorOnHeader, w.header.Get(HeaderContentType))
	}
}

func TestWithJsonBody(t *testing.T) {
	type message struct {
		Data string `json:"data"`
	}
	mesg := message{Data: "hello world"}
	resp := NewResponse(WithJsonBody(mesg))
	w := newMockResponseWriter()
	resp.Write(w)
	if w.statusCode != http.StatusOK {
		t.Errorf(errorOnStatusCode, w.statusCode)
	}
	if w.header.Get(HeaderContentType) != TypeApplicationJson {
		t.Errorf(errorOnHeader, w.header.Get(HeaderContentType))
	}
	if string(w.data) != expectedJsonBody {
		t.Errorf(errorOnBody, string(w.data))
	}
}

func TestNewJsonResponse(t *testing.T) {
	type message struct {
		Data string `json:"data"`
	}
	mesg := message{Data: "hello world"}
	resp := NewJsonResponse(mesg)
	w := newMockResponseWriter()
	resp.Write(w)
	if w.statusCode != http.StatusOK {
		t.Errorf(errorOnStatusCode, w.statusCode)
	}
	if w.header.Get(HeaderContentType) != TypeApplicationJson {
		t.Errorf(errorOnHeader, w.header.Get(HeaderContentType))
	}
	if string(w.data) != expectedJsonBody {
		t.Errorf(errorOnBody, string(w.data))
	}
}

func newMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{
		header: http.Header{},
	}
}

type MockResponseWriter struct {
	header     http.Header
	data       []byte
	statusCode int
}

func (w *MockResponseWriter) Header() http.Header {
	return w.header
}

func (w *MockResponseWriter) Write(d []byte) (int, error) {
	w.data = d
	return len(d), nil
}

func (w *MockResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}
