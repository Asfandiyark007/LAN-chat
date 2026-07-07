package internal

import (
	"encoding/json"
	"fmt"
	"testing"

	"lan-chat/protocol"
)

type WireMessages struct {
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	Timestamp string `json:"timestamp"`
	Content   string `json:"content"`
}

func TestWireMessage(t *testing.T) {
	message := protocol.WireMessage{
		Type:   "USER",
		Sender: "Alice",
		// Timestamp: "2026-07-07T12:00:00Z",
		Content: "Hello|everyone",
	}

	data, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	var decoded protocol.WireMessage
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(data))

	fmt.Println("Type:", decoded.Type)
	fmt.Println("Sender:", decoded.Sender)
	fmt.Println("Timestamp:", decoded.Timestamp)
	fmt.Println("Content:", decoded.Content)
}
