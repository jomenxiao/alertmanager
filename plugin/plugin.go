package plugin

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/unrolled/render"
	"net/http"
	"time"
)

//Run runtime
type Run struct {
	Rdr         *render.Render
	AlertMsgs   chan *AlertData
	KafkaClient sarama.SyncProducer
	Config      PluginConfig
	logger      log.Logger
}

//TransferData transfer alert's message
func (r *Run) TransferData(ad *AlertData) {
	if r.Config.Kafka.KafkaAddress != "" {
		r.TransferToKafka(ad)
	}
	if r.Config.Sms.SmsHost != "" {
		r.TransferToSMS(ad)
	}
}

//Scheduler scheduler
func (r *Run) Scheduler() {
	for {
		lenAlertMsgs := len(r.AlertMsgs)
		if lenAlertMsgs > 0 {
			for i := 0; i < lenAlertMsgs; i++ {
				go r.TransferData(<-r.AlertMsgs)
			}
		}
		time.Sleep(3 * time.Second)
	}
}

//RunPlugin run plugin
func RunPlugin() {
	r := &Run{
		AlertMsgs: make(chan *AlertData, 10000),
	}

	if err := r.ReadConfig(); err != nil {
		level.Error(r.logger).Log(fmt.Printf("config error %v", err))
		return
	}

	if r.Config.Kafka.KafkaAddress != "" {
		if err := r.CreateKafkaProduce(); err != nil {
			log.Errorf("create kafka produce error %v", err)
			return
		}
	}

	go func() {
		log.Infof("create http server")
		r.CreateRender()
		http.ListenAndServe(fmt.Sprintf(":%d", r.Config.ListenPort), r.CreateRouter())
	}()

	go r.Scheduler()
}
