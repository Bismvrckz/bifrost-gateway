package model

type AuthRequest struct {
	GrantType string `json:"grantType"`
}

type AuthResponse struct {
	ResponseCode    int    `json:"responseCode,omitempty"`
	Token           string `json:"token,omitempty"`
	ResponseMessage string `json:"responseMessage,omitempty"`
}
