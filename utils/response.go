package utils

import "net/http"

func Response(w http.ResponseWriter, statusCode int, output []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(output)
}
