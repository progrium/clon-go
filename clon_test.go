package clon

import (
	"reflect"
	"testing"
)

func TestCLON(t *testing.T) {
	for _, tt := range []struct {
		in  []string
		out any
	}{
		{[]string{
			"name=John",
			"email=john@example.com",
		}, map[string]any{
			"name":  "John",
			"email": "john@example.com",
		}},
		{[]string{
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
		{[]string{
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
		{[]string{
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
		{[]string{
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
		{[]string{
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
		{[]string{
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
					"",
					"Mobile",
				},
				"type": "platforms",
			},
		}},
		{[]string{
			"category=tools",
			"search[type]=platforms",
			"search[platforms][]:='[\"Terminal\", \"Desktop\"]'",
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
		{[]string{
			"[]:=1",
			"[]:=2",
			"[]:=3",
		}, []any{
			1,
			2,
			3,
		}},
		{[]string{
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
		got, err := Parse(tt.in)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, tt.out) {
			t.Fatalf("expected '%#v' but got '%#v'", tt.out, got)
		}
	}
}
