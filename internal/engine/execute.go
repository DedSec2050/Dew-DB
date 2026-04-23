package engine

import (
	"fmt"
	"strings"

	"github.com/DedSec2050/dew-db/internal/protocol/resp"
)

func Execute(command []string) []byte {
	if len(command) == 0 {
		return resp.Error("ERR empty command")
	}

	switch strings.ToUpper(command[0]) {
	case "PING":
		return executePing(command)
	default:
		return resp.Error(fmt.Sprintf("ERR unknown command '%s'", strings.ToLower(command[0])))
	}
}

func executePing(command []string) []byte {
	if len(command) == 1 {
		return resp.SimpleString("PONG")
	}

	if len(command) == 2 {
		return resp.BulkString(command[1])
	}

	return resp.Error("ERR wrong number of arguments for 'ping' command")
}