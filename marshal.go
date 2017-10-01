package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type BasicResponse struct {
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, code int, err error) {
	bodyMap := gin.H{"error": gin.H{"message": err.Error()}}
	writeJSON(w, code, bodyMap)
}

func writeJSON(w http.ResponseWriter, code int, obj map[string]interface{}) {
	writeObject(w, code, obj)
}

func writeMessage(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, 200, map[string]interface{}{
		"message": msg,
	})
}

func writeObject(w http.ResponseWriter, code int, obj interface{}) {
	jsonValue, _ := json.Marshal(obj)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write([]byte(jsonValue))
	if err != nil {
		logrus.WithError(err).Errorln("couldn't write error response")
	}
}

func parseJSON(w http.ResponseWriter, r *http.Request, t interface{}) error {
	err := parseJSONReader(r.Body, t)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid request body, bad JSON: %v", err))
		return err
	}
	return nil
}

func parseJSONReader(r io.Reader, t interface{}) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(t)
	return err
}

func bytesToJSON(bs []byte) (string, error) {
	return toJSON(string(bs))
}

func toJSON(v interface{}) (string, error) {
	jsonValue, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(jsonValue), nil
}
