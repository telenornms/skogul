package transformer_test

import (
    "github.com/telenornms/skogul/transformer"
    "testing"

	"github.com/telenornms/skogul"
)

func TestBan(t *testing.T) {
	metric := skogul.Metric{}
    metric.Data = map[string]interface{}{
        "foo": nil,
        "bar": map[string]interface{}{
            "baz": "bar.baz",
        },
        "baz": 1,
    }

	metric.Data["foo"] = map[string]interface{}{
        "bar": map[string]interface{}{
            "baz": "foo.bar.baz",
            "1": map[string]interface{}{
                "baz": "foo.bar.1.baz",
            },
        },
        "foobar": map[string]interface{}{
            "1": map[string]interface{}{
                "baz": map[string]interface{}{
                    "bar": "foo.foobar.1.baz.bar",
                },
            },
        },
    }

    c := skogul.Container{}
	c.Metrics = []*skogul.Metric{&metric}

    ban := &transformer.Ban{}

    ban.DPaths = map[string]interface{}{
        "foo.bar.1.baz": "foo",
        "foo.foobar.1.baz.bar": "foo.foobar.1.baz.bar",
        "bar.baz": "bar.baz",
        "baz": 9,
    }

    err := ban.Transform(&c)
    if err != nil {
        t.Fatalf("error occurred %v", err.Error())
    }

    t.Log(c)
}