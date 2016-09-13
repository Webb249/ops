package opsviewplugin

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// PerfDatum represents one metric to be reported as part of a check
// result.
type PerfData struct {
	label string
	value float64
	unit  string
	min   *float64
	max   *float64
	warn  *float64
	crit  *float64
}

// fmtPerfFloat returns a string representation of n formatted in the
// typical /\d+(\.\d+)/ pattern. The difference from %f is that it
// removes any trailing zeroes (like %g except it never returns
// values in scientific notation).
func fmtPerfFloat(n float64) string {
	return strconv.FormatFloat(n, 'f', -1, 64)
}

// validUnit returns true if the string is a valid UOM; otherwise false.
// It is case-insensitive.
func validUnit(unit string) bool {
	switch strings.ToLower(unit) {
	case "", "us", "ms", "s", "%", "b", "kb", "mb", "gb", "tb", "c":
		return true
	}
	return false
}

// NewPerfData returns a PerfData object suitable to use in a check
// result. unit must a valid Nagios unit, i.e., one of "us", "ms", "s",
// "%", "b", "kb", "mb", "gb", "tb", "c", or the empty string.
//
// Zero to four thresholds may be supplied: min, max, warn and crit.
// Thresholds may be positive infinity, negative infinity, or NaN, in
// which case they will be omitted in check output.
func NewPerfData(label string, unit string, value float64, thresholds ...float64) (*PerfData, error) {
	data := new(PerfData)
	data.label = label
	data.value = value
	data.unit = unit
	if !validUnit(unit) {
		return nil, fmt.Errorf("Invalid unit %v", unit)
	}
	if math.IsInf(value, 0) || math.IsNaN(value) {
		return nil, fmt.Errorf("Perfdata value may not be infinity or NaN: %v.", value)
	}
	if len(thresholds) >= 1 {
		data.min = &thresholds[0]
	}
	if len(thresholds) >= 2 {
		data.max = &thresholds[1]
	}
	if len(thresholds) >= 3 {
		data.warn = &thresholds[2]
	}
	if len(thresholds) >= 4 {
		data.crit = &thresholds[3]
	}
	return data, nil
}

// isThresholdSet returns true if one of min, max, warn or crit are set
// and false otherwise. They are determined to be 'set' if they are not
// a) the nil pointer, b) infinity (positive or negative) or c) NaN.
func isThresholdSet(t *float64) bool {
	switch {
	case t == nil:
		return false
	case math.IsInf(*t, 0):
		return false
	case math.IsNaN(*t):
		return false
	}
	return true
}

// fmtThreshold returns a string representation of min, max, warn or
// crit (whether or not they are set).
func fmtThreshold(t *float64) string {
	if !isThresholdSet(t) {
		return ""
	}
	return fmtPerfFloat(*t)
}

// String returns the string representation of a PerfData, suitable for
// check output.
func (p PerfData) String() string {
	val := fmtPerfFloat(p.value)
	value := fmt.Sprintf("%s=%s%s", p.label, val, p.unit)
	value += fmt.Sprintf(";%s;%s", fmtThreshold(p.warn), fmtThreshold(p.crit))
	value += fmt.Sprintf(";%s;%s", fmtThreshold(p.min), fmtThreshold(p.max))
	return value
}

// RenderPerfdata accepts a slice of PerfData objects and returns their
// concatenated string representations in a form suitable to append to
// the first line of check output.
func RenderPerfdata(perfdata []PerfData) string {
	value := ""
	if len(perfdata) == 0 {
		return value
	}
	// Demarcate start of perfdata in check output.
	value += " |"
	for _, data := range perfdata {
		value += fmt.Sprintf(" %v", data)
	}
	return value
}
