package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"
)

//UR001!
func getInfo() ([]byte, error) {
	var buf bytes.Buffer
	fmt.Println("Message: Get All Info")
	fmt.Fprint(&buf, "UR001!")
	grillResponse, err := sendData(&buf)
	if err != nil {
		return nil, err
	}
	return grillResponse, nil
}

// UT###!
func setGrillTemp(temp int) ([]byte, error) {
	fmt.Println("Message: Set Grill Temp")
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
	fmt.Println("Message: Set Probe Temp")
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
	fmt.Println("Message: Turn Grill On")
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
	fmt.Println("Message: Turn Grill Off")
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
	fmt.Println("Message: Get Grill ID")
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
	fmt.Println("Message: Get Grill FW")
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
	fmt.Println("Connected")

	defer conn.Close()
	fmt.Println("Sending Data..")
	ret, err := conn.Write(b.Bytes())
	if err != nil {
		return nil, errors.New("Failure Sending Payload to Grill")
	}
	fmt.Printf("Bytes Written: %v\n", ret)
	b.Reset()

	fmt.Println("Reading Data..")
	barray := make([]byte, 1024)
	status, err := bufio.NewReader(conn).Read(barray)
	if err != nil {
		return nil, errors.New("Failed Reading Result From Grill")
	}
	// trim null of 1024 byte array
	//barray = bytes.Trim(barray, "\x00")
	barray = barray[:36]

	// print what we got back
	fmt.Println(string(b.Bytes()))
	fmt.Println(string(barray))
	fmt.Println(barray)
	fmt.Println("Bytes Read:", status)
	fmt.Println("Read Buffer Size:", len(barray))
	return barray, nil
}
