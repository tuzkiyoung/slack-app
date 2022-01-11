package module

import (
	"bytes"
	"log"
	"os/exec"
	"slack-app/views"
)

var started = make(chan struct{}, 1)

func ExecDeployment(appArgs *views.AppDeploymentInfo) []byte {
	if appArgs.IsOnline == "production" {
		// Execute the online deployment
		return LaunchAppDeploymentOnline(appArgs)
	} else {
		// Execute the dev deployment
		return LaunchAppDeploymentDev(appArgs)
	}
}

func LaunchAppDeploymentDev(app *views.AppDeploymentInfo) []byte {
	if len(started) == 1 {
		return nil
	}

	// dir := "/data/git.xindong.com/deploy-infra/ansible/HyperJoy-CN-T3Ansible/"
	// s := fmt.Sprintf("ansible-playbook -i inventory/P-HZ-v2-T3Game p-hz-v2-weekly_maintain.yml -D --extra-vars \"seconds_num=1 oss_accesskey=%v oss_accesstoken=%v matchserver_version=%v serverpkg_version=%v\"", oss_accesskey, oss_accesstoken, app_args.Approvor, app_args.Date)

	// bash command for test
	dir := "/root/"
	s := "date"

	started <- struct{}{}
	// time.Sleep(time.Second * 5)
	cmd := exec.Command("/bin/bash", "-c", s)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("ERROR execute ansible command: %v", err)
		return output
	}
	<-started
	n1 := bytes.Index(output, []byte("PLAY RECAP"))
	if n1 == -1 {
		n1 = 0
	}

	n2 := len(output) - 1
	output = output[n1:n2]

	return output
}

func LaunchAppDeploymentOnline(appArgs *views.AppDeploymentInfo) []byte {
	output := []byte{'o', 'n', 'l', 'i', 'n', 'e'}
	return output
}
