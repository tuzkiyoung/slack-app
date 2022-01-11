# Modification Based on [Building Slack Bots in Golang](https://github.com/xNok/slack-go-demo-socketmode)

This project demonstrates how to build a Slackbot in Golang; it uses the [slack-go](https://github.com/slack-go/slack) library and communicates with slack using the [socket mode](https://api.slack.com/apis/connections/socket).


## Test the project

Create a file `/opt/configs/config.yaml` with the following variables:

```
SLACK_BOT_TOKEN=xoxb-xxxxxxxxxxx
SLACK_APP_TOKEN=xapp-1-xxxxxxxxx

accesskey:
accesstoken:

channel:
 - {Channel ID}
applicant:
 - {User ID}
approver:
 - {User ID}
```

Run the application

```
go run main.go
```