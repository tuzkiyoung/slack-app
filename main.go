package main

import (
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
	"net/http"
	"xd-infra-slack/config"
	"xd-infra-slack/drivers"
	"xd-infra-slack/handlers"
	"xd-infra-slack/server"
)

var c config.Config

func init() {
	if err := c.GetConfig(); err != nil {
		panic(err.Error())
	}
}

func main() {
	// init mysql client
	drivers.InitMysqlClient(&c)
	// init Slack app client
	client := drivers.InitSlackClient(&c)

	var n server.Notificator = &server.AlertData{}

	// simple http server for listening alerts
	go func() {
		http.HandleFunc("/api/v1/alert/", server.Ht)
		http.ListenAndServe(":8090", nil)
	}()

	// single goroutine for work
	go func() {
		for alert := range server.AlertChan {
			msgTs := n.Post(c.ChanID, client)
			server.Create(msgTs, alert, drivers.Db)
			if err := n.Call(&c); err != nil {
				log.Printf("Failed to make phone call,%v", err)
			}
		}
	}()

	socketmodeHandler := socketmode.NewSocketmodeHandler(client)

	socketmodeHandler.Handle(socketmode.EventTypeConnecting, middlewareConnecting)
	socketmodeHandler.Handle(socketmode.EventTypeConnectionError, middlewareConnectionError)
	socketmodeHandler.Handle(socketmode.EventTypeConnected, middlewareConnected)
	socketmodeHandler.Handle(socketmode.EventTypeHello, middlewareDoNothing)
	socketmodeHandler.Handle(socketmode.EventTypeIncomingError, middlewareDoNothing)

	//\\ EventTypeEventsAPI //\\
	// Handle all EventsAPI
	socketmodeHandler.Handle(socketmode.EventTypeEventsAPI, middlewareEventsAPI)

	// Handle a specific event from EventsAPI
	socketmodeHandler.HandleEvents(slackevents.AppMention, middlewareAppMentionEvent)

	//\\ EventTypeInteractive //\\
	// Handle all Interactive Events
	socketmodeHandler.Handle(socketmode.EventTypeInteractive, middlewareInteractive)

	// Handle a specific Interaction
	socketmodeHandler.HandleInteraction(slack.InteractionTypeBlockActions, middlewareInteractionTypeBlockActions)

	// Handle click About button in App Homepage
	socketmodeHandler.HandleInteractionBlockAction("git", handlers.MiddlewareServiceDetailModal)
	socketmodeHandler.HandleInteractionBlockAction("nas", handlers.MiddlewareServiceDetailModal)

	// Handle triggered alert messages
	socketmodeHandler.HandleInteractionBlockAction("act_ack", handlers.MiddlewareAlertAck)
	socketmodeHandler.HandleInteractionBlockAction("act_resolve", handlers.MiddlewareAlertResolved)

	// Handle all SlashCommand
	socketmodeHandler.Handle(socketmode.EventTypeSlashCommand, middlewareSlashCommand)
	socketmodeHandler.HandleSlashCommand("/fly", middlewareSlashCommand)

	// socketmodeHandler.HandleDefault(middlewareDefault)

	// Handle App home opened
	socketmodeHandler.HandleEvents(slackevents.AppHomeOpened, handlers.MiddlewareAppHomeOpened)

	socketmodeHandler.RunEventLoop()
}

func middlewareConnecting(evt *socketmode.Event, client *socketmode.Client) {
	log.Println("Connecting to Slack with Socket Mode...")
}

func middlewareConnectionError(evt *socketmode.Event, client *socketmode.Client) {
	log.Println("Connection failed. Retrying later...")
}

func middlewareConnected(evt *socketmode.Event, client *socketmode.Client) {
	log.Println("Connected to Slack with Socket Mode.")
}

func middlewareDoNothing(evt *socketmode.Event, client *socketmode.Client) {
	//fmt.Println("Hallo, meine Freunde! Das ist XD Infrastruktur mannschaft!")
	if e, ok := evt.Data.(slackevents.EventsAPIEvent); !ok {
		//fmt.Println(reflect.TypeOf(evt.Data))
		log.Println("middlewareDoNothing")
	} else {
		fmt.Println(e.Data)
		fmt.Println(e.InnerEvent.Data)
	}

}

func middlewareEventsAPI(evt *socketmode.Event, client *socketmode.Client) {
	log.Println("middlewareEventsAPI")
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return
	}

	fmt.Printf("Event received: %+v\n", eventsAPIEvent)

	client.Ack(*evt.Request)

	switch eventsAPIEvent.Type {
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			fmt.Printf("We have been mentionned in %v", ev.Channel)
			_, _, err := client.Client.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			if err != nil {
				fmt.Printf("failed posting message: %v", err)
			}
		case *slackevents.MemberJoinedChannelEvent:
			fmt.Printf("user %q joined to channel %q", ev.User, ev.Channel)
		}
	default:
		client.Debugf("unsupported Events API event received")
	}
}

func middlewareAppMentionEvent(evt *socketmode.Event, client *socketmode.Client) {

	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return
	}

	client.Ack(*evt.Request)

	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.AppMentionEvent)
	if !ok {
		fmt.Printf("Ignored %+v\n", ev)
		return
	}

	fmt.Printf("We have been mentionned in %v\n", ev.Channel)
	_, _, err := client.Client.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
	if err != nil {
		fmt.Printf("failed posting message: %v", err)
	}
}

func middlewareInteractive(evt *socketmode.Event, client *socketmode.Client) {
	callback, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return
	}

	fmt.Printf("Interaction received: %+v\n", callback)

	var payload interface{}

	switch callback.Type {
	case slack.InteractionTypeBlockActions:
		// See https://api.slack.com/apis/connections/socket-implement#button
		client.Debugf("button clicked!")
	case slack.InteractionTypeShortcut:
	case slack.InteractionTypeViewSubmission:
		// See https://api.slack.com/apis/connections/socket-implement#modal
	case slack.InteractionTypeDialogSubmission:
	default:

	}

	client.Ack(*evt.Request, payload)
}

func middlewareInteractionTypeBlockActions(evt *socketmode.Event, client *socketmode.Client) {
	client.Debugf("button clicked!")
}

func middlewareSlashCommand(evt *socketmode.Event, client *socketmode.Client) {
	cmd, ok := evt.Data.(slack.SlashCommand)
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return
	}

	client.Debugf("Slash command received: %+v", cmd)

	payload := map[string]interface{}{
		"blocks": []slack.Block{
			slack.NewSectionBlock(
				&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: "foo",
				},
				nil,
				slack.NewAccessory(
					slack.NewButtonBlockElement(
						"",
						"somevalue",
						&slack.TextBlockObject{
							Type: slack.PlainTextType,
							Text: "bar",
						},
					),
				),
			),
		}}
	client.Ack(*evt.Request, payload)
}

func middlewareDefault(evt *socketmode.Event, client *socketmode.Client) {
	// fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
}
