package opsviewplugin

import "fmt"

// Accepts a warning and critical level to test against a value given to be evaluated.
// Tests if the value is higher than the warning and critical levels unless reverse is equal to true
// then it will test if the value is lower than warning and critical.
// Returns a status code for OK, WARNING and CRITICAL.
func Evaluate(warning int, critical int, value int, reverse bool) Status {
	if reverse == true {
		if value >= critical {
			return CRITICAL
		} else if value >= warning {
			return WARNING
		}
	} else {
		if value <= critical {
			return CRITICAL
		} else if value <= warning {
			return WARNING
		}
	}
	return OK
}

func GetPercentOf(used int, total int)int{
  floatValue := (float64(used) / float64(total)) * float64(100)
  return int(floatValue)
}

func GetPercentOfF(used float64, total float64)float64{
  return (used / total) * 100
}


func ExitAtError(err error, message string) {
	if err != nil {
		Exit(CRITICAL, message)
	}
}

func ExitAtErrorPrint(err error, message string) {
	if err != nil {
		fmt.Println(err)
		Exit(CRITICAL, message)
	}
}

func removeDecimal(value float64) float64 {
	intValue := int(value)
	return float64(intValue)
}

// Converts Int to kb
func ConvertToKiloBytes(value int) int {
  returnValue := (value / 1024)
  return returnValue
}

// Converts float to kb
func ConvertToKiloBytesf(value float64) float64 {
  returnValue :=  value / 1024
  return returnValue
}

// Converts Int to mb
func ConvertToMegaBytes(value int) int {
  returnValue :=  (value / 1024) / 1024
  return returnValue
}

// Converts float to mb
func ConvertToMegaBytesf(value float64) float64 {
  returnValue :=  (value / 1024) / 1024
  return returnValue
}

// Converts Int to gb
func ConvertToGigaBytes(value int) int {
  returnValue :=  ((value / 1024) / 1024) / 1024
  return returnValue
}

// Converts float to gb
func ConvertToGigaBytesf(value float64) float64 {
  returnValue :=  ((value / 1024) / 1024) / 1024
  return returnValue
}
