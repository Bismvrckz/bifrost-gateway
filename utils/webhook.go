package utils

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	model "middlewareApi/models"
	"net/http"
	"time"
)

type DiscordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type DiscordEmbed struct {
	Title     string         `json:"title"`
	Fields    []DiscordField `json:"fields"`
	Timestamp string         `json:"timestamp"`
}

type DiscordWebhookPayload struct {
	Username string         `json:"username"`
	Embeds   []DiscordEmbed `json:"embeds"`
}

func SendDiscordWebhook(apiKey string, req model.WebhookRequest) (err error) {
	// discordServerID := os.Getenv("VITE_DISCORD_SERVER_ID")
	// discordToken := os.Getenv("VITE_DISCORD_TOKEN")

	// if discordServerID == "" || discordToken == "" {
	// 	return fmt.Errorf("discord credentials not found in environment variables")
	// }

	var apiUser string
	var confUser string
	var confPass string
	var confPort sql.NullString
	var confHost sql.NullString
	var confDefault string

	err = dbMid.QueryRow("CALL GET_USER_API_DETAIL(?, 'webhook')", apiKey).Scan(&apiUser, &confUser, &confPass, &confPort, &confHost, &confDefault)
	if err != nil {
		log.Printf("Failed to execute SP GET_USER_API_DETAIL: %s", err.Error())
		return err
	}

	decryptedToken, err := Decrypt(confPass)
	if err != nil {
		log.Printf("Failed to decrypt token: %s", err.Error())
		return err
	}

	webhookUrl := fmt.Sprintf("https://discordapp.com/api/webhooks/%s/%s", confUser, decryptedToken)

	payload := DiscordWebhookPayload{
		Username: "Quasar Hook",
		Embeds: []DiscordEmbed{
			{
				Title: "New Contact 📬",
				Fields: []DiscordField{
					{
						Name:   "Name",
						Value:  req.WebhookContactName,
						Inline: true,
					},
					{
						Name:   "Email",
						Value:  req.WebhookContactEmail,
						Inline: true,
					},
					{
						Name:  "Message",
						Value: req.WebhookContactMessage,
					},
				},
				Timestamp: time.Now().Format(time.RFC3339),
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord webhook failed with status code: %d", resp.StatusCode)
	}

	Info("Successfully sent Discord webhook for contact: %s", req.WebhookContactName)
	return nil
}
