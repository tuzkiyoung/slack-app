package drivers

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"log"
	"os"
	"xd-infra-slack/config"
)

func InitSlackClient(c *config.Config) *socketmode.Client {
	api := slack.New(
		c.BotToken,
		slack.OptionDebug(c.Debug),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(c.AppToken),
	)

	return socketmode.New(
		api,
		socketmode.OptionDebug(c.Debug),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
}
