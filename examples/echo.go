package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	fbmsg "github.com/stanaka/facebook-messenger"
)

var debug bool
var fb *fbmsg.FacebookMessenger

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Something wrong: %s\n", err.Error())
		return
	}
	if debug {
		log.Println("RecievedMessage Body:", string(b))
	}

	m, _ := url.ParseQuery(r.URL.RawQuery)
	fmt.Println(m["hub.verify_token"])
	if len(m["hub.verify_token"]) > 0 && m["hub.verify_token"][0] == os.Getenv("FB_VERIFY_TOKEN") && len(m["hub.challenge"]) > 0 {
		fmt.Fprintf(w, m["hub.challenge"][0])
		return
	}

	var msg fbmsg.CallbackMessage
	err = json.Unmarshal(b, &msg)
	if err != nil {
		fmt.Printf("Something wrong: %s\n", err.Error())
		return
	}

	for _, event := range msg.Entry[0].Messaging {
		sender := event.Sender.ID
		if event.Message != nil {
			fmt.Printf("Recieved Text: %s\n", event.Message.Text)
			err := fb.SendTextMessage(sender, event.Message.Text)
			if err != nil {
				fmt.Printf("Something wrong: %s\n", err.Error())
			}
		}
	}

}

func main() {
	debug = true
	fb = &fbmsg.FacebookMessenger{
		Token: os.Getenv("FB_TOKEN"),
	}

	http.HandleFunc("/fbbot/callback", callbackHandler)

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}
