package opsviewplugin

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"strconv"
	"errors"
	"encoding/json"
	)

// Used to fetch metrics from a connection
// Give it the address without the http://
func Fetch(HostAddress string) string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://"+HostAddress, nil)
	ExitAtError(err, "Cannot create connection request to "+HostAddress)

	resp, err := client.Do(req)
	ExitAtError(err, "Could Not Connect to "+HostAddress)

	if resp.StatusCode != 200 {
		Exit(CRITICAL, "Could Not Connect to "+HostAddress+" Response Code: "+string(resp.Status))
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	ExitAtError(err, "Error reading output from "+HostAddress)

	return string(bodyText)
}

// Used to fetch metrics from a connection that requires authentication
// Give it the address without the http://
func FetchWithAuth(HostAddress string, Username string, Password string) string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", HostAddress, nil)
	ExitAtError(err, "Cannot create connection request to "+HostAddress)

	if (Username != "") && (Password != "") {
		req.SetBasicAuth(Username, Password)
	}
	resp, err := client.Do(req)
	ExitAtError(err, "Could Not Connect to "+HostAddress)

	if resp.StatusCode != 200 {
		Exit(CRITICAL, "Could Not Connect to "+HostAddress+" Response Code: "+string(resp.Status))
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	ExitAtError(err, "Error reading output from "+HostAddress)

	return string(bodyText)
}

// Accepts a warning and critical level to test against a value given to be evaluated.
// Tests if the value is higher than the warning and critical levels unless reverse is equal to true
// then it will test if the value is lower than warning and critical.
// Returns a status code for OK, WARNING and CRITICAL.
func Evaluate(warning int, critical int, value int, reverse bool) Status {
	if reverse == false {
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
	if (floatValue < 0) || (floatValue > 100) {
		Exit(CRITICAL, "Error converting values to percentage")
	}
  return int(floatValue)
}

func GetPercentOfF(used float64, total float64)float64{
	floatValue := (used / total) * 100
	if (floatValue < 0) || (floatValue > 100) {
		Exit(CRITICAL, "Error converting values to percentage")
	}
  return floatValue
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

func RemoveDecimal(value float64) float64 {
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

// Checks if the port is empty and if it is a value that can be the port number
func CheckPort(portValue string) {
	port, err := strconv.Atoi(portValue)
	ExitAtError(err, "Error converting port to integer")
	if (portValue != "") && ((port < 1) || (port > 65535)) {
		err = errors.New("Port check error")
		ExitAtError(err, "Port must be between 1 and 65535")
	}
}

// Checks the flag has a value
// Arguments are in the order flagValue, flagName, flagCharater repeated
func CheckFlags(args ...string) error {
	for i := range args {
		if args[i] == "" {
			return errors.New(args[i+1]+" (-"+args[i+2]+") is a required argument")
		}
		i =+ 2
	}
	return nil
}

func DecodeJsonToMap(bodyText string) (map[string]interface{}, error){
	var m interface{}

	json.Unmarshal([]byte(bodyText), &m)
	if mapValue, ok := m.(map[string]interface{}); ok {
		return mapValue, nil
	} else {
			return nil, errors.New("Error Converting Body Text to string")
	}
	return nil, nil
}
