package engine

import "testing"

func TestExecutePING(t *testing.T) {
	tests := []struct {
		name    string
		command []string
		want    string
	}{
		{name: "simple ping", command: []string{"PING"}, want: "+PONG\r\n"},
		{name: "ping with payload", command: []string{"PING", "hello"}, want: "$5\r\nhello\r\n"},
		{name: "ping with too many args", command: []string{"PING", "one", "two"}, want: "-ERR wrong number of arguments for 'ping' command\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(Execute(tt.command)); got != tt.want {
				t.Fatalf("unexpected response: got %q want %q", got, tt.want)
			}
		})
	}
}

func TestExecuteSetGetDelExistsAndValidation(t *testing.T) {
	key := "engine:test:key"

	if got := string(Execute([]string{"SET", key, "value"})); got != "+OK\r\n" {
		t.Fatalf("unexpected SET response: %q", got)
	}

	if got := string(Execute([]string{"GET", key})); got != "$5\r\nvalue\r\n" {
		t.Fatalf("unexpected GET response: %q", got)
	}

	if got := string(Execute([]string{"EXISTS", key, "missing"})); got != ":1\r\n" {
		t.Fatalf("unexpected EXISTS response: %q", got)
	}

	if got := string(Execute([]string{"DEL", key})); got != ":1\r\n" {
		t.Fatalf("unexpected DEL response: %q", got)
	}

	if got := string(Execute([]string{"GET", key})); got != "$-1\r\n" {
		t.Fatalf("unexpected GET missing response: %q", got)
	}

	validationTests := []struct {
		name    string
		command []string
		want    string
	}{
		{name: "get missing args", command: []string{"GET"}, want: "-ERR wrong number of arguments for 'get' command\r\n"},
		{name: "set missing args", command: []string{"SET", key}, want: "-ERR wrong number of arguments for 'set' command\r\n"},
		{name: "del missing args", command: []string{"DEL"}, want: "-ERR wrong number of arguments for 'del' command\r\n"},
		{name: "exists missing args", command: []string{"EXISTS"}, want: "-ERR wrong number of arguments for 'exists' command\r\n"},
		{name: "unknown command", command: []string{"NOPE"}, want: "-ERR unknown command 'nope'\r\n"},
		{name: "empty command", command: []string{}, want: "-ERR empty command\r\n"},
	}

	for _, tt := range validationTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(Execute(tt.command)); got != tt.want {
				t.Fatalf("unexpected response: got %q want %q", got, tt.want)
			}
		})
	}
}