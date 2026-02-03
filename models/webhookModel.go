package model

type WebhookRequest struct {
	WebhookContactName    string `json:"name"`
	WebhookContactEmail   string `json:"email"`
	WebhookContactMessage string `json:"message"`
}

type WebhookResponse struct {
	ResponseCode    int    `json:"responseCode,omitempty"`
	ResponseMessage string `json:"responseMessage,omitempty"`
}
