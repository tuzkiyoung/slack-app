package module

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"log"
	"slack-app/views"
)

func CheckPermission(channelID, userID string, permittedChannel, permittedUser []string) bool {
	for _, c := range permittedChannel {
		if channelID == c {
			for _, u := range permittedUser {
				if userID == u {
					return true
				}
			}
		}
	}
	return false
}

func HandlePermissionDenied(name string, channelID string, userID string, clt *socketmode.Client) {
	// create the view using block-kit
	deniedBlocks := views.HandlePermissionDenied(name)

	// Post greeting message (3)
	// We get the Api client from `clt`
	_, err := clt.PostEphemeral(
		channelID,
		userID,
		slack.MsgOptionBlocks(deniedBlocks...),
	)
	//Handle errors
	if err != nil {
		log.Printf("ERROR handle permission denied: %v", err)
	}
}
