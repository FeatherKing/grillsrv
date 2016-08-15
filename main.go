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
	"strings"
	"time"
)

const (
	grillID    = iota
	useWifi    = iota
	serverKey  = iota
	serverMode = 100 // setting high value until this works
	grillInfo  = iota
	externalip = iota
)

// Grill ...
// struct
type Grill struct {
	grillIP     string
	ExternalIP  string `json:"ip"`
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
		fmt.Sprint("0: Print grill id\n",
			"\t1: Switch from ptp to wifi\n",
			"\t2: Generate and Show ServerKey\n",
			"\t3: Send Server Mode Command (wip)\n",
			"\t4: Get Grill Info (wip)\n",
			"\t5: Request External IP From Grill\n\t"))
	flag.IntVar(message, "m", -1, "Shortversion of messageType")
	flag.Parse()
	myGrill := Grill{
		grillIP:  "192.168.0.10",
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
		ptp := false
		iface, err := net.InterfaceAddrs()
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
		for _, ip := range iface {
			if strings.Contains(ip.String(), "192.168.16") {
				ptp = true
			}
		}
		if ptp {
			myGrill.grillIP = "192.168.16.254"
			fmt.Println("Message: PTP to Wifi")
			fmt.Fprintf(&buf, "UH%c%c%s%c%s!", 0, myGrill.ssidlen, myGrill.ssid, myGrill.passwordlen, myGrill.password)
		} else {
			fmt.Println("Need to be connected Ptp to send this message")
			os.Exit(1)
		}
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
		serverKey := []byte(fmt.Sprint(myGrill.serial, myGrill.ExternalIP))
		fmt.Println("Serial:", myGrill.serial)
		fmt.Println("IP:", myGrill.ExternalIP)
		fmt.Println("ServerKey Bytes:", serverKey)
		fmt.Println("ServerKey:", fmt.Sprint(myGrill.serial, myGrill.ExternalIP))

		os.Exit(0) // remove this once key is build correctly
	case grillInfo:
		fmt.Println("Message: Get Grill Temps")
		fmt.Fprint(&buf, "URCV!")
	case externalip:
		fmt.Println("Message: Get External IP")
		fmt.Fprint(&buf, "GMGIP!")
	default:
		fmt.Println("You must choose a message type.")
		os.Exit(1)
	}
	sendData(&buf, myGrill.grillIP)
}

func sendData(b *bytes.Buffer, g string) {
	//b = []byte("UWFM!") // leaveServerMode
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s%s", g, ":8080"), 3*time.Second)
	timeout := time.Now().Add(3 * time.Second)
	conn.SetReadDeadline(timeout)
	if err != nil {
		fmt.Println("Connection Failed")
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

	<meta-data android:name="com.parse.APPLICATION_ID" android:value="0Cc1C1xqnvrsxXkmfNgC9Qu28ejF5iMOePtrrnxA" />
	<meta-data android:name="com.parse.CLIENT_KEY" android:value="iAJbmDULKENYIjliMRBlXd28JKHZaiomU62X1sJG" />

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
