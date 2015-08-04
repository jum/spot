/*
 * This is an unpublished work copyright 2015 Jens-Uwe Mager
 * 30177 Hannover, Germany, jum@anubis.han.de
 */
package spot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	SpotAPI = "https://api.findmespot.com/spot-main-web/consumer/rest-api/2.0/public/feed/%v/message.json?start=%v"
)

type Error struct {
	Code        string
	Text        string
	Description string
}

func (e Error) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Text)
}

type Feed struct {
	Response struct {
		Errors struct {
			Error Error
		}
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
		var d Feed
		var resp *http.Response
		resp, err = http.Get(fmt.Sprintf(SpotAPI, feedID, len(messages)))
		if err != nil {
			return
		}
		defer resp.Body.Close()
		d, err = DecodeFeed(resp.Body)
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

func DecodeFeed(r io.Reader) (f Feed, err error) {
	err = json.NewDecoder(r).Decode(&f)
	if err != nil {
		return
	}
	if len(f.Response.Errors.Error.Code) > 0 {
		err = f.Response.Errors.Error
		return
	}
	return
}
