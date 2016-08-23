package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"
)

// gmg support
// UR[2 Byte Grill Temp][2 Byte food probe Temp][2 Byte Target Temp]
// [skip 22 bytes][2 Byte target food probe][1byte on/off/fan][5 byte tail]
/* byte value map
const (
	grillTemp        = 2
	probeTemp        = 4
	grillSetTemp     = 6
	curveRemainTime  = 20
	warnCode         = 24
	probeSetTemp     = 28
	grillState       = 30
	grillMode        = 31
	fireState        = 32
	fileStatePercent = 33
	profileEnd       = 34
	grillType        = 35
)
var grillStates = map[int]string{
	0: "OFF",
	1: "ON",
	2: "FAN",
	3: "REMAIN",
}
var fireStates = map[int]string{
	0: "DEFAULT",
	1: "OFF",
	2: "STARTUP",
	3: "RUNNING",
	4: "COOLDOWN",
	5: "FAIL",
}
var warnStates = map[int]string{
	0: "FAN_OVERLOADED",
	1: "AUGER_OVERLOADED",
	2: "IGNITOR_OVERLOADED",
	3: "BATTERY_LOW",
	4: "FAN_DISCONNECTED",
	5: "AUGER_DISCONNECTED",
	6: "IGNITOR_DISCONNECTED",
	7: "LOW_PELLET",
}
*/

//UR001!
func getInfo() ([]byte, error) {
	var buf bytes.Buffer
	fmt.Println("Request: Get All Info")
	fmt.Fprint(&buf, "UR001!")
	grillResponse, err := sendData(&buf)
	if err != nil {
		return nil, err
	}
	return grillResponse, nil
}

// UT###!
func setGrillTemp(temp int) ([]byte, error) {
	fmt.Println("Request: Set Grill Temp")
	single := temp % 10
	tens := (temp % 100) / 10
	hundreds := temp / 100
	b := []byte{
		byte(85),
		byte(84),
		byte(hundreds + 48),
		byte(tens + 48),
		byte(single + 48),
		byte(33),
	}
	var buf bytes.Buffer
	buf.Write(b)
	grillResponse, err := sendData(&buf)
	if err != nil {
		return nil, err
	}
	return grillResponse, nil
}

// UF###!
func setProbeTemp(temp int) ([]byte, error) {
	fmt.Println("Request: Set Probe Temp")
	single := temp % 10
	tens := (temp % 100) / 10
	hundreds := temp / 100
	b := []byte{
		byte(85),
		byte(70),
		byte(hundreds + 48),
		byte(tens + 48),
		byte(single + 48),
		byte(33),
	}
	var buf bytes.Buffer
	buf.Write(b)
	grillResponse, err := sendData(&buf)
	if err != nil {
		return nil, err
	}
	return grillResponse, nil
}

// UK001!
func powerOn() ([]byte, error) {
	var buf bytes.Buffer
	fmt.Println("Request: Turn Grill On")
	fmt.Fprint(&buf, "UK001!")
	grillResponse, err := sendData(&buf)
	if err != nil {
		return nil, err
	}
	return grillResponse, nil
}

//UK004!
func powerOff() ([]byte, error) {
	var buf bytes.Buffer
	fmt.Println("Request: Turn Grill Off")
	fmt.Fprint(&buf, "UK004!")
	grillResponse, err := sendData(&buf)
	if err != nil {
		return nil, err
	}
	return grillResponse, nil
}

//UL!
func grillID() ([]byte, error) {
	var buf bytes.Buffer
	fmt.Println("Request: Get Grill ID")
	fmt.Fprint(&buf, "UL!")
	grillResponse, err := sendData(&buf)
	if err != nil {
		return nil, err
	}
	return grillResponse, nil
}

//UN!
func grillFW() ([]byte, error) {
	var buf bytes.Buffer
	fmt.Println("Request: Get Grill FW")
	fmt.Fprint(&buf, "UN!")
	grillResponse, err := sendData(&buf)
	if err != nil {
		return nil, err
	}
	return grillResponse, nil
}

func sendData(b *bytes.Buffer) ([]byte, error) {
	//b = []byte("UWFM!") // leaveServerMode
	if b.Len() == 0 {
		return nil, errors.New("Nothing to Send to Grill")
	}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s", myGrill.grillIP), 3*time.Second)
	timeout := time.Now().Add(3 * time.Second)
	conn.SetReadDeadline(timeout)
	if err != nil {
		return nil, errors.New("Connection to Grill Failed")
	}
	//fmt.Println("Connected")

	defer conn.Close()
	//fmt.Println("Sending Data..")
	_, err = conn.Write(b.Bytes())
	if err != nil {
		return nil, errors.New("Failure Sending Payload to Grill")
	}
	//fmt.Printf("Bytes Written: %v\n", ret)
	b.Reset()

	//fmt.Println("Reading Data..")
	barray := make([]byte, 1024)
	_, err = bufio.NewReader(conn).Read(barray)
	if err != nil {
		return nil, errors.New("Failed Reading Result From Grill")
	}
	// trim null of 1024 byte array
	//barray = bytes.Trim(barray, "\x00")
	barray = barray[:36]

	// print what we got back
	/*
		fmt.Println(string(b.Bytes()))
		fmt.Println(string(barray))
		fmt.Println(barray)
		fmt.Println("Bytes Read:", status)
		fmt.Println("Read Buffer Size:", len(barray))
	*/
	return barray, nil
}
