package sending

import (
	"encoding/json"
	"net/http"
)

const (
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
)

type errorMsg struct {
	ErrMsg string `json:"error"`
}

func SendErrorMsg(w http.ResponseWriter, status int, msg string) {
	body, err := json.Marshal(errorMsg{ErrMsg: msg})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentTypeHeader, contentTypeJSON)
	w.WriteHeader(status)

	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Write bytes to w and set json Content-Typee header
func SendJSONBytes(w http.ResponseWriter, status int, body []byte) {
	w.Header().Add(contentTypeHeader, contentTypeJSON)
	w.WriteHeader(status)

	_, err := w.Write(body)
	if err != nil {
		SendErrorMsg(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func JSONMarshallAndSend(w http.ResponseWriter, status int, obj any) {
	body, err := json.Marshal(obj)
	if err != nil {
		SendErrorMsg(w, http.StatusInternalServerError, err.Error())
		return
	}
	SendJSONBytes(w, status, body)
}
