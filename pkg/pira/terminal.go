package pira

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
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

func (p *Pira) GetBasicData() (*BasicData, error) {
	n, err := p.SendCommand(CmdGetBasicData)
	if err != nil {
		return nil, err
	}
	if n != len(CmdGetBasicData) {
		return nil, fmt.Errorf("failed to send command: %w", err)
	}

	basicData := BasicData{}

	for {
		response, err := p.RecvResponse()
		slog.Debug("response", "response", string(response))

		if err != nil {
			if portErr, ok := err.(*serial.PortError); ok && portErr.Code() == serial.ReadTimeout {
				break
			}
			return nil, err
		}
		key, data, _ := bytes.Cut(response, []byte("\r\n"))
		switch DataKey(string(bytes.TrimSpace(bytes.ToLower(key)))) {
		case KeyFrequency:
			basicData.Frequency, err = parseFloat64(data)
			if err != nil {
				return nil, err
			}
		case KeySignalQuality:
			basicData.SignalQuality, err = parseInt(data)
			if err != nil {
				return nil, err
			}
		case KeyModulationPower:
			basicData.ModulationPower = parseNullableFloat64(data)
		case KeyPilot:
			basicData.Pilot = parseNullableFloat64(data)
		case KeyRDSDeviation:
			basicData.RDSDeviation = parseNullableFloat64(data)
		case KeyRDSPhaseDifference:
			basicData.RDSPhaseDifference = parseNullableFloat64(data)
		case KeyHistogramData:
			basicData.HistogramData, err = parseHistogramData(data)
			if err != nil {
				return nil, err
			}
		case KeyRDSGroupStats:
			basicData.RDSGroupStatsData, err = parseRDSGroupStatsData(data)
			if err != nil {
				return nil, err
			}
		}
	}
	return &basicData, nil
}
