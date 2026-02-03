package model

type SpRequest struct {
	SpName   string `json:"spName"`
	SpParams []any  `json:"spParams"`
}

type SpResponse struct {
	ResponseCode    int              `json:"responseCode,omitempty"`
	ResponseMessage string           `json:"responseMessage,omitempty"`
	Data            []map[string]any `json:"data"`
}
