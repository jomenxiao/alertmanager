package plugin

import (
	"encoding/json"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"io/ioutil"
	"net/http"
	"time"
)

// AlertData alertmanager base struct
type AlertData struct {
	Receiver          string `json:"receiver"`
	Status            string `json:"status"`
	Alerts            Alerts `json:"alerts"`
	GroupLabels       KV     `json:"groupLabels"`
	CommonLabels      KV     `json:"commonLabels"`
	CommonAnnotations KV     `json:"commonAnnotations"`
	ExternalURL       string `json:"externalURL"`
}

// Alert holds one alert for notification templates.
type Alert struct {
	Status       string    `json:"status"`
	Labels       KV        `json:"labels"`
	Annotations  KV        `json:"annotations"`
	StartsAt     time.Time `json:"startsAt"`
	EndsAt       time.Time `json:"endsAt"`
	GeneratorURL string    `json:"generatorURL"`
}

// Alerts is a list of Alert objects.
type Alerts []Alert

// KV is a set of key/value string pairs.
type KV map[string]string

//AlertMsgFromWebhook get alert massage
func (r *Run) AlertMsgFromWebhook(w http.ResponseWriter, hr *http.Request) {
	b, err := ioutil.ReadAll(hr.Body)
	if err != nil {
		level.Error(r.logger).Log("msg", "read http post data", err)
		r.Rdr.Text(w, http.StatusBadRequest, err.Error())
		return
	}
	level.Debug(r.logger).Log("msg", "alert data", string(b))
	defer hr.Body.Close()
	alertData := &AlertData{}
	err = json.Unmarshal(b, alertData)
	if err != nil {
		level.Error(r.logger).Log("msg", "unmarshal http data", err)
		r.Rdr.Text(w, http.StatusBadRequest, err.Error())
		return
	}
	level.Debug(r.logger).Log("msg", "alert data", alertData)
	r.AlertMsgs <- alertData
	r.Rdr.Text(w, http.StatusAccepted, "")
}

//CreateRouter create router
func (r *Run) CreateRouter() *mux.Router {
	m := mux.NewRouter()
	m.HandleFunc(r.Config.WebhookPath, r.AlertMsgFromWebhook).Methods("POST")
	return m
}

// CreateRender for render.
func (r *Run) CreateRender() {
	r.Rdr = render.New(render.Options{
		IndentJSON: true,
	})
}
