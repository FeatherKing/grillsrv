package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {
	ssid := "SSID"
	password := "WIFI_PASS"
	//serverip := "52.26.201.234"
	//port := "8060"
	ssidlen := len(ssid)
	passwordlen := len(password)
	//serveriplen := len(serverip)
	//portlen := len(port)
	s := fmt.Sprintf("UH%c%c%s%c%s!", 0, ssidlen, ssid, passwordlen, password) // connect to wifi
	//s := fmt.Sprintf("UG%c%s%c%s%c%s%c%s!", ssidlen, ssid, passwordlen, password, serveriplen, serverip, portlen, port)
	bytes := []byte(s)
	sendData(bytes)
}

func sendData(b []byte) {
	//b = []byte("UWFM!") // leaveServerMode
	//grillIdCommand := []byte("UL!")
	//str := "UL!" // get grill id
	//str := "GET / HTTP/1.0\r\n\r\n"
	conn, err := net.Dial("tcp", "LAN_IP:PORT")
	//conn, err := net.Dial("tcp", "google.com:80")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Connected")
	timeout := time.Now().Add(time.Second)
	conn.SetReadDeadline(timeout)

	defer conn.Close()
	fmt.Println("Sending Data")
	//ret, _ := conn.Write([]byte(str))
	ret, _ := conn.Write(b)

	//ret, err := fmt.Fprint(conn, buf)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("Bytes Written: %v\n", ret)
	fmt.Println("Reading Data")
	//read := make([]byte, 2048)
	//status, _ := conn.Read(read)
	//status, err := bufio.NewReader(conn).ReadBytes(33)
	status, err := bufio.NewReader(conn).ReadString(33)

	fmt.Println(status)
	fmt.Println("total size:", len(status))

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

		         type myStruct struct {
		         	ssidlen int
		         	ssid string
		         	passwordlen int
		         	password string
		         	serveriplen int
		         	serverip string
		         	portlen int
		         	port string
		         	//ssidlen, ssid, passwordlen, password, serveriplen, serverip, portlen, port)
		         }

		         func main() {
		         	ssid := "SSID"
		         	password := "WIFI_PASS"
		         	serverip := "52.26.201.234"
		         	port := "8060"
		         	var bin_buf bytes.Buffer

		         	x := myStruct{len(ssid),"SSID", len(password), "WIFI_PASS", len(serverip), "52.26.201.234", len(port), "8060"}
		         	binary.Write(&bin_buf, binary.BigEndian, x)
		         	//fmt.Printf("% x", sha1.Sum(bin_buf.Bytes()))
		         	fmt.Printf("UG%!", bin_buf.Bytes())
		         	fmt.Printf(x)
		         }


	*/

}
