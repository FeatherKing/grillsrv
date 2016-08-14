package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	grillID = iota
	useWifi = iota
)

func main() {
	var message = flag.Int("messageType", -1,
		"0 to print grill id. 1 to switch from ptp to wifi")
	flag.Parse()
	var buf bytes.Buffer
	ssid := "SSID"
	password := "WIFI_PASS"
	//serverip := "52.26.201.234"
	//port := "8060"
	ssidlen := len(ssid)
	passwordlen := len(password)
	//serveriplen := len(serverip)
	//portlen := len(port)

	//s := fmt.Sprintf("UH%c%c%s%c%s!", 0, ssidlen, ssid, passwordlen, password) // connect to wifi
	//s := fmt.Sprintf("UG%c%s%c%s%c%s%c%s!", ssidlen, ssid, passwordlen, password, serveriplen, serverip, portlen, port)
	//bytes := []byte(s)
	fmt.Println(*message, grillID, useWifi)
	switch *message {
	case grillID:
		// get grill id
		fmt.Println("Message: Get Grill Id")
		fmt.Fprint(&buf, "UL!")
	case useWifi:
		fmt.Fprintf(&buf, "UH%c%c%s%c%s!", 0, ssidlen, ssid, passwordlen, password)
	default:
		fmt.Println("You must choose a message type.")
		os.Exit(1)
	}

	sendData(&buf)
	//sendData(bytes)
}

func sendData(b *bytes.Buffer) {
	//b = []byte("UWFM!") // leaveServerMode
	//str := "UL!" // get grill id
	//str := "GET / HTTP/1.0\r\n\r\n"
	conn, err := net.Dial("tcp", "LAN_IP:PORT")
	timeout := time.Now().Add(3 * time.Second)
	conn.SetReadDeadline(timeout)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Connected")

	defer conn.Close()
	fmt.Println("Sending Data..")
	ret, err := conn.Write(b.Bytes())
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Bytes Written: %v\n", ret)
	b.Reset()

	fmt.Println("Reading Data..")
	barray := make([]byte, 1024)
	status, err := bufio.NewReader(conn).Read(barray)
	//bufio.NewReader(conn).ReadBytes(33)
	if err != nil {
		// prints a timeout error but data was still read
		fmt.Println(err.Error())
	}
	// trim null of 1024 byte array
	barray = bytes.Trim(barray, "\x00")

	// print what we got back
	fmt.Println(string(b.Bytes()))
	fmt.Println(string(barray))
	fmt.Println(barray)
	fmt.Println("Bytes Read:", status)
	fmt.Println("Read Buffer Size:", len(barray))

	//fmt.Println(grillIdCommand)
	/*this.prefs.edit().putString("serverModeKey", new StringBuilder(String.valueOf(DataManager.getInstance().mGrillId.substring(NUM_PAGES, 11))).append(DataManager.getInstance().mExternalIp).toString()).apply();
	  this.prefs.edit().putString("serverModeKey",
	    new StringBuilder(
	      String.valueOf(DataManager.getInstance().mGrillId.substring(NUM_PAGES, 11))
	      ).append(DataManager.getInstance().mExternalIp
	    ).toString()
	    ).apply();


		fmt.Println(b)
		fmt.Println("String: ", string(b[:]))
	}

	/*
			   case PubDefine.MESSAGE_GRILL_CONNECT_SERVER *115:
			       ssid = DataManager.getInstance().mSsid;
			       password = DataManager.getInstance().mPassword;
			       serverip = PubDefine.Server_AWS;
			       if (DataManager.getInstance().mServerIP.length() > 6) {
			           serverip = DataManager.getInstance().mServerIP;
			       }
			       port = String.format("%d", new Object[]{Integer.valueOf(PubDefine.Server_AWS_Port)});
			       if (DataManager.getInstance().mServerPort.length() > 1) {
			           port = DataManager.getInstance().mServerPort;
			       }
			       ssidlen = ssid.length();
			       passwordlen = password.length();
			       serveriplen = serverip.length();
			       portlen = port.length();
			       return String.format("UG%c%s%c%s%c%s%c%s!", new Object[]{Integer.valueOf(ssidlen), ssid, Integer.valueOf(passwordlen), password, Integer.valueOf(serveriplen), serverip, Integer.valueOf(portlen), port}).getBytes();
		         package main

		         import (
		         	"fmt"
		         	"encoding/binary"
		         	"bytes"
		         )
	*/
}
