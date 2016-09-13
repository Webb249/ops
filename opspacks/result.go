package opsviewplugin

// Plugin exit status.
type Status uint

const (
	OK Status = iota
	WARNING
	CRITICAL
	UNKNOWN
)

// Returns string representation of a Status. Panics if given an invalid
// status (this will be recovered in check.Finish if it has been deferred).
func (s Status) String() string {
	switch s {
	case OK:
		return "OK"
	case WARNING:
		return "WARNING"
	case CRITICAL:
		return "CRITICAL"
	case UNKNOWN:
		return "UNKNOWN"
	}
	panic("Invalid plugin.Status.")
}

// Result is a combination of a Status and infotext. A check can have
// multiple of these, and only the most important (greatest badness)
// will be reported on the first line of output or represented in the
// plugin's exit status.
type Result struct {
	status  Status
	message string
}
