package internal

import "net/http"

func WriteInternalServerErrorResponse(w http.ResponseWriter, err error) {
	WriteJSON(w, http.StatusInternalServerError, "Oops, something went wrong!", nil, err)
}

func WriteBadRequestResponse(w http.ResponseWriter, err error) {
	WriteJSON(w, http.StatusBadRequest, "Bad Request", nil, err)
}

func WriteConflictResponse(w http.ResponseWriter, status string, err error) {
	WriteJSON(w, http.StatusConflict, status, nil, err)
}

func WriteNotFoundResponse(w http.ResponseWriter, err error) {
	WriteJSON(w, http.StatusNotFound, "Not Found", nil, err)
}