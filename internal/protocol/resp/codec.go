package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ReadCommand(reader *bufio.Reader) ([]string, error) {
	firstByte, err := reader.Peek(1)
	if err != nil {
		return nil, err
	}

	if firstByte[0] == '*' {
		return readArrayCommand(reader)
	}

	line, err := readLine(reader)
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	return fields, nil
}

func SimpleString(value string) []byte {
	return []byte("+" + value + "\r\n")
}

func Error(value string) []byte {
	return []byte("-" + value + "\r\n")
}

func BulkString(value string) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value))
}

func NullBulkString() []byte {
	return []byte("$-1\r\n")
}

func Integer(value int) []byte {
	return []byte(fmt.Sprintf(":%d\r\n", value))
}

func readArrayCommand(reader *bufio.Reader) ([]string, error) {
	line, err := readLine(reader)
	if err != nil {
		return nil, err
	}

	if len(line) < 2 || line[0] != '*' {
		return nil, fmt.Errorf("expected array header")
	}

	count, err := strconv.Atoi(line[1:])
	if err != nil || count <= 0 {
		return nil, fmt.Errorf("invalid array length")
	}

	parts := make([]string, 0, count)
	for i := 0; i < count; i++ {
		bulkHeader, err := readLine(reader)
		if err != nil {
			return nil, err
		}

		if len(bulkHeader) < 2 || bulkHeader[0] != '$' {
			return nil, fmt.Errorf("expected bulk string")
		}

		bulkLen, err := strconv.Atoi(bulkHeader[1:])
		if err != nil || bulkLen < 0 {
			return nil, fmt.Errorf("invalid bulk string length")
		}

		payload := make([]byte, bulkLen+2)
		if _, err := io.ReadFull(reader, payload); err != nil {
			return nil, err
		}

		if payload[bulkLen] != '\r' || payload[bulkLen+1] != '\n' {
			return nil, fmt.Errorf("malformed bulk string terminator")
		}

		parts = append(parts, string(payload[:bulkLen]))
	}

	return parts, nil
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	if !strings.HasSuffix(line, "\r\n") {
		return "", fmt.Errorf("line must end with CRLF")
	}

	return strings.TrimSuffix(line, "\r\n"), nil
}