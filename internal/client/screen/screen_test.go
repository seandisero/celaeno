package screen

import (
	"testing"
)

func TestFormatMessage(t *testing.T) {
	screen := &Screen{
		Width: 32,
	}
	t.Log("testing\n")
	message, err := screen.formatMessage("this is hopefully a long enough message for it to count up past 32 characters, we'll see how it goes.")
	if err != nil {
		t.Log("could not get message")
		t.FailNow()
	}
	for _, line := range message {
		t.Logf(" > %s\n", line)
	}
}
