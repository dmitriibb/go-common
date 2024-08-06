package httputils

import (
	"encoding/json"
	"github.com/dmitriibb/go-common/logging"
	"github.com/dmitriibb/go-common/restaurant-common/model"
	"net/http"
)

func ReturnResponseWithError(w http.ResponseWriter, statusCode int, logger logging.Logger, message string) {
	if statusCode < 300 {
		panic("Http error status code must be more or equal 300")
	}
	logger.Error(message)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.CommonErrorResponse{
		Type:    model.CommonErrorTypeInvalidData,
		Message: message,
	})
}
