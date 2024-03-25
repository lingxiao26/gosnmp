package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gosnmp/config"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Message struct {
	MsgType string `json:"msgtype"`
	Text    *Text  `json:"text"`
}

type Text struct {
	Content             string   `json:"content"`
	MentionedMobileList []string `json:"mentioned_mobile_list"`
}

func NewMessage() *Message {
	return &Message{
		MsgType: "text",
		Text:    &Text{},
	}
}

func (msg *Message) Alert(ch *config.Host, body string) {
	msg.Text.Content = fmt.Sprintf("服务器基础资源告警\n\n服务器: %s\n告警内容: %s\n", ch.Addr, body)
	msg.Text.MentionedMobileList = ch.At

	// 序列化 请求body
	marshalMsg, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("marshal sendMsg: %v", err)
	}
	brd := bytes.NewReader(marshalMsg)

	// 发送请求
	resp, err := http.Post(ch.Webhook, "application/json", brd)
	if err != nil {
		log.Errorf("http.Post(): %v", err)
	}
	defer resp.Body.Close()

	// 处理响应消息
	var respData = make([]byte, 1024)
	_, err = resp.Body.Read(respData)
	if err != nil {
		log.Errorf("response body read: %v", err)
	}

	log.Infof("qiwei response: %s", respData)
}
