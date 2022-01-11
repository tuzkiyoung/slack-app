package views

import (
	"embed"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/slack-go/slack"
)

const (
	// AppApprovedActionID Define Action_id as constant, so we can refet to them in the controller
	AppApprovedActionID = "app_launch_approved"
	AppRejectedActionID = "app_launch_rejected"
)

// AppDeploymentInfo define a App deployment struct
type AppDeploymentInfo struct {
	IsOnline string
	User     string
	Time     string
	Game     string
	Date     string
	Comment  string
	Approvor string
}

//go:embed slackCommandAssets/*
var appLaunchAssets embed.FS

func LaunchAppModal() slack.ModalViewRequest {
	str, err := appLaunchAssets.ReadFile("slackCommandAssets/modal.json")
	if err != nil {
		log.Printf("Unable to read view `LaunchModal`: %v", err)
	}
	view := slack.ModalViewRequest{}
	_ = json.Unmarshal(str, &view)

	return view

}

func HandlePermissionDenied(name string) []slack.Block {
	// we need a stuct to hold template arguments
	type args struct {
		User string
	}

	myArgs := args{
		User: name,
	}
	tpl := renderTemplate(appLaunchAssets, "slackCommandAssets/permissiondenied.json", myArgs)

	// we convert the view into a message struct
	view := slack.Msg{}

	str, _ := ioutil.ReadAll(&tpl)
	_ = json.Unmarshal(str, &view)

	// We only return the block because of the way the PostEphemeral function works
	// we are going to use Slack.MsgOptionBlocks in the controller
	return view.Blocks.BlockSet
}

func LaunchAppRequest(appArgs *AppDeploymentInfo) []slack.Block {

	tpl := renderTemplate(appLaunchAssets, "slackCommandAssets/request.json", appArgs)

	// we convert the view into a message struct
	view := slack.Msg{}

	str, _ := ioutil.ReadAll(&tpl)
	_ = json.Unmarshal(str, &view)

	// We only return the block because of the way the PostEphemeral function works
	// we are going to use Slack.MsgOptionBlocks in the controller
	return view.Blocks.BlockSet
}

func LaunchAppApproval(appArgs *AppDeploymentInfo) []slack.Block {

	tpl := renderTemplate(appLaunchAssets, "slackCommandAssets/approval.json", appArgs)

	// we convert the view into a message struct
	view := slack.Msg{}

	str, _ := ioutil.ReadAll(&tpl)
	_ = json.Unmarshal(str, &view)

	// We only return the block because of the way the PostEphemeral function works
	// we are going to use Slack.MsgOptionBlocks in the controller
	return view.Blocks.BlockSet
}

func LaunchAppRejected(user, approvor string) []slack.Block {
	// we need a stuct to hold template arguments
	type args struct {
		User     string
		Approvor string
	}

	myArgs := args{
		User:     user,
		Approvor: approvor,
	}
	tpl := renderTemplate(appLaunchAssets, "slackCommandAssets/rejection.json", myArgs)

	// we convert the view into a message struct
	view := slack.Msg{}

	str, _ := ioutil.ReadAll(&tpl)
	_ = json.Unmarshal(str, &view)

	// We only return the block because of the way the PostEphemeral function works
	// we are going to use Slack.MsgOptionBlocks in the controller
	return view.Blocks.BlockSet
}

func LaunchAppApproved(user, approvor string) []slack.Block {
	// we need a stuct to hold template arguments
	type args struct {
		User     string
		Approvor string
	}

	myArgs := args{
		User:     user,
		Approvor: approvor,
	}
	tpl := renderTemplate(appLaunchAssets, "slackCommandAssets/deployment.json", myArgs)

	// we convert the view into a message struct
	view := slack.Msg{}

	str, _ := ioutil.ReadAll(&tpl)
	_ = json.Unmarshal(str, &view)

	// We only return the block because of the way the PostEphemeral function works
	// we are going to use Slack.MsgOptionBlocks in the controller
	return view.Blocks.BlockSet
}

func LaunchAppResult(user, approvor string, stdOut []byte) []slack.Block {
	stdOutStr := strings.ReplaceAll(string(stdOut), "\n", "\\n")
	stdOutStr = strings.ReplaceAll(stdOutStr, "\t", "  ")
	stdOutStr = strings.ReplaceAll(stdOutStr, "*", "")

	// we need a stuct to hold template arguments
	type args struct {
		User     string
		Approvor string
		Res      string
	}

	myArgs := args{
		User:     user,
		Approvor: approvor,
		Res:      stdOutStr,
	}

	tpl := renderTemplate(appLaunchAssets, "slackCommandAssets/result.json", myArgs)

	// we convert the view into a message struct
	view := slack.Msg{}

	str, _ := ioutil.ReadAll(&tpl)
	_ = json.Unmarshal(str, &view)

	// We only return the block because of the way the PostEphemeral function works
	// we are going to use Slack.MsgOptionBlocks in the controller
	return view.Blocks.BlockSet
}
