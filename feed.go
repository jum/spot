/*
 * This is an unpublished work copyright 2015 Jens-Uwe Mager
 * 30177 Hannover, Germany, jum@anubis.han.de
 */
package spot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	SpotAPI = "https://api.findmespot.com/spot-main-web/consumer/rest-api/2.0/public/feed/%v/message.json?start=%v"
	DEBUG   = true
)

type Error struct {
	Code        string
	Text        string
	Description string
}

func (e Error) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Text)
}

type DebugError struct {
	error
	DebugBody bytes.Buffer
}

func (e DebugError) Error() string {
	return fmt.Sprintf("%s: %s", e.error, e.DebugBody.Bytes())
}

type MessageList []Message

func (m *MessageList) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errors.New("JSON body too short")
	}
	switch b[0] {
	case '[':
		var into []Message // NB: into is not of type MessageList, as that would result in infinite recursion
		err := json.Unmarshal(b, &into)
		*m = into
		return err
	case '{':
		var into Message
		err := json.Unmarshal(b, &into)
		*m = []Message{into}
		return err
	default:
		return errors.New("bad JSON for MessageList")
	}
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
				Message MessageList
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
		if resp.StatusCode/100 != 2 {
			err = fmt.Errorf("HTTP status: %s", resp.Status)
			return
		}
		var r io.Reader = resp.Body
		var de DebugError
		if DEBUG {
			r = io.TeeReader(r, &de.DebugBody)
		}
		d, err = DecodeFeed(r)
		if err != nil {
			if DEBUG {
				de.error = err
				err = de
			}
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
