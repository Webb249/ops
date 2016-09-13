package opsviewplugin

import (
	"fmt"
	"os"
	"strings"
)

const (
	messageSeparator = ", "
)

// Standalone Exit function for simple checks without multiple results or perfdata.
func Exit(status Status, message string) {
	fmt.Printf("%v: %s\n", status, message)
	os.Exit(int(status))
}

// Represents the state of the check.
type Check struct {
	shortname	string
	results  	[]Result
	perfdata 	[]PerfData
	status   	Status
}

// NewCheck returns an empty Check object.
func NewCheck(shortname string) *Check {
	c := new(Check)
	c.shortname = shortname
	return c
}

// AddResult adds a check result. This will not terminate the check. If
// status is the highest yet reported, this will update the check's
// final return status.
func (c *Check) AddResult(status Status, message string) {
	var result Result
	result.status = status
	result.message = message
	c.results = append(c.results, result)
	if result.status > c.status {
		c.status = result.status
	}
}

// AddResultf functions as AddResult, but takes a printf-style format
// string and arguments.
func (c *Check) AddResultf(status Status, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	c.AddResult(status, msg)
}

// AddPerfData adds a metric to the set output by the check. unit must
// be a valid monitoring plugin unit of measurement (UOM): "us", "ms", "s",
// "%", "b", "kb", "mb", "gb", "tb", "c", or the empty string. UOMs are
// not case-sensitive.
//
// Zero or more of the thresholds min, max, warn and crit may be
// supplied; these must be of the same UOM as the value.
//
// A threshold may be positive or negative infinity, in which case it
// will be omitted in the check output. A value may not be either
// infinity.
//
// Returns error on invalid parameters.
func (c *Check) AddPerfData(label, unit string, value float64, thresholds ...float64) error {
	data, err := NewPerfData(label, unit, value, thresholds...)
	if err != nil {
		return err
	}
	c.perfdata = append(c.perfdata, *data)
	return nil
}

// exitInfoText returns the most important result text, formatted for
// the first line of plugin output.
//
// Returns joined string of (messageSeparator-separated) info text from
// results which have a status of at least c.status.
func (c Check) exitInfoText() string {
	importantMessages := make([]string, 0)
	for _, result := range c.results {
		if result.status == c.status {
			importantMessages = append(importantMessages, result.message)
		}
	}
	return strings.Join(importantMessages, messageSeparator)
}

// String representation of the check results, suitable for output and
// parsing by Nagios.
func (c Check) String() string {
	value := fmt.Sprintf("%v: %s", c.status, c.exitInfoText())
	value += RenderPerfdata(c.perfdata)
	return value
}

// Finish ends the check, prints its output (to stdout), and exits with
// the correct status.
func (c *Check) Finish() {
	if r := recover(); r != nil {
		c.Exitf(CRITICAL, "%v - check panicked: %v", c.shortname, r)
	}
	if len(c.results) == 0 {
		c.AddResult(UNKNOWN, "no check result specified")
	}
	fmt.Println(c)
	os.Exit(int(c.status))
}

// Exitf takes a status plus a format string, and a list of
// parameters to pass to Sprintf. It then immediately outputs and exits.
func (c *Check) Exitf(status Status, format string, v ...interface{}) {
	info := fmt.Sprintf(format, v...)
	c.AddResult(status, info)
	c.Finish()
}

// Criticalf is a shorthand function which exits the check with status
// CRITICAL and the message provided.
func (c *Check) Criticalf(format string, v ...interface{}) {
	c.Exitf(CRITICAL, format, v...)
}

// Unknownf is a shorthand function which exits the check with status
// UNKNOWN and the message provided.
func (c *Check) Unknownf(format string, v ...interface{}) {
	c.Exitf(UNKNOWN, format, v...)
}
