package pira

import (
	"bufio"
	"bytes"
	"fmt"
	"time"

	"go.bug.st/serial"
)

type Pira struct {
	port     string
	baudRate int
	conn     serial.Port
	reader   *bufio.Reader
}

func Dial(port string, baudRate int, timeout time.Duration) (*Pira, error) {
	conn, err := serial.Open(port, &serial.Mode{BaudRate: baudRate})
	if err != nil {
		return nil, fmt.Errorf("failed to open port %s: %w", port, err)
	}
	err = conn.SetReadTimeout(timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to set read timeout: %w", err)
	}
	return &Pira{
		port:     port,
		baudRate: baudRate,
		conn:     conn,
		reader:   bufio.NewReader(conn),
	}, nil
}

func (p *Pira) Close() error {
	return p.conn.Close()
}

func (p *Pira) SendCommand(command Command) (int, error) {
	n, err := p.conn.Write([]byte(command))
	if err != nil {
		return 0, err
	}
	err = p.conn.Drain()
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (p *Pira) RecvResponse() ([]byte, error) {
	buf := make([]byte, 0, 1024)
	wb := bytes.NewBuffer(buf)

	lineCount := 0
	for {
		line, err := p.reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		_, err = wb.Write(line)
		lineCount++
		if err != nil {
			return nil, err
		}
		if bytes.Equal(line, []byte("\r\n")) && lineCount > 2 {
			return wb.Bytes(), nil
		}
	}
}
