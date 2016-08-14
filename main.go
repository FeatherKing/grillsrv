package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	grillID    = iota
	useWifi    = iota
	serverKey  = iota
	serverMode = iota
)

// Grill ...
// struct
type Grill struct {
	IP          string
	serial      string
	ssid        string
	password    string
	ssidlen     int
	passwordlen int
	serverip    string
	port        string
	serveriplen int
	portlen     int
}

func main() {
	var message = flag.Int("messageType", -1,
		"0 to print grill id. 1 to switch from ptp to wifi")
	flag.Parse()
	myGrill := Grill{
		serial:   "GMGSERIAL",
		ssid:     "SSID",
		password: "WIFI_PASS",
		serverip: "52.26.201.234",
		port:     "8060",
	}
	myGrill.ssidlen = len(myGrill.ssid)
	myGrill.passwordlen = len(myGrill.password)
	myGrill.serveriplen = len(myGrill.serverip)
	myGrill.portlen = len(myGrill.port)
	var buf bytes.Buffer

	switch *message {
	case grillID:
		// get grill id
		fmt.Println("Message: Get Grill Id")
		fmt.Fprint(&buf, "UL!")
	case useWifi:
		fmt.Println("Message: PTP to Wifi")
		fmt.Fprintf(&buf, "UH%c%c%s%c%s!", 0, myGrill.ssidlen, myGrill.ssid, myGrill.passwordlen, myGrill.password)
	case serverMode:
		fmt.Println("Message: Wifi to Server Mode")
		fmt.Fprintf(&buf, "UG%c%s%c%s%c%s%c%s!", myGrill.ssidlen, myGrill.ssid, myGrill.passwordlen, myGrill.password, myGrill.serveriplen, myGrill.serverip, myGrill.portlen, myGrill.port)
	case serverKey:
		fmt.Println("Message: Create Server Key")
		// curl 'https://api.ipify.org?format=json'
		r, err := http.Get("https://api.ipify.org?format=json")
		if err != nil {
			fmt.Println(err.Error())
		}
		defer r.Body.Close()
		err = json.NewDecoder(r.Body).Decode(&myGrill)
		serverKey := []byte(fmt.Sprint(myGrill.serial, myGrill.IP))
		fmt.Println("Serial:", myGrill.serial)
		fmt.Println("IP:", myGrill.IP)
		fmt.Println("ServerKey Bytes:", serverKey)
		fmt.Println("ServerKey:", fmt.Sprint(myGrill.serial, myGrill.IP))

		os.Exit(0) // remove this once key is build correctly
	default:
		fmt.Println("You must choose a message type.")
		os.Exit(1)
	}
	sendData(&buf)
}

func sendData(b *bytes.Buffer) {
	//b = []byte("UWFM!") // leaveServerMode
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
	if err != nil {
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
