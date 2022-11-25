package config

import (
	"errors"
	"flag"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Slack `yaml:"slack"`
	Mysql `yaml:"mysql"`
	Arms  `yaml:"arms"`
}

type Mysql struct {
	DbPort int    `yaml:"dbPort"`
	DbUser string `yaml:"dbUser"`
	DbPwd  string `yaml:"dbPwd"`
	DbHost string `yaml:"dbHost"`
	DbName string `yaml:"dbName"`
}

type Slack struct {
	Debug    bool   `yaml:"debug"`
	AppToken string `yaml:"appToken"`
	BotToken string `yaml:"botToken"`
	ChanID   string `yaml:"chanID"`
}

type Arms struct {
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	Url             string `yaml:"url"`
}

func (c *Config) GetConfig() error {
	f := flag.String("c", "/config/config.yaml", "default config file path")
	flag.Parse()
	if f != nil {
		b, err := os.ReadFile(*f)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(b, c); err != nil {
			return err
		}
		if c.DbPort == 0 {
			c.DbPort = 3306
		}
		return nil
	}

	d, err := strconv.ParseBool(os.Getenv("debug"))
	if err != nil {
		return errors.New("invalid debug mode")
	}
	c.Debug = d

	c.ChanID = os.Getenv("chanID")
	if c.ChanID == "" {
		return errors.New("SLACK_Channel_ID must be set")
	}

	c.AppToken = os.Getenv("appToken")
	if c.AppToken == "" {
		return errors.New("SLACK_APP_TOKEN must be set")
	}
	if !strings.HasPrefix(c.AppToken, "xapp-") {
		return errors.New("SLACK_APP_TOKEN must have the prefix \"xapp-\"")
	}

	c.BotToken = os.Getenv("botToken")
	if c.BotToken == "" {
		return errors.New("SLACK_BOT_TOKEN must be set")
	}
	if !strings.HasPrefix(c.BotToken, "xoxb-") {
		return errors.New("SLACK_BOT_TOKEN must have the prefix \"xoxb-\"")
	}
	return nil
}
