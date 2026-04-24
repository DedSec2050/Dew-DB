package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:6379", "Dew-DB server address")
	flag.Parse()

	command := flag.Args()
	if len(command) == 0 {
		command = []string{"PING"}
	}

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	if _, err := conn.Write(encodeArray(command)); err != nil {
		fmt.Fprintf(os.Stderr, "write error: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(conn)
	response, err := readRespValue(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(response)
}

func encodeArray(parts []string) []byte {
	var b strings.Builder
	b.WriteString("*")
	b.WriteString(strconv.Itoa(len(parts)))
	b.WriteString("\r\n")

	for _, p := range parts {
		b.WriteString("$")
		b.WriteString(strconv.Itoa(len(p)))
		b.WriteString("\r\n")
		b.WriteString(p)
		b.WriteString("\r\n")
	}

	return []byte(b.String())
}

func readRespValue(reader *bufio.Reader) (string, error) {
	prefix, err := reader.ReadByte()
	if err != nil {
		return "", err
	}

	switch prefix {
	case '+':
		line, err := readLine(reader)
		if err != nil {
			return "", err
		}
		return line, nil
	case '-':
		line, err := readLine(reader)
		if err != nil {
			return "", err
		}
		return "(error) " + line, nil
	case ':':
		line, err := readLine(reader)
		if err != nil {
			return "", err
		}
		return line, nil
	case '$':
		line, err := readLine(reader)
		if err != nil {
			return "", err
		}

		size, err := strconv.Atoi(line)
		if err != nil {
			return "", fmt.Errorf("invalid bulk length: %w", err)
		}
		if size == -1 {
			return "(nil)", nil
		}

		payload := make([]byte, size+2)
		if _, err := io.ReadFull(reader, payload); err != nil {
			return "", err
		}
		if payload[size] != '\r' || payload[size+1] != '\n' {
			return "", errors.New("malformed bulk string terminator")
		}

		return string(payload[:size]), nil
	case '*':
		line, err := readLine(reader)
		if err != nil {
			return "", err
		}

		count, err := strconv.Atoi(line)
		if err != nil {
			return "", fmt.Errorf("invalid array length: %w", err)
		}
		if count == -1 {
			return "(nil)", nil
		}

		items := make([]string, 0, count)
		for i := 0; i < count; i++ {
			item, err := readRespValue(reader)
			if err != nil {
				return "", err
			}
			items = append(items, item)
		}

		return strings.Join(items, " "), nil
	default:
		return "", fmt.Errorf("unsupported RESP prefix %q", prefix)
	}
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	if !strings.HasSuffix(line, "\r\n") {
		return "", errors.New("line must end with CRLF")
	}

	return strings.TrimSuffix(line, "\r\n"), nil
}