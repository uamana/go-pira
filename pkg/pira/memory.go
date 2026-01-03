package pira

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math"
)

type MemoryPart1 struct {
	Frequency                uint16    // 0x01A, khz * 10, raized by 1065
	U1                       [4]byte   // 0x1C
	DIP                      uint16    // 0x020
	U2                       uint16    // 0x22
	PilotDeviation           uint16    // 0x024, Hz * 10
	RDSDeviation             uint16    // 0x026, Hz * 10
	PiotToRDSPhaseDifference int16     // 0x28
	DeviationMax             uint16    //0x2A
	DeviationAverage         uint16    //0x2C
	ModulationPower          uint16    //0x2E
	DeviationMinHold         uint16    //0x30
	RDSPI                    uint16    //0x32
	RDSPS                    [8]byte   //0x34
	RDSPTY                   byte      //0x3C
	U3                       byte      //
	RDSStatus                uint16    //0x3E
	RDSGroupCounters         [32]byte  //0x40
	RDSAFList                [26]byte  //0x60
	RDSEONPI                 [4]uint16 //0x7A
	SignalQuality            byte      //0x82
	U4                       [5]byte   //
	DeviationMaxHold         uint16    //0x88
	U5                       [4]byte   //
	AM                       byte      //0x8E
	U6                       [181]byte //0x11F
	Deviation                uint16    //0x144
	NoiseLevel               uint16    //0x146
	U7                       [84]byte  //
	RDSRT                    [64]byte  //0x19C
	RDSPTYN                  [8]byte   //0x1DC
	RDSCTHour                byte      //0x1E4
	U8                       byte      //0x1E5
	RDSCTMinute              byte      //0x1E6
	U9                       [3]byte   //0x1E7
	RDSMJD                   [3]byte   //0x1EA
	U10                      byte      //0x1ED
	RDSRTPlusGroupType       byte      //0x1EE
	RDSRTPlusStatus          byte      //0x1EF
	RDSRTPlusItem1Type       byte      //0x1F0
	RDSRTPlusItem1Start      byte      //0x1F1
	RDSRTPlusItem1Length     byte      //0x1F2
	RDSRTPlusItem2Type       byte      //0x1F3
	RDSRTPlusItem2Start      byte      //0x1F4
	RDSRTPlusItem2Length     byte      //0x1F5
	U11                      [2]byte   //0x1F6
	RDSPINDay                byte      //0x1F8
	RDSPINHour               byte      //0x1F9
	RDSPINMinute             byte      //0x1FA
	RDSLIC                   byte      //0x1FB
	RDSECC                   byte      //0x1FC
	RDSCTLocalTimeOffset     byte      //0x1FD
}

type MemoryPart2 struct {
	InstantModulationPower uint16      //0x48C
	U1                     [64]byte    //0x48E
	Alarms                 [13]byte    //0x4CE
	U2                     [151]byte   //
	HistogramData          [122]uint16 //0x572
	U3                     [266]byte   //0x5F4
	RDSLongPS              [32]byte    //0x770
}

func (p *Pira) Load(addr int, structure any) error {
	size := binary.Size(structure)
	if addr > 0xFFF || size > 0xFFF {
		panic("invalid address or size")
	}

	command := fmt.Sprintf("%03X,%03X?h", addr, size)
	n, err := p.SendCommand(Command(command))
	if err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}
	if n != len(command) {
		return fmt.Errorf("failed to send command: %w", err)
	}

	data, err := p.RecvResponse()
	if err != nil {
		return fmt.Errorf("failed to receive response: %w", err)
	}
	lines := bytes.Split(data, []byte("\r\n"))
	if len(lines) < 3 {
		return fmt.Errorf("invalid response: %s", string(data))
	}
	slog.Debug("response", "hex", string(lines[2]), "len", len(lines[2]))

	buffer := make([]byte, size)

	n, err = hex.Decode(buffer, lines[2])
	if err != nil {
		return fmt.Errorf("failed to decode hex memory (%d bytes): %w", len(buffer), err)
	}
	if n != size {
		return fmt.Errorf("invalid response: %s", string(data))
	}
	n, err = binary.Decode(buffer, binary.LittleEndian, structure)
	if err != nil {
		return fmt.Errorf("failed to decode hex memory (%d bytes): %w", len(buffer), err)
	}
	if n != size {
		return fmt.Errorf("invalid response: %s", string(data))
	}
	return nil
}

// parseFrequency converts frequency from kHz * 10 + 1065to kHz
func parseFrequency(frequency uint16) uint32 {
	return uint32(frequency-1065) / 10
}

func parsePiotToRDSPhaseDifference(piotToRDSPhaseDifference int16) int16 {
	return piotToRDSPhaseDifference - 90
}

// parseDeviation converts deviation from kHz * 10 to Hz
func parseDeviation(d uint16) uint32 {
	return uint32(d) * 100
}

// parseModulationPower converts linear modulation power into dBr
func parseModulationPower(modulationPower uint16) float64 {
	return 10 * math.Log10(float64(modulationPower)*100)
}

func parseRDSStatus(rdsStatus uint16) (status *RDSStatus) {
	var rtType RTType
	if rdsStatus&0b0100000000 != 0 {
		rtType = RTTypeA
	} else {
		rtType = RTTypeB
	}
	status = &RDSStatus{
		CT:     rdsStatus&0b10000000000 != 0,
		RT:     rdsStatus&0b01000000000 != 0,
		RTType: rtType,
		AF:     rdsStatus&0b00010000000 != 0,
		TP:     rdsStatus&0b00001000000 != 0,
		TA:     rdsStatus&0b00000100000 != 0,
		MS:     rdsStatus&0b00000010000 != 0,
		DI:     byte(rdsStatus & 0b00000001111),
	}
	return status
}

func (p *Pira) GetFrequency() (uint32, error) {
	var f uint16
	err := p.Load(0x01A, &f)
	if err != nil {
		return 0, fmt.Errorf("failed to get frequency: %w", err)
	}
	return parseFrequency(f), nil
}

type DeviationType int

const (
	DeviationPilot   DeviationType = 0x024
	DeviationRDS     DeviationType = 0x026
	DeviationMax     DeviationType = 0x02A
	DeviationAve     DeviationType = 0x02C
	DeviationMinHold DeviationType = 0x030
	DeviationMaxHold DeviationType = 0x088
	Deviation        DeviationType = 0x144
)

func (p *Pira) GetDeviation(dt DeviationType) (uint32, error) {
	var pd uint16
	err := p.Load(int(dt), &pd)
	if err != nil {
		return 0, fmt.Errorf("failed to get deviation: %w", err)
	}
	return parseDeviation(pd), nil
}

func (p *Pira) GetRDSPhaseDifference() (int16, error) {
	var rdsPhaseDifference int16
	err := p.Load(0x028, &rdsPhaseDifference)
	if err != nil {
		return 0, fmt.Errorf("failed to get rds phase difference: %w", err)
	}
	return parsePiotToRDSPhaseDifference(rdsPhaseDifference), nil
}

func (p *Pira) GetModulationPower() (float64, error) {
	var mp uint16
	err := p.Load(0x02E, &mp)
	if err != nil {
		return 0, fmt.Errorf("failed to get modulation power: %w", err)
	}
	return parseModulationPower(mp), nil
}

func (p *Pira) GetRDSPI() (uint16, error) {
	var rdsPI uint16
	err := p.Load(0x032, &rdsPI)
	if err != nil {
		return 0, fmt.Errorf("failed to get rds pi: %w", err)
	}
	return rdsPI, nil
}

func (p *Pira) GetRDSPS() (string, error) {
	var rdsPS [8]byte
	err := p.Load(0x034, &rdsPS)
	if err != nil {
		return "", fmt.Errorf("failed to get rds ps: %w", err)
	}
	return string(rdsPS[:]), nil
}

func (p *Pira) GetRDSPTY() (byte, error) {
	var rdsPTY byte
	err := p.Load(0x03C, &rdsPTY)
	if err != nil {
		return 0, fmt.Errorf("failed to get rds pty: %w", err)
	}
	return rdsPTY, nil
}

func (p *Pira) GetRDSStatus() (*RDSStatus, error) {
	var rdsStatus uint16
	err := p.Load(0x03E, &rdsStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to get rds status: %w", err)
	}
	return parseRDSStatus(rdsStatus), nil
}

func (p *Pira) GetRDSGroupCounters() ([32]byte, error) {
	var rdsGroupCounters [32]byte
	err := p.Load(0x040, &rdsGroupCounters)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to get rds group counters: %w", err)
	}
	return rdsGroupCounters, nil
}

func (p *Pira) GetRDSAFList() ([26]byte, error) {
	var rdsAFList [26]byte
	err := p.Load(0x060, &rdsAFList)
	if err != nil {
		return [26]byte{}, fmt.Errorf("failed to get rds af list: %w", err)
	}
	return rdsAFList, nil
}

func (p *Pira) GetRDSEONPI() ([4]uint16, error) {
	var rdsEONPI [4]uint16
	err := p.Load(0x07A, &rdsEONPI)
	if err != nil {
		return [4]uint16{}, fmt.Errorf("failed to get rds eonpi: %w", err)
	}
	return rdsEONPI, nil
}

func (p *Pira) GetSignalQuality() (int, error) {
	var signalQuality byte
	err := p.Load(0x082, &signalQuality)
	if err != nil {
		return 0, fmt.Errorf("failed to get signal quality: %w", err)
	}
	return int(signalQuality), nil
}

func (p *Pira) GetAM() (byte, error) {
	var am byte
	err := p.Load(0x08E, &am)
	if err != nil {
		return 0, fmt.Errorf("failed to get am: %w", err)
	}
	return am, nil
}

func (p *Pira) GetNoiseLevel() (uint16, error) {
	var noiseLevel uint16
	err := p.Load(0x146, &noiseLevel)
	if err != nil {
		return 0, fmt.Errorf("failed to get noise level: %w", err)
	}
	return noiseLevel, nil
}

func (p *Pira) GetRDSRT() (string, error) {
	var rdsRT [64]byte
	err := p.Load(0x19C, &rdsRT)
	if err != nil {
		return "", fmt.Errorf("failed to get rds rt: %w", err)
	}
	return string(rdsRT[:]), nil
}

func (p *Pira) GetRDSPTYN() (string, error) {
	var rdsPTYN [8]byte
	err := p.Load(0x1DC, &rdsPTYN)
	if err != nil {
		return "", fmt.Errorf("failed to get rds ptyn: %w", err)
	}
	return string(rdsPTYN[:]), nil
}

func (p *Pira) GetRDSCT() (*RDSCT, error) {
	var (
		data   [3]byte
		offset byte
		err    error
	)
	err = p.Load(0x1E4, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to get rds ct: %w", err)
	}
	err = p.Load(0x1FD, &offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get rds ct local time offset: %w", err)
	}
	rdsCT := &RDSCT{
		Hour:            data[0],
		Minute:          data[2],
		LocalTimeOffset: offset,
	}
	return rdsCT, nil
}

func (p *Pira) GetRDSMJD() ([3]byte, error) {
	var rdsMJD [3]byte
	err := p.Load(0x1EA, &rdsMJD)
	if err != nil {
		return [3]byte{}, fmt.Errorf("failed to get rds mjd: %w", err)
	}
	return rdsMJD, nil
}

func (p *Pira) GetRDSRTPlus() (*RTPlus, error) {
	var rdsRTPlus RTPlus
	err := p.Load(0x1EE, &rdsRTPlus)
	if err != nil {
		return nil, fmt.Errorf("failed to get rds rt plus: %w", err)
	}
	return &rdsRTPlus, nil
}

func (p *Pira) GetRDSPIN() (*RDSPIN, error) {
	var rdsPIN RDSPIN
	err := p.Load(0x1F8, &rdsPIN)
	if err != nil {
		return nil, fmt.Errorf("failed to get rds pin: %w", err)
	}
	return &rdsPIN, nil
}

func (p *Pira) GetRDSLIC() (byte, error) {
	var rdsLIC byte
	err := p.Load(0x1FB, &rdsLIC)
	if err != nil {
		return 0, fmt.Errorf("failed to get rds lic: %w", err)
	}
	return rdsLIC, nil
}

func (p *Pira) GetRDSECC() (byte, error) {
	var rdsECC byte
	err := p.Load(0x1FC, &rdsECC)
	if err != nil {
		return 0, fmt.Errorf("failed to get rds ecc: %w", err)
	}
	return rdsECC, nil
}

func (p *Pira) GetRDSCTLocalTimeOffset() (byte, error) {
	var rdsCTLocalTimeOffset byte
	err := p.Load(0x1FD, &rdsCTLocalTimeOffset)
	if err != nil {
		return 0, fmt.Errorf("failed to get rds ct local time offset: %w", err)
	}
	return rdsCTLocalTimeOffset, nil
}

func (p *Pira) GetHistogram() ([]uint16, error) {
	var histogram [122]uint16
	err := p.Load(0x572, &histogram)
	if err != nil {
		return nil, fmt.Errorf("failed to get histogram: %w", err)
	}
	return histogram[:], nil
}

func (p *Pira) GetRDSLongPS() (string, error) {
	var rdsLongPS [32]byte
	err := p.Load(0x770, &rdsLongPS)
	if err != nil {
		return "", fmt.Errorf("failed to get rds long ps: %w", err)
	}
	return string(rdsLongPS[:]), nil
}

func (p *Pira) GetFMInfo(fmi *FMInfo) (err error) {
	mem1 := MemoryPart1{}
	err = p.Load(0x01A, &mem1)
	if err != nil {
		return fmt.Errorf("failed to get fm info: %w", err)
	}
	mem2 := MemoryPart2{}
	err = p.Load(0x48C, &mem2)
	if err != nil {
		return fmt.Errorf("failed to get fm info: %w", err)
	}

	fmi.Frequency = parseFrequency(mem1.Frequency)
	fmi.PilotDeviation = parseDeviation(mem1.PilotDeviation)
	fmi.RDSDeviation = parseDeviation(mem1.RDSDeviation)
	fmi.RDSPhaseDifference = parsePiotToRDSPhaseDifference(mem1.PiotToRDSPhaseDifference)
	fmi.DeviationMax = parseDeviation(mem1.DeviationMax)
	fmi.DeviationAverage = parseDeviation(mem1.DeviationAverage)
	fmi.ModulationPower = parseModulationPower(mem1.ModulationPower)
	fmi.DeviationMinHold = parseDeviation(mem1.DeviationMinHold)

	fmi.RDS.Status = *parseRDSStatus(mem1.RDSStatus)
	fmi.RDS.Groups = mem1.RDSGroupCounters
	fmi.RDS.EONPI = mem1.RDSEONPI
	fmi.RDS.RT = string(mem1.RDSRT[:])
	fmi.RDS.PTYN = string(mem1.RDSPTYN[:])
	fmi.RDS.CT.Hour = mem1.RDSCTHour
	fmi.RDS.CT.Minute = mem1.RDSCTMinute
	fmi.RDS.CT.LocalTimeOffset = mem1.RDSCTLocalTimeOffset
	fmi.RDS.MJD = mem1.RDSMJD
	fmi.RDS.PIN.Day = mem1.RDSPINDay
	fmi.RDS.PIN.Hour = mem1.RDSPINHour
	fmi.RDS.PIN.Minute = mem1.RDSPINMinute
	fmi.RDS.LIC = mem1.RDSLIC
	fmi.RDS.ECC = mem1.RDSECC

	fmi.RDS.RTPlus.GroupType = mem1.RDSRTPlusGroupType
	fmi.RDS.RTPlus.Status = mem1.RDSRTPlusStatus
	fmi.RDS.RTPlus.Item1.Type = mem1.RDSRTPlusItem1Type
	fmi.RDS.RTPlus.Item1.Start = mem1.RDSRTPlusItem1Start
	fmi.RDS.RTPlus.Item1.Length = mem1.RDSRTPlusItem1Length
	fmi.RDS.RTPlus.Item2.Type = mem1.RDSRTPlusItem2Type
	fmi.RDS.RTPlus.Item2.Start = mem1.RDSRTPlusItem2Start
	fmi.RDS.RTPlus.Item2.Length = mem1.RDSRTPlusItem2Length

	fmi.RDS.LongPS = string(mem2.RDSLongPS[:])

	fmi.SignalQuality = int(mem1.SignalQuality)
	fmi.DeviationMaxHold = parseDeviation(mem1.DeviationMaxHold)
	fmi.AM = mem1.AM
	fmi.Deviation = parseDeviation(mem1.Deviation)
	fmi.NoiseLevel = mem1.NoiseLevel
	fmi.Histogram = mem2.HistogramData
	return nil
}
