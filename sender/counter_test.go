package sender_test

import (
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/parser"
	"github.com/KristianLyng/skogul/sender"
	"testing"
	"time"
)

func TestCounter(t *testing.T) {
	one := &(sender.Test{})
	two := &(sender.Test{})

	h := skogul.Handler{Sender: one, Parser: parser.JSON{}}
	cnt := sender.Counter{Next: two, Stats: h, Period: time.Duration(50 * time.Millisecond)}
	two.TestQuick(t, &cnt, &validContainer, 1)
	if one.Received != 0 {
		t.Errorf("Stats received too early.")
	}
	two.TestQuick(t, &cnt, &validContainer, 1)
	two.TestQuick(t, &cnt, &validContainer, 1)
	two.TestQuick(t, &cnt, &validContainer, 1)
	if one.Received != 0 {
		t.Errorf("Stats received too early 2.")
	}
	time.Sleep(time.Duration(50 * time.Millisecond))
	if one.Received != 0 {
		t.Errorf("Stats received too early 2.")
	}
	two.TestQuick(t, &cnt, &validContainer, 1)
	if one.Received != 1 {
		t.Errorf("Correct stats not received")
	}
	two.TestQuick(t, &cnt, &validContainer, 1)
	two.TestQuick(t, &cnt, &validContainer, 1)
	time.Sleep(time.Duration(50 * time.Millisecond))
	two.TestQuick(t, &cnt, &validContainer, 1)
	if one.Received != 2 {
		t.Errorf("Correct stats not received")
	}
}
