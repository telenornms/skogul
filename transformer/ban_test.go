package transformer_test

import (
	"testing"

	"github.com/telenornms/skogul/transformer"

	"github.com/telenornms/skogul"
)

func TestBan(t *testing.T) {
	metric := skogul.Metric{
		Data: map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": map[string]interface{}{
					"baz": "foo.bar.baz",
					"1": map[string]interface{}{
						"baz": "foo.bar.1.baz",
						"ble": "ble123",
						"321": "12345",
					},
					"2": map[string]interface{}{
						"baz": "foo.bar.2.baz",
						"ble": "ble123",
						"321": "12345",
					},
					"3": map[string]interface{}{
						"baz": "foo.bar.1.baz",
						"ble": "ble123",
						"321": "12345",
					},
					"4": map[string]interface{}{
						"baz": "foo.bar.4.baz",
						"ble": "ble1234",
						"321": "12345",
					},
				},
			},
		},
	}

	metric2 := skogul.Metric{
		Data: map[string]interface{}{
			"baz": map[string]interface{}{
				"bar": map[string]interface{}{
					"baz": "baz.bar.baz",
				},
				"foobar": map[string]interface{}{
					"1": map[string]interface{}{
						"baz": map[string]interface{}{
							"bar": "foo.foobar.1.baz.bar",
						},
					},
				},
			},
		},
	}

	metric3 := skogul.Metric{
		Data: map[string]interface{}{
			"bar": map[string]interface{}{
				"bar": map[string]interface{}{
					"baz": "foo.bar.baz",
					"1": map[string]interface{}{
						"baz": "foo.bar.1.baz",
						"ble": "ble123",
						"321": "12345",
					},
					"2": map[string]interface{}{
						"baz": "foo.bar.2.baz",
						"ble": "ble123",
						"321": "12345",
					},
					"3": map[string]interface{}{
						"baz": "foo.bar.1.baz",
						"ble": "ble123",
						"321": "12345",
					},
					"4": map[string]interface{}{
						"baz": "foo.bar.4.baz",
						"ble": "ble1234",
						"321": "12345",
					},
				},
				"foobar": map[string]interface{}{
					"1": map[string]interface{}{
						"baz": map[string]interface{}{
							"bar": "foo.foobar.1.baz.bar",
						},
					},
				},
			},
			"fee": "5",
		},
	}

	metric4 := skogul.Metric{
		Metadata: map[string]interface{}{
			"funny": "notfunny",
		},
	}
	metric5 := skogul.Metric{
		Metadata: map[string]interface{}{
			"funny": "",
		},
	}
	metric6 := skogul.Metric{
		Metadata: map[string]interface{}{
			"funny": "",
		},
	}
	metric7 := skogul.Metric{
		Metadata: map[string]interface{}{
			"bar2": "hmm",
		},
		Data: map[string]interface{}{
			"foo2": "dette er 1339, two steps ahead",
		},
	}
	metric8 := skogul.Metric{
		Metadata: map[string]interface{}{
			"bar2": "hmm",
		},
		Data: map[string]interface{}{
			"foo2": "dette er 1337, akkurat passe",
		},
	}
	metric9 := skogul.Metric{
		Metadata: map[string]interface{}{
			"bar2": "1234578901337890",
		},
		Data: map[string]interface{}{
			"foo2": "dette er 1335, litt veikt",
		},
	}
	metric10 := skogul.Metric{
		Metadata: map[string]interface{}{
			"bar2": []byte("1234578901337890"),
		},
		Data: map[string]interface{}{
			"foo2": "dette er 1335, litt veikt",
		},
	}

	c := skogul.Container{}
	c.Metrics = []*skogul.Metric{&metric, &metric2, &metric3, &metric4, &metric5, &metric6, &metric7, &metric8, &metric9, &metric10}

	ban := &transformer.Ban{}

	ban.LookupData = map[string]interface{}{
		"/baz/bar/baz":   "baz.bar.baz",
		"/bar/bar/1/321": "12345",
		"/foo/bar/baz":   "foo.bar.baz",
	}

	ban.LookupMetadata = map[string]interface{}{
		"/funny": "",
	}

	ban.RegexpData = map[string]string{
		"/foo2": ".*1337.*",
	}
	ban.RegexpMetadata = map[string]string{
		"/bar2": ".*1337.*",
	}

	err := ban.Transform(&c)
	if err != nil {
		t.Fatalf("error occurred %v", err.Error())
	}
	if len(c.Metrics) != 2 {
		for _, x := range c.Metrics {
			t.Logf("metric left: %#v", x)
		}
		t.Fatalf("expected exactly 1 metric to remain, got %d", len(c.Metrics))
	}
	if cap(c.Metrics) != 2 {
		t.Fatalf("expected exactly len(metrics) == 1, got %d", cap(c.Metrics))
	}
}
