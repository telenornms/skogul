package sender_test

import (
	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/sender"
	"testing"
	"time"
)

// BackTester will fail until it has failed fails times.
type BackTester struct {
	fails int
}

func (bt *BackTester) Send(c *skogul.Container) error {
	if bt.fails > 0 {
		bt.fails--
		return skogul.Error{Source: "back tester", Reason: "still failing"}
	}
	return nil
}

// TestBackoff tests if backoff works at least a little bit
func TestBackoff(t *testing.T) {
	te := BackTester{fails: 1}
	bo := sender.Backoff{Next: skogul.SenderRef{S: &te},
		Base:    skogul.Duration{Duration: time.Duration(time.Millisecond * 10)},
		Retries: 2}
	err := bo.Send(&validContainer)
	if err != nil {
		t.Errorf("Got error from bo.Send(): %v", err)
	}
	te.fails = 10
	err = bo.Send(&validContainer)
	if err == nil {
		t.Errorf("Didn't get error from bo.Send()")
	}
}
