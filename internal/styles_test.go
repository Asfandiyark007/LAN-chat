package internal

import (
	"fmt"
	"testing"
)

func TestStyles(t *testing.T) {
	fmt.Println(OwnMessageStyle.Render("Me: Hello"))
	fmt.Println(OtherMessageStyle.Render("Alice: Hi"))
	fmt.Println(SystemMessageStyle.Render("Bob joined the chat"))
}
