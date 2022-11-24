package clon

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestCLON(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   []string
		out  any
	}{
		{"basic", []string{
			"name=John",
			"email=john@example.com",
		}, map[string]any{
			"name":  "John",
			"email": "john@example.com",
		}},

		{"basic raw values", []string{
			"name=John",
			"age:=29",
			"married:=false",
			"hobbies:='[\"http\", \"pies\"]'",
			"favorite:='{\"tool\": \"HTTPie\"}'",
		}, map[string]any{
			"age":     29,
			"hobbies": []any{"http", "pies"},
			"married": false,
			"name":    "John",
			"favorite": map[string]any{
				"tool": "HTTPie",
			},
		}},

		{"basic nested", []string{
			"platform[name]=HTTPie",
			"platform[about][mission]='Make APIs simple and intuitive'",
			"platform[about][homepage]=httpie.io",
			"platform[about][stars]:=54000",
			"platform[apps][]=Terminal",
			"platform[apps][]=Desktop",
			"platform[apps][]=Web",
			"platform[apps][]=Mobile",
		}, map[string]any{
			"platform": map[string]any{
				"name": "HTTPie",
				"about": map[string]any{
					"mission":  "Make APIs simple and intuitive",
					"homepage": "httpie.io",
					"stars":    54000,
				},
				"apps": []any{
					"Terminal",
					"Desktop",
					"Web",
					"Mobile",
				},
			},
		}},

		{"nested example", []string{
			"category=tools",
			"search[type]=id",
			"search[id]:=1",
		}, map[string]any{
			"category": "tools",
			"search": map[string]any{
				"id":   1,
				"type": "id",
			},
		}},

		{"append example", []string{
			"category=tools",
			"search[type]=keyword",
			"search[keywords][]=APIs",
			"search[keywords][]=CLI",
		}, map[string]any{
			"category": "tools",
			"search": map[string]any{
				"keywords": []any{
					"APIs",
					"CLI",
				},
				"type": "keyword",
			},
		}},

		{"indexed example", []string{
			"category=tools",
			"search[type]=keyword",
			"search[keywords][1]=APIs",
			"search[keywords][0]=CLI",
		}, map[string]any{
			"category": "tools",
			"search": map[string]any{
				"keywords": []any{
					"CLI",
					"APIs",
				},
				"type": "keyword",
			},
		}},

		{"indexed and appened example", []string{
			"category=tools",
			"search[type]=platforms",
			"search[platforms][]=Terminal",
			"search[platforms][1]=Desktop",
			"search[platforms][3]=Mobile",
		}, map[string]any{
			"category": "tools",
			"search": map[string]any{
				"platforms": []any{
					"Terminal",
					"Desktop",
					nil,
					"Mobile",
				},
				"type": "platforms",
			},
		}},

		{"raw and appended example", []string{
			"category=tools",
			"search[type]=platforms",
			"search[platforms]:='[\"Terminal\", \"Desktop\"]'",
			"search[platforms][]=Web",
			"search[platforms][]=Mobile",
		}, map[string]any{
			"category": "tools",
			"search": map[string]any{
				"platforms": []any{
					"Terminal",
					"Desktop",
					"Web",
					"Mobile",
				},
				"type": "platforms",
			},
		}},

		{"top level array", []string{
			"[]:=1",
			"[]:=2",
			"[]:=3",
		}, []any{
			1,
			2,
			3,
		}},

		{"top level array shorthand", []string{
			":1",
			"string",
			":true",
			`:{"foo": "bar"}`,
		}, []any{
			1,
			"string",
			true,
			map[string]any{"foo": "bar"},
		}},

		{"top level array nested", []string{
			"[0][type]=platform",
			"[0][name]=terminal",
			"[1][type]=platform",
			"[1][name]=desktop",
		}, []any{
			map[string]any{
				"type": "platform",
				"name": "terminal",
			},
			map[string]any{
				"type": "platform",
				"name": "desktop",
			},
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.out) {
				t.Fatalf("expected '%s' but got '%s'", dump(tt.out), dump(got))
			}
		})

	}
}

func dump(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
