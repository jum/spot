/*
 * This is an unpublished work copyright 2015 Jens-Uwe Mager
 * 30177 Hannover, Germany, jum@anubis.han.de
 */
package spot

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	SpotAPI = "https://api.findmespot.com/spot-main-web/consumer/rest-api/2.0/public/feed/%v/message.json?start=%v"
)

type SpotFeed struct {
	Response struct {
		FeedMessageResponse struct {
			Count int
			Feed  struct {
				Id                   string
				Name                 string
				Description          string
				Status               string
				Usage                int
				DaysRange            int
				DetailedMessageShown bool
			}
			TotalCount   int
			AcivityCount int
			Messages     struct {
				Message []Message
			}
		}
	}
}

func RetrieveMessages(feedID string) (messages []Message, err error) {
	for {
		var d SpotFeed
		var resp *http.Response
		resp, err = http.Get(fmt.Sprintf(SpotAPI, feedID, len(messages)))
		if err != nil {
			return
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&d)
		if err != nil {
			return
		}
		messages = append(messages, d.Response.FeedMessageResponse.Messages.Message...)
		//fmt.Printf("count %v, total %v\n", len(messages), d.Response.FeedMessageResponse.TotalCount)
		if len(messages) >= d.Response.FeedMessageResponse.TotalCount {
			break
		}
	}
	return
}
