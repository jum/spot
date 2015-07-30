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
