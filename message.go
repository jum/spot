/*
 * This is an unpublished work copyright 2015 Jens-Uwe Mager
 * 30177 Hannover, Germany, jum@anubis.han.de
 */
package spot

type Message struct {
	ClientUnixTime int64
	Id             int64
	MessengerId    string
	MessengerName  string
	UnixTime       int64
	MessageType    string
	Latitude       float64
	Longitude      float64
	ModelId        string
	ShowCustomMsg  string
	DateTime       string
	BatteryState   string
	Hidden         int
	MessageContent string
}

func MergeMessages(db []Message, n []Message) []Message {
	for _, e := range n {
		found := false
		for _, old := range db {
			if old.Id == e.Id {
				found = true
				break
			}
		}
		if !found {
			db = append(db, e)
		}
	}
	return db
}

type MessageTimeSorter []Message

func (p MessageTimeSorter) Len() int           { return len(p) }
func (p MessageTimeSorter) Less(i, j int) bool { return p[i].UnixTime < p[j].UnixTime }
func (p MessageTimeSorter) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
