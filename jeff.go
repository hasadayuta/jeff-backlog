package main

import (
	"fmt"
	"os"
	"strings"

	// https://[team-name].slack.com/apps/manage/custom-integrations -> Bots
	"github.com/nlopes/slack"
)

func main() {

	token := os.Getenv("SLACK_TOKEN")
	slackApi := slack.New(token)
	slackApi.SetDebug(true)

	rtm := slackApi.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")

			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)

				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					basicRespond(rtm, ev, prefix)
					backlogInfo(rtm, ev, prefix)
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				//Take no action
			}
		}
	}
}

func basicRespond(rtm *slack.RTM, msg *slack.MessageEvent, prefix string) {
	var response string
	text := msg.Text
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	acceptedGreetings := map[string]bool{
		"あなたは誰？": true,
		// "what's up?": true,
		// "yo":         true,
	}
	acceptedHowAreYou := map[string]bool{
		"調子はどう？": true,
		// "how are ya?":     true,
		// "feeling okay?":   true,
	}

	if acceptedGreetings[text] {
		response = "Back Logのお手伝いするジェフ・ベゾスだ"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	} else if acceptedHowAreYou[text] {
		response = "とても良いです！ありがとう！"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	}
}

func backlogInfo(rtm *slack.RTM, msg *slack.MessageEvent, prefix string) {
	var response string
	text := msg.Text
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	// backlog_api_key := os.Getenv("BACKLOG_API_KEY")
	// resource := "space"
	// backlog_space_url := "https://tenso.backlog.jp/api/v2/" + resource + "?apiKey=" + backlog_api_key

	basicInfomation := map[string]bool{
		"基本情報": true,
	}

	if basicInfomation[text] {
		response = "Backlog APIのレスポンスをparseして返します"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	}
}
