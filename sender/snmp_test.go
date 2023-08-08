package sender_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/config"
	"github.com/telenornms/skogul/parser"
	"github.com/telenornms/skogul/receiver"
	"github.com/telenornms/skogul/sender"
)

func TestSnmpSenderTest(t *testing.T) {
	t.Skip()
	c, _ := loadJsonFile(t, "ble ble")

	sconf := fmt.Sprintln(`
{
  "receivers": {
    "http": {
      "type": "http",
      "address": "localhost:1337",
      "handlers": {
				"/": "h"
			}
    }
  },
	"transformers": {
		"flatten_data": {
			"type": "data",
      "flattenSeparator": "drop",
      "flatten": [
        [
          "kek"
        ]
      ]
		},
    "remove_data": {
      "type": "data",
      "remove": [
        "kek"
      ]
    }
  },
  "handlers": {
    "h": {
      "parser": "skogul",
      "transformers": [
        "now",
				"flatten_data",
				"remove_data"
      ],
      "sender": "snmp"
    }
  },
  "senders": {
    "snmp": {
      "type": "snmp",
			"port": 7331,
			"community": "xxx",
			"version": "2c",
			"target": "localhost",
			"oidmap": {
				"kek": "1.3.3.7",
			}
    }
  }
}`)

	conf, err := config.Bytes([]byte(sconf))

	if err != nil {
		t.Errorf("Failed to load config: %v", err)
		return
	}

	snmpSender := conf.Senders["snmp"].Sender.(*sender.SNMP)
	httpRcv := conf.Receivers["http"].Receiver.(*receiver.HTTP)

	for _, hr := range httpRcv.Handlers {
		hr.H.Transform(c)
	}

	if httpRcv == nil {
		t.Errorf("failed to get receiver")
		return
	}

	go httpRcv.Start()
	time.Sleep(time.Duration(20 * time.Millisecond))

	err = snmpSender.Send(c)

	if err != nil {
		log.Println(err)
	}
}

func loadJsonFile(t *testing.T, file string) (*skogul.Container, []byte) {
	b, _ := ioutil.ReadFile(file)

	container, _ := parser.SkogulJSON{}.Parse(b)

	return container, b
}
