package views

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"log"
)

func AppHomeTabView() slack.HomeTabViewRequest {
	str, err := slackAssets.ReadFile("assets/appHomeView.json")
	if err != nil {
		log.Printf("Unable to read view `AppHomeView`: %v", err)
	}
	view := slack.HomeTabViewRequest{}
	if err := json.Unmarshal(str, &view); err != nil {
		return slack.HomeTabViewRequest{}
	}
	return view
}

func AppServiceDetailModal(svcName string) slack.ModalViewRequest {
	asset := fmt.Sprintf("assets/serviceModal-%s.json", svcName)
	str, err := slackAssets.ReadFile(asset)
	if err != nil {
		log.Printf("Unable to read view `ModalHomeView`: %v", err)
	}
	view := slack.ModalViewRequest{}
	if err := json.Unmarshal(str, &view); err != nil {
		return slack.ModalViewRequest{}
	}
	return view
}
