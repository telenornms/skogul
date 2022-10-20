package parser_test

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/telenornms/skogul/gen/usp"
	"github.com/telenornms/skogul/parser"
)

func readFile(file string, t *testing.T) []byte {
	t.Helper()
	b := make([]byte, 9000)
	f, err := os.Open(file)
	if err != nil {
		t.Fatalf("unable to open protobuf packet file: %v", err)
	}
	defer f.Close()
	n, err := f.Read(b)
	if err != nil {
		t.Fatalf("unable to read protobuf packet file: %v", err)
	}
	if n == 0 {
		t.Fatalf("read 0 bytes from protobuf packet file....")
	}
	return b[0:n]
}

func TestUSPParseFile(t *testing.T) {
	d := readFile("testdata/usp.bin", t)

	x := parser.USP_Parser{}
	container, err := x.Parse(d)

	if err != nil {
		t.Error("Error while parsing", err)
	}

	if container.Metrics[0].Metadata == nil || container.Metrics[0].Data == nil {
		t.Error("Either metadata or data is missing")
	}
}

func TestUSPGetUspRecord(t *testing.T) {
	d := readFile("testdata/usp.bin", t)

	unmarshaledMessage := &usp.Record{}
	if err := proto.Unmarshal(d, unmarshaledMessage); err != nil {
		t.Error("Error while unmarshalling data", err)
	}
}

func TestUSPGetRecordMsgPayload(t *testing.T) {
	msgPayload := &usp.Msg{}

	d := readFile("testdata/usp.bin", t)

	unmarshaledMessage := &usp.Record{}
	if err := proto.Unmarshal(d, unmarshaledMessage); err != nil {
		t.Error("Error while unmarshalling record", err)
	}

	payload := unmarshaledMessage.GetNoSessionContext().GetPayload()
	if err := proto.Unmarshal(payload, msgPayload); err != nil {
		t.Error("Failed to unmarshall record payload")
	}
}

func TestUSPExtractJSON(t *testing.T) {
	msgPayload := &usp.Msg{}

	d := readFile("testdata/usp.bin", t)

	unmarshaledMessage := &usp.Record{}
	if err := proto.Unmarshal(d, unmarshaledMessage); err != nil {
		t.Error("Error while unmarshalling record", err)
	}

	payload := unmarshaledMessage.GetNoSessionContext().GetPayload()
	if err := proto.Unmarshal(payload, msgPayload); err != nil {
		t.Error("Failed to unmarshall record payload")
	}

	input := []byte(msgPayload.Body.GetRequest().GetNotify().GetEvent().GetParams()["Data"])

	var k map[string]interface{}

	if err := json.Unmarshal(input, &k); err != nil {
		t.Error("Failed to unmarshall json")
	}
}

func TestUSPTransformData(t *testing.T) {
	msgPayload := &usp.Msg{}

	d := readFile("testdata/usp.bin", t)

	unmarshaledMessage := &usp.Record{}
	if err := proto.Unmarshal(d, unmarshaledMessage); err != nil {
		t.Error("Error while unmarshalling record", err)
	}

	payload := unmarshaledMessage.GetNoSessionContext().GetPayload()
	if err := proto.Unmarshal(payload, msgPayload); err != nil {
		t.Error("Failed to unmarshall record payload")
	}

	input := []byte(msgPayload.Body.GetRequest().GetNotify().GetEvent().GetParams()["Data"])

	var k map[string]interface{}

	if err := json.Unmarshal(input, &k); err != nil {
		t.Error("Failed to unmarshall json")
	}

	transformData(k["Report"])
}

func createValues(keys []string, index int,d map[string]interface{}) map[string]interface{} {
	tmp := make(map[string]interface{})
	for _, k := range keys {
		spl := strings.Split(k, ".")

		if len(spl) == 3 {
			in, _ := strconv.Atoi(spl[1])

			if in - 1 < 0 {
				in = 0
			} else {
				in = in - 1
			}

			if in == index {
				tmp[spl[2]] = d[k]
			}
		}
	}

	return tmp
}

func unflattenValues(d map[string]interface{}) map[string]interface{} {
	keys := make([]string, 0, len(d))

	for k := range d {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	tmp := map[string]interface{}{}

	for _, k := range keys {
		spl := strings.Split(k, ".")

		if len(spl) == 1 {
			tmp[spl[0]] = d[k]
		} else if len(spl) == 2 {
			if _, ok := tmp[spl[0]]; !ok {
				tmp[spl[0]] = make(map[string]interface{})
			}
			t := tmp[spl[0]].(map[string]interface{})
			t[spl[1]] = d[k]
			tmp[spl[0]] = t
		} else if len(spl) == 3 {
			if _, ok := tmp[spl[0]]; !ok {
				tmp[spl[0]] = make([]map[string]interface{}, 10)
			}

			x := tmp[spl[0]].([]map[string]interface{})
			index, _ := strconv.Atoi(spl[1])

			if index - 1 < 0 {
				index = 0
			} else {
				index = index - 1
			}

			x[index] = createValues(keys, index, d)

			tmp[spl[0]] = x
		}
	}

	return tmp
}

func transformData(j interface{}) {
	unflatten := unflattenValues(j.([]interface{})[0].(map[string]interface{}))

	jsonData, _ := json.MarshalIndent(unflatten, "", "	")

	log.Println(string(jsonData))
}
