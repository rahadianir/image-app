package xhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"image-app/internal/pkg/pagination"
	"io"
	"net/http"
	"reflect"
)

type BaseResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type BaseListResponse struct {
	Error    string              `json:"error,omitempty"`
	Message  string              `json:"message,omitempty"`
	Data     any                 `json:"data,omitempty"`
	Metadata pagination.Metadata `json:"metadata,omitempty"`
}

func BindJSONRequest(request *http.Request, destination any) error {
	defer request.Body.Close()

	if reflect.ValueOf(destination).Kind() != reflect.Ptr {
		return errors.New("destination has to be pointer")
	}

	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bodyBytes, &destination)
}

func SendJSONResponse(w http.ResponseWriter, data any, code int) {
	dj, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s", dj)
}
