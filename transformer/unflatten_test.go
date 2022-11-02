package transformer_test

import (
	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/transformer"
	"testing"
)

func TestTransformData(t *testing.T) {
	testData := map[string]interface{}{
		"foo": "bar",
		"bar.baz": "bar.baz",
		"bar.foo": "bar.foo",
		"bar.1.baz.1.foo": "bar.1.baz.1.foo",
		"bar.1.baz.1.bar": "bar.1.baz.1.bar",
		"bar.1.foo.1.baz.1.foo": "bar.1.foo.1.baz.1.foo",
	}

	metric := skogul.Metric{
		Time: nil,
		Metadata: map[string]interface{}{
			"foo": "bar",
		},
		Data: testData,
	}

	metrics := []*skogul.Metric{
		0: nil,
	}
	metrics[0] = &metric

	container := &skogul.Container{
		Template: nil,
		Metrics:  metrics,
	}

	u := &transformer.Unflatten{}
	err := u.Transform(container)

	if err != nil {
		t.Errorf("%v", err)
	}

    data := container.Metrics[0].Data["bar"].(map[string]interface{})

    if _, ok := data["1"]; !ok {
        t.Error("faile to create nested structure")
    }
}
