package plugin

import (
	//"fmt"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/go-kit/kit/log/level"
	"strings"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

//KafkaMsg represent kafka msg
type KafkaMsg struct {
	Title       string `json:"title"`
	Source      string `json:"source"`
	Node        string `json:"node"`
	Expr        string `json:"expr"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Level       string `json:"level"`
	Note        string `json:"note"`
	Value       string `json:"value"`
	Time        string `json:"time"`
}

//CreateKafkaProduce create a kafka produce
func (r *Run) CreateKafkaProduce() error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewManualPartitioner
	var err error
	r.KafkaClient, err = sarama.NewSyncProducer(strings.Split(r.Config.Kafka.KafkaAddress, ","), config)
	return err
}

//PushKafkaMsg push message to kafka cluster
func (r *Run) PushKafkaMsg(msg string) error {
	kafkaMsg := &sarama.ProducerMessage{
		Topic: r.Config.Kafka.KafkaTopic,
		Value: sarama.StringEncoder(msg),
	}
	level.Debug(r.logger).Log("msg", "kafka message", msg)
	_, _, err := r.KafkaClient.SendMessage(kafkaMsg)
	return err
}

//TransferToKafka transfer alert to kafka string
func (r *Run) TransferToKafka(ad *AlertData) {
	for _, at := range ad.Alerts {
		kafkaMsg := &KafkaMsg{
			Title:       getValue(at.Labels, "alertname"),
			Description: getValue(at.Annotations, "description"),
			Expr:        getValue(at.Labels, "expr"),
			Level:       getValue(at.Labels, "level"),
			Node:        getValue(at.Labels, "instance"),
			Source:      getValue(at.Labels, "env"),
			Value:       getValue(at.Annotations, "value"),
			Note:        getValue(at.Annotations, "summary"),
			URL:         at.GeneratorURL,
			Time:        at.StartsAt.Format(timeFormat),
		}
		atByte, err := json.Marshal(kafkaMsg)
		if err != nil {
			level.Error(r.logger).Log("msg", "marshal data", err)
			continue
		}

		if err := r.PushKafkaMsg(string(atByte)); err != nil {
			level.Error(r.logger).Log("msg", "push message to kafka", err)
		}
	}
}
