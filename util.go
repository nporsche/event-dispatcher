package dispatcher

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	defaultChannel  = "default_ch"
	updateTopicsDur = time.Minute * 2
	topicPrefix     = "x7x6y9-"
)

var idx uint16 = 0

type ListTopicResponse struct {
	StatusCode int    `json:"status_code"`
	StatusTxt  string `json:"status_txt"`
	Data       Topics `json:"data"`
}

type Topics struct {
	Topics []string `json:"topics"`
}

type TopicNameBuilder func(e EventHeader) string

func DefaultTopicNameBuilder(e EventHeader) string {
	return e.String()
}

func listTopics(lookupdEndpoints []string) (topics []string, err error) {
	//{"status_code":200,"status_txt":"OK","data":{"topics":[]}}
	if len(lookupdEndpoints) == 0 {
		return nil, errors.New("No lookupd endpoints")
	}
	endpoint := lookupdEndpoints[idx%uint16(len(lookupdEndpoints))]
	resp, err := http.Get(fmt.Sprintf("http://%s/topics", endpoint))
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var topicsResponse ListTopicResponse
	err = json.Unmarshal(b, &topicsResponse)
	if err != nil {
		fmt.Println(string(b))
		return nil, err
	}

	return topicsResponse.Data.Topics, nil
}

func EventHeaderToString(e *EventHeader) string {
	var buffer bytes.Buffer
	buffer.WriteString(topicPrefix)
	buffer.WriteString(fmt.Sprintf("%s-%s-%s-%s", e.King, e.Noble, e.Knight, e.Peasant))
	for _, tag := range e.Tags {
		buffer.WriteString(fmt.Sprintf("-%s", tag))
	}

	return buffer.String()
}

func StringToEventHeader(s string) (*EventHeader, error) {
	if !strings.HasPrefix(s, topicPrefix) {
		return nil, errors.New("skip this topic")
	}
	s = s[len(topicPrefix):]
	strs := strings.Split(s, "-")
	if len(strs) < 4 {
		return nil, errors.New("skip this topic because - is less than 4")
	}
	eh := &EventHeader{
		King:    strs[0],
		Noble:   strs[1],
		Knight:  strs[2],
		Peasant: strs[3],
		Tags:    strs[4:],
	}
	return eh, nil
}
