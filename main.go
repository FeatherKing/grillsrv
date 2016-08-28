package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

// gmg support
// UR[2 Byte Grill Temp][2 Byte food probe Temp][2 Byte Target Temp]
// [skip 22 bytes][2 Byte target food probe][1byte on/off/fan][5 byte tail]
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
	databaseHost     = "DB_IP"
	databaseName     = "DB_NAME"
	databaseUser     = "DB_USER"
	databasePassword = "DB_PASS"
	databasePort     = 5432
)

type payload struct {
	Cmd    string `json:"cmd"`
	Params string `json:"params"`
}

type food struct {
	Food     string  `json:"food"`
	Weight   float64 `json:"weight"`
	Interval int     `json:"interval"`
}

type temperature struct {
	Grill int `json:"grill"`
	Probe int `json:"probe"`
}

type grill struct {
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
var myGrill = grill{
	grillIP: "LAN_IP:PORT",
	//grillIP:  "FQDN:PORT",
	serial:   "GMGSERIAL",
	ssid:     "SSID",
	password: "WIFI_PASS",
	serverip: "52.26.201.234",
	port:     "8060",
}

func main() {
	// TODO make these read from a file or runtime flags?
	myGrill.ssidlen = len(myGrill.ssid)
	myGrill.passwordlen = len(myGrill.password)
	myGrill.serveriplen = len(myGrill.serverip)
	myGrill.portlen = len(myGrill.port)

	http.HandleFunc("/temp", allTemp)                // all temps GET UR001!
	http.HandleFunc("/temp/grill", singleTemp)       // grill temp GET
	http.HandleFunc("/temp/probe", singleTemp)       // probe temp GET
	http.HandleFunc("/temp/grilltarget", singleTemp) // grill target temp GET/POST UT###!
	http.HandleFunc("/temp/probetarget", singleTemp) // probe target temp GET/POST UF###!
	http.HandleFunc("/power", powerSrv)              // power POST on/off UK001!/UK004!
	http.HandleFunc("/id", idSrv)                    // grill id GET UL!
	http.HandleFunc("/info", infoSrv)                // all fields GET UR001!
	http.HandleFunc("/firmware", fwSrv)              // firmware GET UN!
	http.HandleFunc("/log", log)                     // start grill and log POST
	http.HandleFunc("/cmd", cmd)                     // cmd POST

	http.ListenAndServe(":8000", nil)
}

func singleTemp(w http.ResponseWriter, req *http.Request) {
	requestedTemp := req.URL.Path[6:]
	if req.Method == "GET" {
		grillResponse, err := getInfo()
		if err != nil {
			http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
			return
		}
		var writebuf bytes.Buffer
		fmt.Fprint(&writebuf, "{ ")
		switch requestedTemp {
		case "grill":
			fmt.Fprintf(&writebuf, "\"grilltemp\" : %v ", grillResponse[grillTemp])
		case "grilltarget":
			fmt.Fprintf(&writebuf, "\"grilltarget\" : %v ", grillResponse[grillSetTemp])
		case "probe":
			fmt.Fprintf(&writebuf, "\"probetemp\" : %v ", grillResponse[probeTemp])
		case "probetarget":
			fmt.Fprintf(&writebuf, "\"probetarget\" : %v ", grillResponse[probeSetTemp])
		}
		fmt.Fprint(&writebuf, " }")
		w.Write(writebuf.Bytes())
	} else if req.Method == "POST" {
		defer req.Body.Close()
		var t = temperature{Grill: -1, Probe: -1}
		err := json.NewDecoder(req.Body).Decode(&t)
		if err != nil {
			http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 400)
			return
		}
		switch requestedTemp {
		case "grilltarget":
			if t.Grill == -1 {
				http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", "Grill Target Not Set"), 400)
				return
			}
			grillResponse, err := setGrillTemp(t.Grill)
			if err != nil {
				http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
				return
			}
			w.Write(bytes.Trim(grillResponse, "\x00"))
		case "probetarget":
			if t.Probe == -1 {
				http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", "Probe Target Not Set"), 400)
				return
			}
			grillResponse, err := setProbeTemp(t.Probe)
			if err != nil {
				http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
				return
			}
			w.Write(bytes.Trim(grillResponse, "\x00"))
		}
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), 405)
		return
	}
}

func allTemp(w http.ResponseWriter, req *http.Request) {
	grillResponse, err := getInfo()
	if err != nil {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
		return
	}
	var writebuf bytes.Buffer
	fmt.Fprint(&writebuf, "{ ")
	fmt.Fprintf(&writebuf, "\"grilltemp\" : %v , ", grillResponse[grillTemp])
	fmt.Fprintf(&writebuf, "\"grilltarget\" : %v , ", grillResponse[grillSetTemp])
	fmt.Fprintf(&writebuf, "\"probetemp\" : %v , ", grillResponse[probeTemp])
	fmt.Fprintf(&writebuf, "\"probetarget\" : %v", grillResponse[probeSetTemp])
	fmt.Fprint(&writebuf, " }")
	w.Write(writebuf.Bytes())
}

func log(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	var f food
	fmt.Println("Message: Start Logging Food")

	// check for post method validate request
	if req.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), 405)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&f)

	//validate request
	if f.Food == "" || f.Weight == 0 {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", "Missing Required Values"), 400)
		return
	}
	if f.Interval == 0 {
		f.Interval = 5
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 400)
		return
	}

	// power on grill read OK from grill
	// check if grill is already on, it might be on from preheating
	checkPower, err := getInfo()
	if checkPower[grillState] != 1 {
		_, err = powerOn()
		if err != nil {
			http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
			return
		}
	}
	// connect to persistent storage
	url := fmt.Sprintf("DB_USER://%s:%s@%s:%d/%s?sslmode=disable",
		databaseUser, databasePassword, databaseHost, databasePort, databaseName)
	db, err := sql.Open("DB_USER", url)
	if err != nil {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", "Error Connecting to Database"), 500)
		return
	}
	// kick off a go routine to log on interval
	go writeTemp(&f, db)
	// inform client that logging was started
	fmt.Fprint(&buf, "{\"logging_started\": true }")
	w.Write(buf.Bytes())
}

func writeTemp(f *food, db *sql.DB) error {
	var lastInsertID int
	startTime := time.Now().UTC()
	defer db.Close()
	fmt.Printf("INSERT INTO item(food,weight,starttime) VALUES('%s','%v','%s') returning id;\n",
		f.Food, f.Weight, startTime.Format(time.RFC3339))
	query := `INSERT INTO item(food,weight,starttime)
	VALUES($1,$2,$3)
	RETURNING id`
	stmt, qerr := db.Prepare(query)
	if qerr != nil {
		fmt.Println(qerr.Error())
	}
	defer stmt.Close()
	qerr = stmt.QueryRow(f.Food, f.Weight, startTime.Format(time.RFC3339)).Scan(&lastInsertID)
	if qerr != nil {
		fmt.Println(qerr.Error())
	}

	// loop on interval
	// the loop will not end on failure to read from grill
	// it will only end if the grill gets turned off
	for {
		time.Sleep(time.Minute * time.Duration(f.Interval))
		// get current temps
		grillResponse, err := getInfo()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// stop if the grill turned off or in fan mode
		if grillResponse[grillState] == 0 || grillResponse[grillState] == 2 {
			fmt.Printf("UPDATE item set endtime = '%s' where id = '%v'\n", time.Now().Format(time.RFC3339), lastInsertID)
			query := "UPDATE item SET endtime = $1 where id = $2"
			stmt, _ := db.Prepare(query)
			defer stmt.Close()
			_, nerr := stmt.Exec(time.Now().Format(time.RFC3339), lastInsertID)
			if err != nil {
				fmt.Println(nerr.Error())
			}
			break
		}

		// insert temp
		fmt.Printf("INSERT INTO log(item,logtime,foodtemp,grilltemp) VALUES('%v','%s','%v' '%v')\n",
			lastInsertID, time.Now().Format(time.RFC3339), grillResponse[probeTemp], grillResponse[grillTemp])
		query := `INSERT INTO log(item,logtime,foodtemp,grilltemp)
		VALUES($1,$2,$3,$4)`
		stmt, _ := db.Prepare(query)
		defer stmt.Close()
		result, err := stmt.Exec(lastInsertID, time.Now().Format(time.RFC3339),
			grillResponse[probeTemp], grillResponse[grillTemp])
		if err != nil {
			fmt.Println(err.Error())
		}
		result.RowsAffected()
	}
	//   write to storage on interval
	//   if interval not specified, default 5 mins
	//   if grill read fails, try again, then move on?
	return nil
}

func idSrv(w http.ResponseWriter, req *http.Request) {
	grillResponse, err := grillID()
	if err != nil {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
		return
	}
	w.Write(bytes.Trim(grillResponse, "\x00"))
}

func fwSrv(w http.ResponseWriter, req *http.Request) {
	grillResponse, err := grillFW()
	if err != nil {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
		return
	}
	w.Write(bytes.Trim(grillResponse, "\x00"))
}

func cmd(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	fmt.Println("Message: Execute Command")
	// change broadcast to client
	// ptp to Wifi
	// server mode
	if req.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), 405)
		return
	}
	defer req.Body.Close()
	var pay payload
	err := json.NewDecoder(req.Body).Decode(&pay)
	fmt.Printf("Decoded Request: %s %s %s\n", &pay, pay.Cmd, pay.Params)
	if err != nil {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 400)
		return
	}
	switch pay.Cmd {
	case "btoc": // flip grill from broadcast to client mode
		fmt.Fprintf(&buf, "UH%c%c%s%c%s!", 0, myGrill.ssidlen, myGrill.ssid, myGrill.passwordlen, myGrill.password)
	}
	grillResponse, err := sendData(&buf)
	if err != nil {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
		return
	}
	w.Write(bytes.Trim(grillResponse, "\x00"))
}

func infoSrv(w http.ResponseWriter, req *http.Request) {
	grillResponse, err := getInfo()
	if err != nil {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
		return
	}
	var writebuf bytes.Buffer
	// gmg support
	// UR[2 Byte Grill Temp][2 Byte food probe Temp][2 Byte Target Temp][skip 22 bytes][2 Byte target food probe][1byte on/off/fan][5 byte tail]
	fmt.Fprint(&writebuf, "{ ")
	fmt.Fprintf(&writebuf, "\"grilltemp\" : %v , ", grillResponse[grillTemp])
	fmt.Fprintf(&writebuf, "\"grilltarget\" : %v , ", grillResponse[grillSetTemp])
	fmt.Fprintf(&writebuf, "\"probetemp\" : %v , ", grillResponse[probeTemp])
	fmt.Fprintf(&writebuf, "\"probetarget\" : %v ,", grillResponse[probeSetTemp])
	fmt.Fprintf(&writebuf, "\"curveremaintime\" : %v ,", grillResponse[curveRemainTime])
	fmt.Fprintf(&writebuf, "\"warncode\" : \"%s\" ,", warnStates[int(grillResponse[warnCode])])
	fmt.Fprintf(&writebuf, "\"grillstate\" : \"%s\" ,", grillStates[int(grillResponse[grillState])])
	fmt.Fprintf(&writebuf, "\"firestate\" : \"%s\" ,", fireStates[int(grillResponse[fireState])])
	fmt.Fprintf(&writebuf, "\"filestatepercent\" : %v, ", grillResponse[fileStatePercent])
	fmt.Fprintf(&writebuf, "\"profileend\" : %v, ", grillResponse[profileEnd])
	fmt.Fprintf(&writebuf, "\"grilltype\" : %v", grillResponse[grillType])
	fmt.Fprint(&writebuf, " }")
	w.Write(writebuf.Bytes())
}

func powerSrv(w http.ResponseWriter, req *http.Request) {
	var grillResponse []byte
	var err error
	var pay payload
	if req.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), 405)
		return
	}
	defer req.Body.Close()
	err = json.NewDecoder(req.Body).Decode(&pay)
	fmt.Printf("Decoded Request: %s %s %s\n", &pay, pay.Cmd, pay.Params)
	if err != nil {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 400)
		return
	}
	switch pay.Cmd {
	case "on":
		grillResponse, err = powerOn()
	case "off":
		grillResponse, err = powerOff()
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
		return
	}
	w.Write(bytes.Trim(grillResponse, "\x00"))
}

/*
* This stuff is old, it could be used to switch the grill from ptp to wifi
	http.HandleFunc("/",
		func(w http.ResponseWriter, req *http.Request) {
			requestedFile := req.URL.Path[1:]
			switch requestedFile {
			case "usewifi":
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
				}
			case "servermode":
				fmt.Println("Message: Wifi to Server Mode")
				fmt.Fprintf(&buf, "UG%c%s%c%s%c%s%c%s!", myGrill.ssidlen, myGrill.ssid, myGrill.passwordlen, myGrill.password, myGrill.serveriplen, myGrill.serverip, myGrill.portlen, myGrill.port)
			case "serverkey":
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
			case "grillinfo":
				fmt.Println("Message: Get Grill Temps?")
				fmt.Fprint(&buf, "URCV!")
			case "externalip":
				fmt.Println("Message: Get External IP")
				fmt.Fprint(&buf, "GMGIP!")
			default:
				w.WriteHeader(404)
			}
		})
*/
