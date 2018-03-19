package plugin

import (
	"github.com/juju/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const (
	//ConfigFileName is hardcode
	ConfigFileName = "alertmanger_plugin.yaml"
)

//KafkaConfig kafka configure
type KafkaConfig struct {
	KafkaAddress string `yaml:"kafka_address"`
	KafkaTopic   string `yaml:"kafka_topic"`
}

//SmsConfig sms config
type SmsConfig struct {
	SmsHost      string `yaml:"sms_host"`
	SmsPort      string `yaml:"sms_port"`
	smsBranchno  string `yaml:"sms_branchno"`
	SmsChanno    string `yaml:"sms_channo"`
	SmsImmedflag string `yaml:"sms_immedflag"`
	SmsMobile    string `yaml:"sms_mobile"`
	SmsOperator  string `yaml:"sms_operator"`
	SmsUsername  string `yaml:"sms_username"`
	SmsPassword  string `yaml:"sms_passsword"`
	SmsRetval    string `yaml:"sms_retval"`
	SmsTranscode string `yaml:"sms_transcode"`
	SmsPackfmt   string `yaml:"sms_packfmt"`
	SmsSendtime  string `yaml:"sms_sendtime"`
	SmsMessage   string `yaml:"sms_message"`
}

//PluginConfig plugin config file
type PluginConfig struct {
	ListenPort  int         `yaml:"listen_port"`
	WebhookPath string      `yaml:"webhook_path"`
	Kafka       KafkaConfig `yaml:"kafka_config"`
	Sms         SmsConfig   `yaml:"sms_config"`
}

//ReadConfig read config file
func (r *Run) ReadConfig() error {
	if _, err := os.Stat(ConfigFileName); os.IsNotExist(err) {
		return err
	}
	c, err := ioutil.ReadFile(ConfigFileName)
	if err != nil {
		return errors.Trace(err)
	}
	return yaml.Unmarshal(c, &r.Config)
}
