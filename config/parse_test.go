package config

import (
	"testing"
	//"fmt"
)

func TestFile(t *testing.T) {
	c, err := File("test.json")
	if err != nil {
		t.Errorf("File() failed: %v", err)
	}
	if c == nil {
		t.Errorf("File() returned nil config")
	}
}

func TestByte_ok(t *testing.T) {
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
