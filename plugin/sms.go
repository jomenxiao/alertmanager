package plugin

import (
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/juju/errors"
	"net"
	"strings"
	"time"
)

func (r *Run) sendData(data string) error {
	sc := r.Config.Sms
	n, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", sc.SmsHost, sc.SmsPort), 10*time.Second)
	if err != nil {
		return errors.Trace(err)
	}
	defer n.Close()
	level.Debug(r.logger).Log("msg", "send message", data)
	if _, err := n.Write([]byte(data)); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (r *Run) gendata(msg string) (string, error) {
	/* data format
	   0189VZF01|1|77|sa|sa|15810270712|您data|1|0||system|
	data format */
	sc := r.Config.Sms
	data := []string{
		fmt.Sprintf("%s%s", sc.SmsPackfmt, sc.SmsTranscode),
		fmt.Sprintf("%s", sc.SmsRetval),
		fmt.Sprintf("%s", sc.SmsChanno),
		fmt.Sprintf("%s", sc.SmsUsername),
		fmt.Sprintf("%s", sc.SmsPassword),
		fmt.Sprintf("%s", sc.SmsMobile),
		fmt.Sprintf("%s", msg),
		fmt.Sprintf("%s", sc.SmsImmedflag),
		fmt.Sprintf("%s", sc.SmsSendtime),
		fmt.Sprintf("%s", sc.smsBranchno),
		fmt.Sprintf("%s", sc.SmsOperator),
	}

	dataString := strings.Join(data, "|")
	return fmt.Sprintf("%04d%s", len(dataString), dataString), nil
}

func (r *Run) handlerMsg(originMsg string) error {
	dataString, err := r.gendata(originMsg)
	if err != nil {
		return errors.Trace(err)
	}
	return r.sendData(dataString)

}

//TransferToSMS send alert's message to sms
func (r *Run) TransferToSMS(ad *AlertData) {
	for _, at := range ad.Alerts {
		omsg := fmt.Sprintf("name:%s url:%s", getValue(at.Labels, "alertname"), at.GeneratorURL)
		if len(omsg) > 1024 {
			omsg = omsg[:1024]
		}
		if err := r.handlerMsg(omsg); err != nil {
			level.Error(r.logger).Log("msg", "send alert omsg", omsg, err)
		}

	}
}
