package config

import (
	"encoding/json"
	"reflect"
	"testing"
	
	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/receiver"
)

func TestFile(t *testing.T) {
	defer func() {
		if skogul.AssertErrors > 0 {
			t.Errorf("File() paniced")
		}
	}()
	c, err := File("testdata/test.json")
	if err != nil {
		t.Errorf("File() failed: %v", err)
	}
	if c == nil {
		t.Errorf("File() returned nil config")
	}
}

func TestByte_ok(t *testing.T) {
	defer func() {
		if skogul.AssertErrors > 0 {
			t.Errorf("Byte() paniced")
		}
	}()
	var okData []byte
	okData = []byte(`
{
  "senders": {
    "tnet_alarms": {
      "type": "sql",
      "Driver": "mysql",
      "ConnStr": "root:lol@/mydb",
      "Query": "INERT INTO tnet_Alarms values(${foo})"
    },
    "mysql_log": {
      "type": "sql",
      "Driver": "mysql",
      "ConnStr": "root:lol@/mydb",
      "Query": "INSERT INTO liksomlog VALUES(${timestmap.timestamp},${metadata.name},${key})"
    },
    "forward": {
	    "type": "http",
	    "url": "http://localhost:8888"
    },
    "duplicate": {
	    "type": "dupe",
	    "next": ["forward","mysql_log","tnet_alarms"]
    },
    "batch": {
      "type": "batch",
      "next": "duplicate",
      "interval": "5s"
    },
    "det": {
	    "type": "detacher",
	    "next": "batch"
    }
  },
  "handlers": {
    "plain": {
      "parser": "json",
      "sender": "batch",
      "transformers": []
    }
  },
  "receivers": {
    "http": {
      "type": "http",
      "handlers": {
	      "/": "plain"
      }
    }
  }
}
`)
	c, err := Bytes(okData)
	if err != nil {
		t.Errorf("Bytes() failed: %v", err)
	}
	if c == nil {
		t.Errorf("Bytes() returned nil config")
	}
	badData := []byte(`{ "senders": { "x": { "type": "sql", "ConnStr": 5 } } }`)
	c, err = Bytes(badData)
	if err == nil {
		t.Errorf("Bytes() test 2 failed, sent bad data, didn't get error.")
	}
	if c != nil {
		t.Errorf("Bytes() with bad data returned valid config.")
	}
	noURL := []byte(`{ "senders": { "x" : { "type": "http" }}}`)
	c, err = Bytes(noURL)
	if err == nil {
		t.Errorf("Bytes() test 3 failed, http sender with no URL didn't get error.")
	}
	if c != nil {
		t.Errorf("Bytes() with bad data returned valid config.")
	}

}

func TestHelpSender(t *testing.T) {
	_, err := HelpSender("sql")
	if err != nil {
		t.Errorf("HelpSender(\"sql\") didn't work: %v", err)
	}
}

func TestFindSuperfluousReceiverConfigProperties(t *testing.T) {
	rawConfig := []byte(`{"receivers": {
		"foo": {
		  "type": "udp",
		  "address": "[::1]:5015",
		  "superfluousField": "this is not needed"
		}
	  }
	}`)

	var config map[string]interface{}
	err := json.Unmarshal(rawConfig, &config)

	if err != nil {
		t.Error("Failed to parse config")
	}

	configStruct := reflect.TypeOf(receiver.UDP{})
	superfluousProperties := verifyOnlyRequiredConfigProps(&config, "receivers", "foo", configStruct)

	if len(superfluousProperties) != 1 {
		t.Errorf("Expected 1 superfluous property but got %d", len(superfluousProperties))
	}

	if superfluousProperties[0] != "superfluousField" {
		t.Errorf("Expected to find '%s' in the superfluous fields list", "superfluousField")
	}
}
