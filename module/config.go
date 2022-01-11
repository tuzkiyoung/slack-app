package module

import (
	"flag"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

// Conf Declear a permission config variable
var Conf Config

var confFlag = flag.String(
	"confPath",
	"/opt/configs/config.yaml",
	"optional path of config file",
)

type Config struct {
	Accesskey        string   `yaml:"accesskey"`
	Accesstoken      string   `yaml:"accesstoken"`
	AppToken         string   `yaml:"apptoken"`
	BotToken         string   `yaml:"bottoken"`
	Approver         []string `yaml:"approver"`
	Applicant        []string `yaml:"applicant"`
	PermittedChannel []string `yaml:"channel"`
}

func (c *Config) GetConf() *Config {
	flag.Parse()
	yamlFile, err := ioutil.ReadFile(*confFlag)
	if err != nil {
		log.Printf("Unable to read yaml file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, c)

	if err != nil {
		log.Printf("Unable to Unmarshal yaml file: %v", err)
	}

	return c
}

// func getKubernetesConfig(path string) *string {
// 	var kubeconfig *string
// 	home := homedir.HomeDir()
// 	if home != "" {
// 		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
// 	} else {
// 		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
// 	}
// 	flag.Parse()
// 	return *kubeconfig
// }
