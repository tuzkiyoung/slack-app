package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"slack-app/controllers"
	"slack-app/drivers"
	"slack-app/module"
)

func main() {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("Panic: %v", err)
		}
	}()
	module.Conf.GetConf()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// err := godotenv.Load("./t3_slack.env")
	// if err != nil {
	// 	log.Fatal().Msg("Error loading .env file")
	// }

	// Instanciate deps
	client, err := drivers.ConnectToSlackViaSocketmode(module.Conf.AppToken, module.Conf.BotToken)
	if err != nil {
		log.Error().
			Str("error", err.Error()).
			Msg("Unable to connect to slack")

		os.Exit(1)
	}

	// Inject Deps in router
	socketmodeHandler := drivers.NewsSocketmodeHandler(client)

	// This if for Separate articles and demos. You can run there separatly or all together
	// Build Slack Slash Command in Golang Using Socket Mode
	controllers.NewAppLaunchController(socketmodeHandler)

	socketmodeHandler.RunEventLoop()
}
