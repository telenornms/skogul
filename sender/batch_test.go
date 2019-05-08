package sender_test

import (
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/sender"
	"testing"
	"time"
)

func TestBatch(t *testing.T) {
	c := skogul.Container{}
	m := skogul.Metric{}
	n := time.Now()
	m.Time = &n
	m.Metadata = make(map[string]interface{})
	m.Data = make(map[string]interface{})
	m.Data["test"] = 5
	m.Metadata["key"] = "value"

	c.Metrics = []*skogul.Metric{&m}

	one := &(testSender{})

	batch := sender.Batch{Next: one}

	for i := 0; i < 9; i++ {
		err := batch.Send(&c)
		if err != nil {
			t.Errorf("batch.Send() failed: %v", err)
		}
		if one.received != 0 {
			t.Errorf("batch.Send(), sender 1 expected %d recevied, got %d", 0, one.received)
		}
	}
	err := batch.Send(&c)
	if err != nil {
		t.Errorf("batch.Send() failed: %v", err)
	}
	time.Sleep(time.Duration(100 * time.Millisecond))
	if one.received != 1 {
		t.Errorf("batch.Send(), sender 1 expected %d recevied, got %d", 1, one.received)
	}

	for i := 0; i < 9; i++ {
		err := batch.Send(&c)
		if err != nil {
			t.Errorf("batch.Send() failed: %v", err)
		}
		if one.received != 1 {
			t.Errorf("batch.Send(), sender 1 expected %d recevied, got %d", 1, one.received)
		}
	}
	err = batch.Send(&c)
	if err != nil {
		t.Errorf("batch.Send() failed: %v", err)
	}
	time.Sleep(time.Duration(100 * time.Millisecond))
	if one.received != 2 {
		t.Errorf("batch.Send(), sender 1 expected %d recevied, got %d", 2, one.received)
	}
	err = batch.Send(&c)
	if err != nil {
		t.Errorf("batch.Send() failed: %v", err)
	}
	time.Sleep(time.Duration(100 * time.Millisecond))
	if one.received != 2 {
		t.Errorf("batch.Send(), sender 1 expected %d recevied, got %d", 2, one.received)
	}
	time.Sleep(time.Duration(1 * time.Second))
	if one.received != 3 {
		t.Errorf("batch.Send(), expected %d recevied, got %d after expected interval expiry", 3, one.received)
	}

	c.Metrics = []*skogul.Metric{&m, &m, &m, &m, &m, &m, &m, &m, &m}
	err = batch.Send(&c)
	if err != nil {
		t.Errorf("batch.Send() failed: %v", err)
	}
	time.Sleep(time.Duration(5 * time.Millisecond))
	if one.received != 3 {
		t.Errorf("batch.Send(), sender 1 expected %d recevied, got %d", 3, one.received)
	}
	err = batch.Send(&c)
	if err != nil {
		t.Errorf("batch.Send() failed: %v", err)
	}
	time.Sleep(time.Duration(5 * time.Millisecond))
	if one.received != 4 {
		t.Errorf("batch.Send(), sender 1 expected %d recevied, got %d", 4, one.received)
	}
}
