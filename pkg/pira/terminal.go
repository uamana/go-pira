package pira

import (
	"bytes"
	"fmt"
	"log/slog"

	"go.bug.st/serial"
)

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
