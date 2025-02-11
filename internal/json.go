package internal

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int     `json:"code"`
	Status  string  `json:"status"`
	Data    any 	`json:"data"`
	Error   string  `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, code int, status string, data interface{}, err error) error {
	if data == nil {
		data = []interface{}{}
	}
	
	response := Response{
		Code: code,
		Status: status,
		Data: data,
	}

	if err != nil {
		response.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(response)
}

func ReadJson(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := int64(1_048_576) // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func WriteOK(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusOK, "OK", data, nil)
}