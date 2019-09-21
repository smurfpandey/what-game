package main

import (
	"fmt"
	"github.com/imroc/req"
)

const (
	API_URL = "https://now.smurfpandey.me/playing"
)

func NotifyGameStarted(startedGame Game) {
	rawResp, err := req.Post(API_URL, req.BodyJSON(&startedGame))

	if err != nil {
		fmt.Println("Error making POST request")
		return
	}
	resp := rawResp.Response()

	if resp.StatusCode != 200 {
		fmt.Println("API responded with not OK code", resp.StatusCode)
		return
	}

	fmt.Println("Notified via API: Started")
}

func NotifyGameExited() {
	rawResp, err := req.Delete(API_URL)

	if err != nil {
		fmt.Println("Error making DELETE request")
		return
	}
	resp := rawResp.Response()

	if resp.StatusCode != 200 {
		fmt.Println("API responded with not OK code", resp.StatusCode)
		return
	}

	fmt.Println("Notified via API: Exited")
}
