package api

import (
	"encoding/json"
	"io"
	"net/http"
)

func ParseJsonRequest[T any](r *http.Request) (*T, error) {
	if r.Body == nil {
		return nil, http.ErrBodyReadAfterClose
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var obj *T
	err = json.Unmarshal(b, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func BadRequest(w http.ResponseWriter, err error) {
	resp := NewJsonResponse(NewErrorResponse(err), WithStatus(http.StatusBadRequest))
	resp.Write(w)
}
