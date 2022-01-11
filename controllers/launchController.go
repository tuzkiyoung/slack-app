package controllers

import (
	"log"
	"slack-app/drivers"
	"slack-app/module"
	"slack-app/views"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// AppLaunchController We create a sctucture to let us use dependency injection
type AppLaunchController struct {
	EventHandler *drivers.SocketmodeHandler
}

var (
	// Declare a parent message ChannelID pointer variable
	responseChannelID *string
	// Declare a parent message TimeStamp pointer variable
	responseTimeStamp *string
	// Declare a appDeploymentInfo pointer variable
	appArgs *views.AppDeploymentInfo
)

func NewAppLaunchController(eventhandler *drivers.SocketmodeHandler) AppLaunchController {
	// we need to cast our socketmode.Event into a SlashCommand
	c := AppLaunchController{
		EventHandler: eventhandler,
	}

	// Register callback for the command /app
	c.EventHandler.HandleSlashCommand(
		"/app",
		c.launchAppModal,
	)
	// Launch a request for App deployment
	c.EventHandler.HandleInteraction(
		slack.InteractionTypeViewSubmission,
		c.launchAppRequest,
	)

	// The App launch is approved
	c.EventHandler.HandleInteractionBlockAction(
		views.AppApprovedActionID,
		c.launchAppApproved,
	)

	// The App launch is rejected
	c.EventHandler.HandleInteractionBlockAction(
		views.AppRejectedActionID,
		c.launchAppRejected,
	)

	return c

}

// This method is used to launch a slack modal to get info about the deployment
func (c *AppLaunchController) launchAppModal(evt *socketmode.Event, clt *socketmode.Client) {
	// we need to cast our socketmode.Event
	command, _ := evt.Data.(slack.SlashCommand)

	// Make sure to respond to the server to avoid an error
	clt.Ack(*evt.Request)

	permissionRes := module.CheckPermission(command.ChannelID, command.UserID, module.Conf.PermittedChannel, module.Conf.Applicant)

	if permissionRes {
		// create the view using block-kit
		view := views.LaunchAppModal()

		// Open Modal (13)
		//_, err := clt.GetApiClient().OpenView(command.TriggerID, view)
		_, err := clt.OpenView(command.TriggerID, view)

		//Handle errors
		if err != nil {
			log.Printf("ERROR launch App Modal: %v", err)
		}
	} else {
		module.HandlePermissionDenied(command.UserName, command.ChannelID, command.UserID, clt)
	}
}

// This method is used to launch an approval request to Slack channel, then send a "waiting for response" message to the thread
func (c *AppLaunchController) launchAppRequest(evt *socketmode.Event, clt *socketmode.Client) {
	interaction := evt.Data.(slack.InteractionCallback)

	clt.Ack(*evt.Request)
	var approvor string
	if interaction.View.State.Values["game"]["game_id"].SelectedOption.Value == "Human Fall Flat" {
		approvor = "dingxing"
	} else if interaction.View.State.Values["game"]["game_id"].SelectedOption.Value == "Terraria" {
		approvor = "zhaoyongfeng"
	} else {
		approvor = ""
	}
	appArgs = &views.AppDeploymentInfo{
		User:     interaction.User.Name,
		Time:     interaction.View.State.Values["time"]["timepicker"].SelectedTime,
		Game:     interaction.View.State.Values["game"]["game_id"].SelectedOption.Value,
		IsOnline: interaction.View.State.Values["is_online"]["platform"].SelectedOption.Value,
		Date:     interaction.View.State.Values["date"]["datepicker"].SelectedDate,
		// Approvor: interaction.View.State.Values["mention"]["user_id"].SelectedUsers,
		Approvor: approvor,
		Comment:  interaction.View.State.Values["additional_comment"]["comment"].Value,
	}
	// channelID := interaction.ViewSubmissionCallback.ResponseURLs[len(interaction.ViewSubmissionCallback.ResponseURLs)-1].ChannelID
	// responseURL := interaction.ViewSubmissionCallback.ResponseURLs[len(interaction.ViewSubmissionCallback.ResponseURLs)-1].ResponseURL

	blocks := views.LaunchAppRequest(appArgs)

	//client := clt.GetApiClient()

	rC, rT, err := clt.PostMessage(
		interaction.ViewSubmissionCallback.ResponseURLs[len(interaction.ViewSubmissionCallback.ResponseURLs)-1].ChannelID,
		slack.MsgOptionBlocks(blocks...),
		// slack.MsgOptionResponseURL(responseURL, slack.ResponseTypeInChannel),
	)

	// Get the memory address of parent message ChannelID's value
	responseChannelID = &rC

	// Get the memory address of parent message TimeStamp's value
	responseTimeStamp = &rT

	// Handle errors
	if err != nil {
		log.Printf("ERROR while launch App Request: %v", err)
	}

	// Post the "Wating For Response" message to the thread
	approvalBlocks := views.LaunchAppApproval(appArgs)

	_, _, err0 := clt.PostMessage(
		*responseChannelID,
		// slack.MsgOptionPostMessageParameters(postMsgParams),
		slack.MsgOptionTS(*responseTimeStamp),
		slack.MsgOptionBlocks(approvalBlocks...),
	)

	// Handle errors
	if err != nil {
		log.Printf("ERROR while launch App Approval: %v", err0)
	}
}

// This method is used to post a message to thread after the request was rejected
func (c *AppLaunchController) launchAppRejected(evt *socketmode.Event, clt *socketmode.Client) {
	// cast our socketmode.Event into an App approved callback
	interaction := evt.Data.(slack.InteractionCallback)

	// Make sure to respond to the server to avoid an error
	clt.Ack(*evt.Request)

	permissionRes := module.CheckPermission(interaction.Container.ChannelID, interaction.User.ID, module.Conf.PermittedChannel, module.Conf.Approver)

	if permissionRes {
		blocks := views.LaunchAppRejected(appArgs.User, interaction.User.Name)

		_, _, err := clt.PostMessage(
			interaction.Container.ChannelID,
			slack.MsgOptionBlocks(blocks...),
			// slack.MsgOptionResponseURL(interaction.ResponseURL, slack.ResponseTypeInChannel),
			slack.MsgOptionReplaceOriginal(interaction.ResponseURL),
		)

		// Handle errors
		if err != nil {
			log.Printf("ERROR while sending message for /rocket: %v", err)
		}
	} else {
		module.HandlePermissionDenied(interaction.User.Name, interaction.Container.ChannelID, interaction.User.ID, clt)
	}
}

// This method is used to launch App server deployment after having the clearance
func (c *AppLaunchController) launchAppApproved(evt *socketmode.Event, clt *socketmode.Client) {
	// cast our socketmode.Event into an App approved callback
	interaction := evt.Data.(slack.InteractionCallback)

	// Make sure to respond to the server to avoid an error
	clt.Ack(*evt.Request)

	permissionRes := module.CheckPermission(interaction.Container.ChannelID, interaction.User.ID, module.Conf.PermittedChannel, module.Conf.Approver)

	if permissionRes {
		blocks := views.LaunchAppApproved(appArgs.User, interaction.User.Name)

		_, _, err := clt.PostMessage(
			interaction.Container.ChannelID,
			slack.MsgOptionBlocks(blocks...),
			// slack.MsgOptionResponseURL(interaction.ResponseURL, slack.ResponseTypeInChannel),
			slack.MsgOptionReplaceOriginal(interaction.ResponseURL),
		)

		// Handle errors
		if err != nil {
			log.Printf("ERROR while sending message \"approved\": %v", err)
		}

		// // execute the deployment
		// stdOut := module.ExecDeployment(app_args)

		// if stdOut == nil {
		// 	return
		// }

		// resultBlocks := views.LaunchAppResult(app_args.User, interaction.User.Name, stdOut)

		// _, _, err0 := client.PostMessage(
		// 	*responseChannelID,
		// 	// slack.MsgOptionPostMessageParameters(postMsgParams),
		// 	slack.MsgOptionTS(*responseTimeStamp),
		// 	slack.MsgOptionBlocks(resultBlocks...),
		// )

		// // Handle errors
		// if err0 != nil {
		// 	log.Printf("ERROR while post resultBlocks : %v", err0)
		// }

	} else {
		module.HandlePermissionDenied(interaction.User.Name, interaction.Container.ChannelID, interaction.User.ID, clt)
	}
}
