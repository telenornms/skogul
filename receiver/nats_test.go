package receiver_test

import (
	"github.com/telenornms/skogul/config"
	"github.com/telenornms/skogul/receiver"
	"testing"
)

func TestNats(t *testing.T) {
	config, err := config.Bytes([]byte(`
	{
	  "receivers": {
	    "nats_r": {
	      "type": "nats",
	      "servers": "nats://0.0.0.0:4222",
	      "queue": "skogul_queue",
	      "subject": "test.subject",
	      "handler": "test_h"
	    }
	  },
	  "handlers": {
	    "test_h": {
	      "parser": "skogul",
	      "transformers": [],
	      "sender": "test"
	    }
	  },
	  "senders": {
	    "test": {
	      "type": "test"
	    }
	  }
	}`))

	if err != nil {
		t.Errorf("Failed to load config: %s", err)
	}

	nr := config.Receivers["nats_r"].Receiver.(*receiver.Nats)

	if nr == nil {
		t.Error("Bad config")
	}
}
