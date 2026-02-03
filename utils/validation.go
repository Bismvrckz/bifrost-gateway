package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
)

func ValidateAuth(contentType, apiKey, authorization, reqSignature, strToSign string) (response string, err error) {
	if !strings.HasPrefix(authorization, "Bearer ") {
		err = errors.New("invalid token auth format")
		return response, err
	}

	token := strings.TrimPrefix(authorization, "Bearer ")

	secret, err := GetSecret(apiKey)
	if err != nil {
		return response, err
	}

	stringToSign := fmt.Sprintf("%s%s", secret, strToSign)

	hasher := sha256.New()
	hasher.Write([]byte(stringToSign))
	computeSignature := hex.EncodeToString(hasher.Sum(nil))

	err = dbMid.QueryRow("CALL VALIDATE_AUTH(?, ?, ?, ?, ?)", contentType, apiKey, token, reqSignature, computeSignature).Scan(&response)
	if err != nil {
		log.Printf("Failed to execute stored procedure VALIDATE_AUTH: %s", err.Error())
		return response, err
	}

	return response, nil
}
