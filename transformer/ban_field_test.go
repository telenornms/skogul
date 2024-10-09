package transformer_test

import (
	"testing"

	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/transformer"
)

func TestBanField(t *testing.T) {
	metric := skogul.Metric{}
	metric.Metadata = make(map[string]interface{})
	metric.Metadata["foofoo"] = "barBAR"
	metric.Data = make(map[string]interface{})
	metric.Data["foo"] = "BAR"
	metric.Data["baz"] = "foobar"
	c := skogul.Container{}
	c.Metrics = []*skogul.Metric{&metric}

	ban := transformer.BanField{
		SourceData:     "foo",
		RegexpData:     "BAR",
		SourceMetadata: "foofoo",
		RegexpMetadata: "barBAR",
	}

	t.Logf("Container before transform:\n%v", c)
	err := ban.Transform(&c)
	if err != nil {
		t.Errorf("ban_field returned non-nil err: %v", err)
	}

	t.Logf("Container after transform:\n%v", c)

	if _, ok := c.Metrics[0].Metadata["foofoo"]; ok {
		t.Fatal("ban_field transformer failed to ban key-value pair")
	}
	if _, ok := c.Metrics[0].Data["foo"]; ok {
		t.Fatal("ban_field transformer failed to ban key-value pair")
	}
}
