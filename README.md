# Modification Based on [Building Slack Bots in Golang](https://github.com/slack-go/slack/blob/master/examples/socketmode_handler/socketmode_handler.go)

This project demonstrates how to build a Slackbot in Golang; it uses the [slack-go](https://github.com/slack-go/slack) library and communicates with slack using the [socket mode](https://api.slack.com/apis/connections/socket).


## Run the project

Create a file `/config/config.yaml` with the following variables:

```
slack:
  appToken: "xapp-xxxxxx"
  botToken: "xoxb-xxxxxx"
  debug: false
  chanID: "xxxxxxx"
mysql:
  dbUser: xxxx
  dbPwd: "xxxxxxxxxxx"
  dbHost: xxxxxxxx
  dbPort: xxxx    // default 3306
  dbName: xxxxxxxx
arms:
  url: "https://alerts.aliyuncs.com/api/v1/integrations/custom/xxxxxx"
```

Run the application

```
go run main.go
```