package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Err string `json:"error"`
}

type Response struct {
	Msg string `json:"msg"`
}

func SendData(w http.ResponseWriter, resp interface{}, statusCode int) error {
	respData, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(respData)

	return nil
}

func SendError(w http.ResponseWriter, resp interface{}, statusCode int) {
	respData, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(respData)

	log.Printf("Error response %#v", resp)
}
