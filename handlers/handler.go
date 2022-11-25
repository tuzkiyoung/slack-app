package handlers

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
	"reflect"
	"xd-infra-slack/drivers"
	"xd-infra-slack/server"
	"xd-infra-slack/views"
)

func MiddlewareAppHomeOpened(evt *socketmode.Event, client *socketmode.Client) {
	var user string
	payload := views.AppHomeTabView()
	evtApi, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		log.Printf("ERROR converting event to slackevents.EventsAPIEvent")
	}
	evtAppHomeOpened, ok := evtApi.InnerEvent.Data.(slackevents.AppHomeOpenedEvent)

	if !ok {
		log.Printf("ERROR converting inner event to slackevents.AppHomeOpenedEvent")
		//Patch the fact that we are not able to cast evt_api.InnerEvent.Data to AppHomeOpenedEvent
		user = reflect.ValueOf(evtApi.InnerEvent.Data).Elem().FieldByName("User").Interface().(string)
	} else {
		user = evtAppHomeOpened.User
	}

	_, err := client.PublishView(user, payload, "")
	if err != nil {
		return
	}
}

func MiddlewareServiceDetailModal(evt *socketmode.Event, client *socketmode.Client) {
	// we need to cast our socketmode.Event
	interaction := evt.Data.(slack.InteractionCallback)

	// Make sure to respond to the server to avoid an error
	client.Ack(*evt.Request)

	// create the view using block-kit
	payload := views.AppServiceDetailModal(interaction.ActionCallback.BlockActions[0].Value)

	// Open Modal (13)
	_, err := client.OpenView(interaction.TriggerID, payload)

	//Handle errors
	if err != nil {
		log.Printf("ERROR openCreateStickieNoteModal: %v", err)
	}
}

func MiddlewareAlertAck(evt *socketmode.Event, client *socketmode.Client) {
	interaction := evt.Data.(slack.InteractionCallback)
	client.Ack(*evt.Request)
	server.Update(interaction.Message.Msg.Timestamp, "Acknowledged", drivers.Db)
	data := server.Retrieve(interaction.Message.Msg.Timestamp, drivers.Db)
	payload := views.AlertAcknowledged(data)
	if _, _, err := client.PostMessage(interaction.Container.ChannelID, slack.MsgOptionAttachments(payload), slack.MsgOptionReplaceOriginal(interaction.ResponseURL)); err != nil {
		log.Printf("ERROR MiddlewareAlertAck: %v", err)
	}
}

func MiddlewareAlertResolved(evt *socketmode.Event, client *socketmode.Client) {
	interaction := evt.Data.(slack.InteractionCallback)
	client.Ack(*evt.Request)
	server.Update(interaction.Message.Msg.Timestamp, "Resolved", drivers.Db)
	data := server.Retrieve(interaction.Message.Msg.Timestamp, drivers.Db)
	payload := views.AlertResolved(data)
	if _, _, err := client.PostMessage(interaction.Container.ChannelID, slack.MsgOptionAttachments(payload), slack.MsgOptionReplaceOriginal(interaction.ResponseURL)); err != nil {
		log.Printf("ERROR MiddlewareAlertAck: %v", err)
	}
}

//func MiddlewarePostAlert(chanID string, data server.AlertData, client *socketmode.Client) {
//	payload := views.AlertTriggered(data)
//	_, _, err := client.PostMessage(chanID, slack.MsgOptionAttachments(payload))
//	if err != nil {
//		log.Printf("ERROR MiddlewarePostAlert: %v", err)
//	}
//}

func MiddlewareGreeting(evt *socketmode.Event, client *socketmode.Client) {}
