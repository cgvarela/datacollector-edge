package dpm

import (
	"bytes"
	"encoding/json"
	"github.com/streamsets/dataextractor/container/common"
	"io/ioutil"
	"log"
	"net/http"
)

type ClientEvent struct {
	EventId      string   `json:"eventId"`
	Destinations []string `json:"destinations"`
	RequiresAck  bool     `json:"requiresAck"`
	IsAckEvent   bool     `json:"ackEvent"`
	EventTypeId  int      `json:"eventTypeId"`
	Payload      string   `json:"payload"`
	OrgId        string   `json:"orgId"`
}

type ServerEvent struct {
	EventId      string `json:"eventId"`
	From         string `json:"from"`
	RequiresAck  bool   `json:"requiresAck"`
	IsAckEvent   bool   `json:"isAckEvent"`
	EventTypeId  int    `json:"eventTypeId"`
	Payload      string `json:"payload"`
	ReceivedTime int64  `json:"receivedTime"`
	OrgId        string `json:"orgId"`
}

type SDCInfoEvent struct {
	SdeId           string            `json:"sdcId"`
	HttpUrl         string            `json:"httpUrl"`
	Sdc2GoBuildInfo *common.BuildInfo `json:"sdc2goBuildInfo"`
	Labels          []string          `json:"labels"`
}

func SendEvent(dpmConfig Config, buildInfo *common.BuildInfo, runtimeInfo *common.RuntimeInfo) {
	if dpmConfig.Enabled && dpmConfig.AppAuthToken != "" {

		sdcInfoEvent := SDCInfoEvent{
			SdeId:           runtimeInfo.ID,
			HttpUrl:         runtimeInfo.HttpUrl,
			Sdc2GoBuildInfo: buildInfo,
			Labels:          dpmConfig.JobLabels,
		}

		sdcInfoEventJson, _ := json.Marshal(sdcInfoEvent)

		clientEvent := ClientEvent{
			EventId:      runtimeInfo.HttpUrl,
			Destinations: []string{dpmConfig.EventsRecipient},
			RequiresAck:  false,
			IsAckEvent:   false,
			EventTypeId:  2001,
			Payload:      string(sdcInfoEventJson),
			OrgId:        "",
		}

		log.Println("Client Event JSON:")
		jsonValue, err := json.Marshal([]ClientEvent{clientEvent})
		if err != nil {
			log.Println(err)
		}

		log.Println(string(jsonValue))

		var eventsUrl = dpmConfig.BaseUrl + "/messaging/rest/v1/events"

		req, err := http.NewRequest("POST", eventsUrl, bytes.NewBuffer(jsonValue))
		req.Header.Set("X-SS-App-Auth-Token", dpmConfig.AppAuthToken)
		req.Header.Set("X-SS-App-Component-Id", runtimeInfo.ID)
		req.Header.Set("X-Requested-By", "SDC2GO")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		log.Println("DPM Event Status:", resp.Status)
		if resp.StatusCode != 200 {
			panic("DPM Send event failed")
		}
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("response Body:", string(body))
		runtimeInfo.DPMEnabled = true
	} else {
		runtimeInfo.DPMEnabled = false
	}
}
