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

func QueryController(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-KEY")
	contentType := r.Header.Get("Content-Type")
	authorization := r.Header.Get("Authorization")
	signature := r.Header.Get("X-SIGNATURE")

	if apiKey == "" || contentType == "" || authorization == "" {
		utils.WriteJSONResponse(w, http.StatusBadRequest, model.QueryResponse{ResponseCode: 4002200, ResponseMessage: "Invalid Headers"})
		return
	}

	// Read the raw body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, model.QueryResponse{ResponseCode: 5000000, ResponseMessage: "Failed to read request body"})
		return
	}
	defer r.Body.Close()

	// Restore the body so it can be read again
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	rawBody := string(bodyBytes)

	response, err := utils.ValidateAuth(contentType, apiKey, authorization, signature, rawBody)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, model.QueryResponse{ResponseCode: 5000000, ResponseMessage: err.Error()})
		return
	}

	if response != "2002100" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write([]byte(response))
		return
	}

	if !strings.HasPrefix(authorization, "Bearer ") {
		utils.WriteJSONResponse(w, http.StatusBadRequest, model.QueryResponse{ResponseCode: 4002200, ResponseMessage: "Invalid Authorization header format"})
		return
	}

	// encodedCredentials := strings.TrimPrefix(authorization, "Bearer ")

	// decodedBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
	// if err != nil {
	// 	utils.WriteJSONResponse(w, http.StatusBadRequest, model.CallQueryResponse{ResponseCode: 4002200, ResponseMessage: "Invalid base64 in Authorization header"})
	// 	return
	// }
	// decodedCredentials := string(decodedBytes)

	var reqQuery model.QueryRequest
	if err := json.Unmarshal([]byte(rawBody), &reqQuery); err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, model.QueryResponse{ResponseCode: 5002200, ResponseMessage: "Internal server error"})
		return
	}

	fmt.Printf("Sql: %s\n", reqQuery.Sql)
	fmt.Printf("Params: %+v\n", reqQuery.Params...)

	resExecQuery, err := utils.ExecuteQuery(apiKey, reqQuery.Sql, reqQuery.Params)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, model.QueryResponse{ResponseCode: 5002200, ResponseMessage: "Internal server error"})
		return
	}

	// var resCallSp model.CallQueryResponse
	// if err := json.Unmarshal([]byte(QueryResponseJSON), &resCallSp); err != nil {
	// 	utils.WriteJSONResponse(w, http.StatusInternalServerError, model.CallQueryResponse{ResponseCode: 5002200, ResponseMessage: "Internal server error on unmarshal"})
	// 	return
	// }

	// if resCallSp.ResponseCode != 2002100 {
	// 	utils.WriteJSONResponse(w, http.StatusBadRequest, resCallSp)
	// 	return
	// }

	utils.Info("Successfully executed query from DB: ")

	if resExecQuery == nil {
		resExecQuery = make([]map[string]any, 0)
	}

	utils.WriteJSONResponse(w, http.StatusOK, model.QueryResponse{ResponseCode: 2002200, ResponseMessage: "Success", Data: resExecQuery})
}
