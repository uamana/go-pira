package pira

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
)

func parseFloat64(data []byte) (float64, error) {
	v, _, _ := bytes.Cut(bytes.TrimSpace(data), []byte(" "))
	value, err := strconv.ParseFloat(string(v), 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func parseNullableFloat64(data []byte) Nullable[float64] {
	value, err := parseFloat64(data)
	if err != nil {
		return Nullable[float64]{Value: 0, Valid: false}
	}
	return Nullable[float64]{Value: value, Valid: true}
}

func parseInt(data []byte) (int, error) {
	v, _, _ := bytes.Cut(bytes.TrimSpace(data), []byte(" "))
	value, err := strconv.ParseInt(string(v), 10, 32)
	if err != nil {
		return 0, err
	}
	return int(value), nil
}

func parseNullableInt(data []byte) Nullable[int] {
	value, err := parseInt(data)
	if err != nil {
		return Nullable[int]{Value: 0, Valid: false}
	}
	return Nullable[int]{Value: value, Valid: true}
}

func parseString(data []byte) (string, error) {
	v := bytes.TrimSpace(data)
	return string(v), nil
}

func parseNullableString(data []byte) Nullable[string] {
	value, err := parseString(data)
	if err != nil {
		return Nullable[string]{Value: "", Valid: false}
	}
	return Nullable[string]{Value: value, Valid: true}
}

func parseBool(data []byte) (bool, error) {
	value, err := parseInt(data)
	if err != nil {
		return false, err
	}
	return value != 0, nil
}

func parseNullableBool(data []byte) Nullable[bool] {
	value, err := parseBool(data)
	if err != nil {
		return Nullable[bool]{Value: false, Valid: false}
	}
	return Nullable[bool]{Value: value, Valid: true}
}

func parseHistogramData(data []byte) (Histogram, error) {
	histogram := make(Histogram, 0, 122)
	var pair [2][]byte
	count := 0
	for b := range bytes.FieldsFuncSeq(data, func(r rune) bool {
		return !unicode.IsNumber(r) && r != '.'
	}) {
		pair[count] = b
		count++
		if count == 2 {
			count = 0
			binStr := string(pair[0])
			bin, err := strconv.ParseInt(string(binStr), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid bin string: %s", string(pair[0]))
			}
			value, err := strconv.ParseInt(string(pair[1]), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid value string: %s", string(pair[1]))
			}
			histogram = append(histogram, []int64{bin, value})
		}
	}
	return histogram, nil
}

func parseRDSGroupStatsData(data []byte) (RDSGroupStatsData, error) {
	stats := make(RDSGroupStatsData, 0, 30)
	var pair [2][]byte
	count := 0
	for b := range bytes.FieldsFuncSeq(data, func(r rune) bool {
		return unicode.IsSpace(r) || r == ';' || r == '%'
	}) {
		pair[count] = b
		count++
		if count == 2 {
			count = 0
			group := string(pair[0])
			value, err := parseFloat64(pair[1])
			if err != nil {
				return nil, fmt.Errorf("invalid value string: %s", string(pair[1]))
			}
			stats = append(stats, RDSGroupStatsDataItem{Group: group, Percent: float64(value)})
		}
	}
	return stats, nil
}
