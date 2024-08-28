package webUtils

import "net/http"

func EnableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func ContentTypeJson(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func HandleOptionsRequest(w http.ResponseWriter, allowOrigins string, allowMethods string) {
	w.Header().Set("Access-Control-Allow-Origin", allowOrigins)
	w.Header().Set("Access-Control-Allow-Methods", allowMethods)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusNoContent)
}
