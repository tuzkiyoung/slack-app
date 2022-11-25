package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"io"
	"log"
	"net/http"
	"xd-infra-slack/config"
	"xd-infra-slack/views"
)

type Notificator interface {
	Receive(w http.ResponseWriter, r *http.Request)
	Post(chanID string, client *socketmode.Client) string
	Call(c *config.Config) error
}

type AlertData struct {
	Title        string `json:"title"`
	Platform     string `json:"platform"`
	Region       string `json:"region"`
	InstanceId   string `json:"instanceId"`
	InstanceName string `json:"instanceName"`
	AlertId      string `json:"alertId"`
	Assigned     string `json:"assigned"`
	Service      string `json:"service"`
	Metric       string `json:"metric"`
	Value        string `json:"value"`
	Priority     string `json:"priority"`
	//Others       map[string]interface{} `json:"others"`
}

var AlertChan = make(chan AlertData)

func Ht(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}

	a := AlertData{}
	if err := json.Unmarshal(b, &a); err != nil {
		log.Printf("Failed to assemble data %v", err)
		fmt.Fprintf(w, "%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	a.Receive2()
}

func (a AlertData) Receive2() {
	AlertChan <- a
}

func (a AlertData) Receive(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}

	//a := AlertData{}
	if err := json.Unmarshal(b, &a); err != nil {
		log.Printf("Failed to assemble data %v", err)
		fmt.Fprintf(w, "%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	AlertChan <- a
}

func (a AlertData) Post(chanID string, client *socketmode.Client) string {
	payload := views.AlertTriggered(a)
	_, msgTs, err := client.PostMessage(chanID, slack.MsgOptionAttachments(payload))
	if err != nil {
		log.Printf("ERROR MiddlewarePostAlert: %v", err)
		return ""
	}
	return msgTs
}

func (a AlertData) Call(c *config.Config) error {
	payload, err := json.Marshal(a)
	if err != nil {
		return err
	}

	if _, err := http.Post(c.Url, "application/json", bytes.NewReader(payload)); err != nil {
		return err
	}
	return nil
}
