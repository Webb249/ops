# opsviewplugin

Package for writing Opsview plugins in Go (golang).

# Usage example

The general usage pattern looks like this:

```go
func main() {
	// Initialize the check - this will return an UNKNOWN result
	// until more results are added.
	check := plugin.NewCheck()
	// If we exit early or panic() we'll still output a result.
	defer check.Finish()

	// obtain data here

	// Add an 'OK' result - if no 'worse' check results have been
	// added, this is the one that will be output.
	check.AddResult(opsviewplugin.OK, "everything looks shiny, cap'n")
	// Add some perfdata too (label, unit, value, min, max,
	// warn, crit). The math.Inf(1) will be parsed as 'no
	// maximum'.
	check.AddPerfDatum("badness", "kb", 3.14159, 0.0, math.Inf(1), 8000.0, 9000.0)

	// Parse an range from the command line and the more severe
	// results if they match.
	warnRange, err := opsviewplugin.ParseRange( "1:2" )
	if err != nil {
		check.AddResult(plugin.UNKNOWN, "error parsing warning range")
	}
	if warnRange.Check( 3.14159 ) {
		check.AddResult(opsviewplugin.WARNING, "Are we crashing again?")
	}
}
```

