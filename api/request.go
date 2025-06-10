package api

import (
	"encoding/json"
	"io"
	"net/http"
)

func ParseJsonRequest[T any](r *http.Request) (*T, error) {
	body, err := r.GetBody()
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var obj *T
	err = json.Unmarshal(b, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func BadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	resp := NewJsonResponse(&Response{
		Status: http.StatusBadRequest,
		Body:   []byte(`{"error": "` + err.Error() + `"}`),
	})
	resp.Write(w)
}
