package model

type QueryRequest struct {
	Sql    string `json:"sql"`
	Params []any  `json:"params"`
}

type QueryResponse struct {
	ResponseCode    int              `json:"responseCode,omitempty"`
	ResponseMessage string           `json:"responseMessage,omitempty"`
	Data            []map[string]any `json:"data"`
}
