package models

import "time"

type ChannelType string

const (
	Email ChannelType = "EMAIL"
	Slack ChannelType = "SLACK"
	InApp ChannelType = "IN_APP"
)

type NotificationRequest struct {
	UserID       string                 `json:"user_id"`
	Channel      ChannelType            `json:"channel"`
	TemplateName string                 `json:"template_name"`
	CustomBody   string                 `json:"custom_body"`
	Data         map[string]interface{} `json:"data"`
	ScheduledAt  *time.Time             `json:"scheduled_at"`
}
