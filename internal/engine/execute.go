package engine

import (
	"fmt"
	"strings"

	"github.com/DedSec2050/dew-db/internal/protocol/resp"
	"github.com/DedSec2050/dew-db/internal/storage"
)

var kv = storage.NewStore()

func Execute(command []string) []byte {
	if len(command) == 0 {
		return resp.Error("ERR empty command")
	}

	switch strings.ToUpper(command[0]) {
	case "PING":
		return executePing(command)
	case "GET":
		return executeGet(command)
	case "SET":
		return executeSet(command)
	case "DEL":
		return executeDel(command)
	case "EXISTS":
		return executeExists(command)
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

func executeGet(command []string) []byte {
	if len(command) != 2 {
		return resp.Error("ERR wrong number of arguments for 'get' command")
	}

	value, ok := kv.Get(command[1])
	if !ok {
		return resp.NullBulkString()
	}

	return resp.BulkString(value)
}

func executeSet(command []string) []byte {
	if len(command) != 3 {
		return resp.Error("ERR wrong number of arguments for 'set' command")
	}

	kv.Set(command[1], command[2])
	return resp.SimpleString("OK")
}

func executeDel(command []string) []byte {
	if len(command) < 2 {
		return resp.Error("ERR wrong number of arguments for 'del' command")
	}

	return resp.Integer(kv.Del(command[1:]...))
}

func executeExists(command []string) []byte {
	if len(command) < 2 {
		return resp.Error("ERR wrong number of arguments for 'exists' command")
	}

	return resp.Integer(kv.Exists(command[1:]...))
}