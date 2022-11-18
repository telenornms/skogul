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

	c := skogul.Container{}
	c.Metrics = []*skogul.Metric{&metric, &metric2, &metric3, &metric4}

	ban := &transformer.Ban{}

	ban.Lookup = map[string]interface{}{
		"/baz/bar/baz":   "baz.bar.baz",
		"/bar/bar/1/321": "12345",
		"/funny":         "notfunny",
	}

	t.Log(c)
	err := ban.Transform(&c)
	if err != nil {
		t.Fatalf("error occurred %v", err.Error())
	}
	t.Log(c)

	if baz, ok := c.Metrics[1].Data["baz"]; ok {
		r := baz.(map[string]interface{})
		r2 := r["bar"].(map[string]interface{})

		if _, ok := r2["baz"]; ok {
			t.Fatal("/baz/bar/baz has not been deleted")
		}
	}

}
