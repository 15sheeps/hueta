package main

import (
	"strings"
	"fmt"
	"log"
	"gopkg.in/yaml.v3"
	"net/http"
	"github.com/gorilla/websocket"
)

func (usr *Instance) DecideNextWsMessage(msg string, c chan<- string) {
	// у челов криво все я ебал
	msg = strings.Replace(msg, "alternateIPs: ,\n", "", -1)

	var data struct {
		Type         string      `yaml:"type"`
		MatchID      string      `yaml:"matchId"`
		IP           string      `yaml:"ip"`
		Port         int         `yaml:"port"`
		CustomAttribute string	 `yaml:"customAttribute"`
	}

	err := yaml.Unmarshal([]byte(msg), &data)
	if err != nil {
		log.Println(err)
	}

	switch data.Type {
	case "connectNotif":
		c <- "type: partyCreateRequest\nid: my_unique_request_id"
	case "partyCreateResponse": 
		c <- "type: partyInfoRequest\nid: my_unique_request_id"
	case "partyInfoResponse":
		c <- "startMatchmakingRequest\nid: my_unique_request_id\ngameMode: recruitarena\nlatencies: " + `{"us-east-2":991,"me-central-1":991,"eu-west-2":1,"ap-southeast-1":991,"sa-east-1":991,"us-west-2":991}`
	case "matchmakingNotif":
		c <- "type: setReadyConsentRequest\nid: my_unique_request_id\nmatchId: " + data.MatchID
	case "dsNotif":
		grpcAddr := fmt.Sprintf("%s:%d", data.IP, data.Port)

		usr.Play(grpcAddr, data.CustomAttribute, true)
	}
}

func (usr *Instance) StartMatchmaking() {
    wsUrl := "wss://api.voxies.io/lobby/"

    header := make(http.Header)
    header.Set("Authorization", "Bearer " + usr.token)

    c, _, err := websocket.DefaultDialer.Dial(wsUrl, header)
    if err != nil {
    	log.Println("failed to dial:", err)
        return
    }
    defer c.Close()

    messageOut := make(chan string)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			log.Println("RECEIVED:", string(message))
			fmt.Println()
			usr.DecideNextWsMessage(string(message), messageOut)
		}
	}()

	for {
		select {
		case <-done:
			return
		case m := <-messageOut:
			err := c.WriteMessage(websocket.TextMessage, []byte(m))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			log.Println("SEND:", m)
			fmt.Println()
		}
	}
}