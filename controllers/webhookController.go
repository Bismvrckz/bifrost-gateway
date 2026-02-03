package controller

import (
	"bytes"
	"encoding/json"
	"io"
	model "middlewareApi/models"
	utils "middlewareApi/utils"
	"net/http"
	"strings"
)

func WebhookController(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-Api-Key")
	contentType := r.Header.Get("Content-Type")
	authorization := r.Header.Get("Authorization")
	signature := r.Header.Get("X-Signature")

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

	var reqSendWebhook model.WebhookRequest
	if err := json.Unmarshal([]byte(rawBody), &reqSendWebhook); err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, model.SpResponse{ResponseCode: 5002200, ResponseMessage: "Internal server error"})
		return
	}

	if err := utils.SendDiscordWebhook(apiKey, reqSendWebhook); err != nil {
		utils.Error("Failed to send discord webhook: %v", err)
		utils.WriteJSONResponse(w, http.StatusInternalServerError, model.SpResponse{ResponseCode: 5002200, ResponseMessage: "Failed to send webhook"})
		return
	}

	// resCallSp, err := utils.CallSp(apiKey, reqCallSp.SpName, reqCallSp.SpParams)
	// if err != nil {
	// 	utils.WriteJSONResponse(w, http.StatusInternalServerError, model.SpResponse{ResponseCode: 5002200, ResponseMessage: "Internal server error"})
	// 	return
	// }

	// var resCallSp model.CallSpResponse
	// if err := json.Unmarshal([]byte(spResponseJSON), &resCallSp); err != nil {
	// 	utils.WriteJSONResponse(w, http.StatusInternalServerError, model.CallSpResponse{ResponseCode: 5002200, ResponseMessage: "Internal server error on unmarshal"})
	// 	return
	// }

	// if resCallSp.ResponseCode != 2002100 {
	// 	utils.WriteJSONResponse(w, http.StatusBadRequest, resCallSp)
	// 	return
	// }

	utils.WriteJSONResponse(w, http.StatusOK, model.WebhookResponse{ResponseCode: 2002200, ResponseMessage: "Success"})
}
