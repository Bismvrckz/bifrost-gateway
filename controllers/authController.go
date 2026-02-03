package controller

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	model "middlewareApi/models"
	utils "middlewareApi/utils"
	"net/http"
	"strings"
)

func AuthController(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-Api-Key")
	contentType := r.Header.Get("Content-Type")
	authorization := r.Header.Get("Authorization")

	if apiKey == "" || contentType == "" || authorization == "" {
		utils.WriteJSONResponse(w, http.StatusBadRequest, model.AuthResponse{ResponseCode: 4002100, ResponseMessage: "Invalid Headers"})
		return
	}

	if !strings.HasPrefix(authorization, "Basic ") {
		utils.WriteJSONResponse(w, http.StatusBadRequest, model.AuthResponse{ResponseCode: 4002100, ResponseMessage: "Invalid Authorization header format"})
		return
	}

	encodedCredentials := strings.TrimPrefix(authorization, "Basic ")

	decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, model.AuthResponse{ResponseCode: 4002100, ResponseMessage: "Invalid base64 in Authorization header"})
		return
	}

	credentials := strings.SplitN(string(decodedCredentials), ":", 2)
	if len(credentials) != 2 {
		utils.WriteJSONResponse(w, http.StatusBadRequest, model.AuthResponse{ResponseCode: 4002100, ResponseMessage: "Invalid credentials format"})
		return
	}

	username := credentials[0]
	password := credentials[1]

	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedPassword := strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))

	strTosign := username + ":" + hashedPassword

	newAuthString := "Basic " + base64.StdEncoding.EncodeToString([]byte(strTosign))

	spResponseJSON, err := utils.GetToken(contentType, apiKey, strings.ToUpper(newAuthString))
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, model.AuthResponse{ResponseCode: 5002100, ResponseMessage: "Internal server error"})
		return
	}

	var resCallSp model.AuthResponse
	if err := json.Unmarshal([]byte(spResponseJSON), &resCallSp); err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, model.AuthResponse{ResponseCode: 5002100, ResponseMessage: "Internal server error on unmarshal"})
		return
	}

	if resCallSp.ResponseCode != 2002100 {
		utils.WriteJSONResponse(w, http.StatusUnauthorized, resCallSp)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, resCallSp)
}
