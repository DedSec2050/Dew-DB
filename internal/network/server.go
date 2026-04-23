package network

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/DedSec2050/dew-db/internal/engine"
	"github.com/DedSec2050/dew-db/internal/protocol/resp"
)

type Server struct {
	Addr string
}

func (s Server) Run(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer ln.Close()

	go func() {
		<-ctx.Done()
		_ = ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return fmt.Errorf("accept: %w", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		command, err := resp.ReadCommand(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			_, _ = conn.Write(resp.Error("ERR protocol error: " + err.Error()))
			return
		}

		_, err = conn.Write(engine.Execute(command))
		if err != nil {
			return
		}
	}
}