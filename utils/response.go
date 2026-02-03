package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, httpStatus int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	requestId := w.Header().Get("X-Request-ID")
	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"responseCode":5000000,"responseMessage":"encode error"}`))
		return
	}

	Response(fmt.Sprintf("%s | [RequestId=%s]", string(body), requestId))
	w.Write(body)
}
