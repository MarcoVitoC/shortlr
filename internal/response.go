package internal

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error,omitempty"`
}

func WriteJSONResponse(w http.ResponseWriter, code int, status string, data interface{}, err error) {
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
	json.NewEncoder(w).Encode(response)
}