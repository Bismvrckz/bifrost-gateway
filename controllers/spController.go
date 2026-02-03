package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	model "middlewareApi/models"
	utils "middlewareApi/utils"
	"net/http"
	"strings"
)

func SpController(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-KEY")
	contentType := r.Header.Get("Content-Type")
	authorization := r.Header.Get("Authorization")
	signature := r.Header.Get("X-SIGNATURE")

	if apiKey == "" || contentType == "" || authorization == "" {
		utils.WriteJSONResponse(w, http.StatusBadRequest, model.SpResponse{ResponseCode: 4002200, ResponseMessage: "Invalid Headers"})
		return
	}

	// Read the raw body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, model.SpResponse{ResponseCode: 5000000, ResponseMessage: "Failed to read request body"})
		return
	}
	defer r.Body.Close()

	// Restore the body so it can be read again
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	rawBody := string(bodyBytes)

	response, err := utils.ValidateAuth(contentType, apiKey, authorization, signature, rawBody)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, model.SpResponse{ResponseCode: 5000000, ResponseMessage: err.Error()})
		return
	}

	if response != "2002100" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write([]byte(response))
		return
	}

	if !strings.HasPrefix(authorization, "Bearer ") {
		utils.WriteJSONResponse(w, http.StatusBadRequest, model.SpResponse{ResponseCode: 4002200, ResponseMessage: "Invalid Authorization header format"})
		return
	}

	// encodedCredentials := strings.TrimPrefix(authorization, "Bearer ")

	// decodedBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
	// if err != nil {
	// 	utils.WriteJSONResponse(w, http.StatusBadRequest, model.CallSpResponse{ResponseCode: 4002200, ResponseMessage: "Invalid base64 in Authorization header"})
	// 	return
	// }
	// decodedCredentials := string(decodedBytes)

	var reqCallSp model.SpRequest
	if err := json.Unmarshal([]byte(rawBody), &reqCallSp); err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, model.SpResponse{ResponseCode: 5002200, ResponseMessage: "Internal server error"})
		return
	}

	fmt.Printf("SpName: %s\n", reqCallSp.SpName)
	fmt.Printf("SpParams: %v\n", reqCallSp.SpParams)

	resCallSp, err := utils.CallSp(apiKey, reqCallSp.SpName, reqCallSp.SpParams)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, model.SpResponse{ResponseCode: 5002200, ResponseMessage: "Internal server error"})
		return
	}

	// var resCallSp model.CallSpResponse
	// if err := json.Unmarshal([]byte(spResponseJSON), &resCallSp); err != nil {
	// 	utils.WriteJSONResponse(w, http.StatusInternalServerError, model.CallSpResponse{ResponseCode: 5002200, ResponseMessage: "Internal server error on unmarshal"})
	// 	return
	// }

	// if resCallSp.ResponseCode != 2002100 {
	// 	utils.WriteJSONResponse(w, http.StatusBadRequest, resCallSp)
	// 	return
	// }

	utils.Info("Successfully call sp: " + reqCallSp.SpName)

	if resCallSp == nil {
		resCallSp = make([]map[string]any, 0)
	}

	utils.WriteJSONResponse(w, http.StatusOK, model.SpResponse{ResponseCode: 2002200, ResponseMessage: "Success", Data: resCallSp})
}
