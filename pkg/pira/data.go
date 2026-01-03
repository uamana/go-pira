package pira

import (
	"encoding/json"
	"fmt"
)

type Nullable[T any] struct {
	Value T
	Valid bool
}

func (n *Nullable[T]) IsValid() bool {
	return n.Valid
}

func (n *Nullable[T]) IsNull() bool {
	return !n.Valid
}

func (n Nullable[T]) String() string {
	if !n.Valid {
		return "not set"
	}
	return fmt.Sprintf("%v", n.Value)
}

// MarshalJSON implements the json.Marshaler interface for Nullable[T]
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Value)
}

// UnmarshalJSON implements the json.Unmarshaler interface for Nullable[T]
func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		var zero T
		n.Value = zero
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		n.Valid = false
		var zero T
		n.Value = zero
		return err
	}
	n.Value = v
	n.Valid = true
	return nil
}

type Command string

const (
	CmdGetBasicData Command = "?B"
	CmdGetFrequency Command = "?F"
)

type DataKey string

const (
	KeyFrequency          DataKey = "frequency:"
	KeyPilot              DataKey = "pilot:"
	KeyRDSDeviation       DataKey = "rds deviation:"
	KeyRDSGroupStats      DataKey = "rds group statistics:"
	KeyDeviationMax       DataKey = "max:"
	KeyDeviationMin       DataKey = "min:"
	KeyDeviationAvg       DataKey = "ave:"
	KeySignalQuality      DataKey = "signal quality:"
	KeyModulationPower    DataKey = "modulation power:"
	KeyRDSPhaseDifference DataKey = "rds phase difference:"
	KeyHistogramData      DataKey = "histogram data:"
	KeyFFT                DataKey = "fft:"
	KeyDeviationMaxHold   DataKey = "max hold:"
	KeySignalLevel        DataKey = "signal level:"
	KeyAM                 DataKey = "am:"
	KeyBalance            DataKey = "r/l:"
)

type Histogram [][]int64

type BasicData struct {
	Frequency          float64
	SignalQuality      int
	Pilot              Nullable[float64]
	RDSDeviation       Nullable[float64]
	RDSPhaseDifference Nullable[float64]
	ModulationPower    Nullable[float64]
	HistogramData      Histogram
	RDSGroupStatsData  RDSGroupStatsData
}

type RDSGroupStatsDataItem struct {
	Group   string
	Percent float64
}

type RDSGroupStatsData []RDSGroupStatsDataItem

type RTType byte

const (
	RTTypeA RTType = 0
	RTTypeB RTType = 1
)

func (rt RTType) String() string {
	if rt == RTTypeA {
		return "RTA"
	}
	return "RTB"
}

type RDSStatus struct {
	CT     bool   `json:"ct"`
	RT     bool   `json:"rt"`
	RTType RTType `json:"rt_type"`
	AF     bool   `json:"af"`
	TP     bool   `json:"tp"`
	TA     bool   `json:"ta"`
	MS     bool   `json:"ms"`
	DI     byte   `json:"di"`
}

// RTPlusItem is a single RT+ item
type RTPlusItem struct {
	Type   byte `json:"type"`
	Start  byte `json:"start"`
	Length byte `json:"length"`
}

// RTPlus is the RT+ struct
type RTPlus struct {
	GroupType byte       `json:"group_type"`
	Status    byte       `json:"status"`
	Item1     RTPlusItem `json:"item_1"`
	Item2     RTPlusItem `json:"item_2"`
}

type RDSPIN struct {
	Day    byte `json:"day"`
	Hour   byte `json:"hour"`
	Minute byte `json:"minute"`
}

func (r *RDSPIN) String() string {
	return fmt.Sprintf("%02d %02d:%02d", r.Day, r.Hour, r.Minute)
}

// RDSCT is the RDS Current Time
type RDSCT struct {
	Hour            byte `json:"hour"`
	Minute          byte `json:"minute"`
	LocalTimeOffset byte `json:"local_time_offset"`
}

func (r *RDSCT) String() string {
	return fmt.Sprintf("%02d:%02d (%02d)", r.Hour, r.Minute, r.LocalTimeOffset)
}

type RDSInfo struct {
	PI     uint16      `json:"pi"`
	PS     string      `json:"ps"`
	PTY    byte        `json:"pty"`
	Status RDSStatus   `json:"status"`
	Groups [32]byte    `json:"groups"`
	AFList [26]float64 `json:"af_list"`
	EONPI  [4]uint16   `json:"eonpi"`
	RT     string      `json:"rt"`
	PTYN   string      `json:"ptyn"`
	CT     RDSCT       `json:"ct"`
	MJD    [3]byte     `json:"mjd"`
	RTPlus RTPlus      `json:"rt_plus"`
	PIN    RDSPIN      `json:"pin"`
	LIC    byte        `json:"lic"`
	ECC    byte        `json:"ecc"`
	LongPS string      `json:"long_ps"`
}

func (r *RDSInfo) String() string {
	return fmt.Sprintf(
		"PI: %d, PS: %s, PTY: %d, Status: %v, Groups: %v, AFList: %v, "+
			"EONPI: %v, RT: %s, PTYN: %s, CT: %v, MJD: %v, RTPlus: %v, PIN: %v, "+
			"LIC: %d, ECC: %d",
		r.PI, r.PS, r.PTY, r.Status, r.Groups, r.AFList, r.EONPI, r.RT,
		r.PTYN, r.CT, r.MJD, r.RTPlus, r.PIN, r.LIC, r.ECC,
	)
}

type FMInfo struct {
	Frequency          uint32
	PilotDeviation     uint32
	RDSDeviation       uint32
	RDSPhaseDifference int16
	DeviationMax       uint32
	DeviationAverage   uint32
	ModulationPower    float64
	DeviationMinHold   uint32
	RDS                RDSInfo
	SignalQuality      int
	DeviationMaxHold   uint32
	AM                 byte
	Deviation          uint32
	NoiseLevel         uint16
	Histogram          [122]uint16
}
