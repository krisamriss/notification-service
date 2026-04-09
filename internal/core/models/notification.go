package models

import "time"

type ChannelType string

const (
	Email ChannelType = "EMAIL"
	Slack ChannelType = "SLACK"
	InApp ChannelType = "IN_APP"
	SMS   ChannelType = "SMS"
)

type NotificationRequest struct {
	UserID       string                 `json:"user_id"`
	Channel      ChannelType            `json:"channel"`
	TemplateName string                 `json:"template_name"` // For predefined
	CustomBody   string                 `json:"custom_body"`   // For user-defined
	Data         map[string]interface{} `json:"data"`          // Dynamic data for templates
	ScheduledAt  *time.Time             `json:"scheduled_at"`  // Nil means instant
}
