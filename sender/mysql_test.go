package sender

import (
	"github.com/KristianLyng/skogul"
	"testing"
	"time"
)

func TestMysql(t *testing.T) {
	m := Mysql{Query: "INSERT INTO test VALUES(${timestamp.timestamp},${metadata.src},${name},${data});", ConnStr: "root:lol@/skogul"}
	err := m.Init()
	if err != nil {
		t.Errorf("Mysql.Init failed: %v", err)
	}
	want := "INSERT INTO test VALUES(?,?,?,?);"
	if want != m.q {
		t.Errorf("Mysql.Init wanted %s got %s", want, m.q)
	}

	c := skogul.Container{}
	me := skogul.Metric{}
	n := time.Now()
	me.Time = &n
	me.Metadata = make(map[string]interface{})
	me.Data = make(map[string]interface{})
	me.Metadata["src"] = "Test"
	me.Data["name"] = "Foo Bar"
	me.Data["data"] = "something"
	c.Metrics = []*skogul.Metric{&me}

	err = m.Send(&c)
	if err != nil {
		t.Errorf("Mysql.Send failed: %v", err)
	}
}

func ExampleMysql() {
	m := skogul.Mysql{Query: "INSERT INTO test VALUES(${timestamp.timestamp},${metadata.src},${name},${data});", ConnStr: "user:password@/db"}
	// Actual query to be prepared will be INSERT INTO test VALUES(?, // ?, ?, ?);
}
