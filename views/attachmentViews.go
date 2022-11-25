package views

import (
	"embed"
	"encoding/json"
	"github.com/slack-go/slack"
	"io"
	"log"
)

//go:embed assets/*
var slackAssets embed.FS

func AlertTriggered(data interface{}) slack.Attachment {
	tpl := renderTemplate(slackAssets, "assets/alert_triggered.json", data)
	view := slack.Attachment{}
	str, err := io.ReadAll(&tpl)
	if err != nil {
		log.Printf("Unable to read view `alert_triggered`: %v", err)
	}
	if err := json.Unmarshal(str, &view); err != nil {
		log.Printf("AlertTriggered,%v\n", err)
		return slack.Attachment{}
	}
	return view
}

func AlertAcknowledged(data interface{}) slack.Attachment {
	tpl := renderTemplate(slackAssets, "assets/alert_ack.json", data)
	str, err := io.ReadAll(&tpl)
	if err != nil {
		log.Printf("Unable to read view `alert_ack`: %v", err)
	}
	view := slack.Attachment{}
	if err := json.Unmarshal(str, &view); err != nil {
		log.Printf("AlertAcknowledged,%v\n", err)
		return slack.Attachment{}
	}
	return view
}

func AlertResolved(data interface{}) slack.Attachment {
	tpl := renderTemplate(slackAssets, "assets/alert_resolved.json", data)
	str, err := io.ReadAll(&tpl)
	if err != nil {
		log.Printf("Unable to read view `alert_resolved`: %v", err)
	}
	view := slack.Attachment{}
	if err := json.Unmarshal(str, &view); err != nil {
		log.Printf("AlertResolved,%v\n", err)
		return slack.Attachment{}
	}
	return view
}
