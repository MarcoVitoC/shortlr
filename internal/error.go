package internal

import "net/http"

func WriteInternalServerErrorResponse(w http.ResponseWriter, err error) {
	WriteJSONResponse(w, http.StatusInternalServerError, "Oops, something went wrong!", nil, err)
}

func WriteBadRequestResponse(w http.ResponseWriter, err error) {
	WriteJSONResponse(w, http.StatusBadRequest, "Bad Request", nil, err)
}

func WriteConflictResponse(w http.ResponseWriter, err error) {
	WriteJSONResponse(w, http.StatusConflict, "Conflict Request", nil, err)
}

func WriteNotFoundResponse(w http.ResponseWriter, err error) {
	WriteJSONResponse(w, http.StatusNotFound, "Not Found", nil, err)
}