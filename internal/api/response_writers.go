package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, status int) error {
	msg := ""
	switch status {
	case http.StatusMethodNotAllowed:
		msg = "method not allowed"
	case http.StatusBadRequest:
		msg = "bad request"
	case http.StatusInternalServerError:
		msg = "internal server error"
	default:
		status = http.StatusInternalServerError
		msg = "internal server error"
	}
	w.WriteHeader(status)
	_, err := fmt.Fprint(w, msg)
	if err != nil {
		return fmt.Errorf("WriteErrorResponse: %s", err)
	}
	return nil
}

func WriteErrorResponsePayload(w http.ResponseWriter, statusCode int, code int, message string, details *Detailes) error {
	if details == nil {
		details = &Detailes{}
	}
	p := ErrorPayload{
		Code:     code,
		Message:  message,
		Detailes: details,
	}
	w.WriteHeader(statusCode)
	if data, err := json.Marshal(p); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError)
		return fmt.Errorf("WriteErrorResponsePayload.%s", err.Error())
	} else {
		if _, err := w.Write(data); err != nil {
			return fmt.Errorf("WriteErrorResponsePayload.Write: %s", err.Error())
		}
	}
	return nil
}

func WriteJson(w http.ResponseWriter, data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError)
		return fmt.Errorf("WriteJson.json.Marshal: %s", err.Error())
	}
	if _, err := w.Write(jsonData); err != nil {
		return fmt.Errorf("WriteJson.Write: %s", err.Error())
	}
	return nil
}
