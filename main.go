package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

// gmg support
// UR[2 Byte Grill Temp][2 Byte food probe Temp][2 Byte Target Temp]
// [skip 22 bytes][2 Byte target food probe][1byte on/off/fan][5 byte tail]
const (
	grillTemp        = 2
	grillTempHigh    = 3
	probeTemp        = 4
	probeTempHigh    = 5
	grillSetTemp     = 6
	grillSetTempHigh = 7
	curveRemainTime  = 20
	warnCode         = 24
	probeSetTemp     = 28
	probeSetTempHigh = 29
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
	VERSION          = 1.3
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
	GrillIP     string
	ExternalIP  string `json:"ip"`
	Serial      string
	Ssid        string
	Password    string
	ssidlen     int
	passwordlen int
	Serverip    string
	Port        string
	serveriplen int
	portlen     int
	ListenPort  string
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

var myGrill = grill{}

func main() {
	// TODO read fromd a file or runtime flags?
	loadConfig()
	router := httprouter.New()
	router.GET("/temp", allTemp)           // all temps GET UR001!
	router.GET("/temp/:name", singleTemp)  // all temps GET UR001!
	router.POST("/temp/:name", singleTemp) // all temps GET UR001!
	router.POST("/power", powerSrv)        // power POST on/off UK001!/UK004!
	router.GET("/id", idSrv)               // grill id GET UL!
	router.GET("/info", infoSrv)           // all fields GET UR001!
	router.GET("/firmware", fwSrv)         // firmware GET UN!
	router.POST("/log", logSrv)            // start grill and log POST
	router.POST("/cmd", cmd)               // cmd POST
	router.GET("/history/:id", historySrv) // history GET

	router.GET("/", index)
	router.ServeFiles("/assets/*filepath", http.Dir("assets"))
	http.ListenAndServe(":"+myGrill.ListenPort, router)
}

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error Opening config.json")
		os.Exit(1)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&myGrill)
	if err != nil {
		fmt.Println("Error Reading config.json")
		os.Exit(1)
	}
	fmt.Printf("GrillIP: %s\nGMGSerial: %s\nSSID: %s\nPASS: %s\n",
		myGrill.GrillIP, myGrill.Serial, myGrill.Ssid, myGrill.Password)
}

func index(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	indexTemplate := template.New("index")
	indexTemplate.ParseFiles("assets/index.html")
	//indexTemplate.Parse(`{{define "T"}}Hello, {{ range . }}Name: {{.Name}}, ID: {{.ID}} {{end}}{{end}}`)
	//_ = indexTemplate.ExecuteTemplate(w, "T", historyItems())
	_ = indexTemplate.ExecuteTemplate(w, "index.html", historyItems())
	//http.ServeFile(w, req, "assets/index.html")
}

func singleTemp(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method == "GET" {
		grillResponse, err := getInfo()
		if err != nil {
			http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
			return
		}
		var writebuf bytes.Buffer
		fmt.Fprint(&writebuf, "{ ")
		switch ps.ByName("name") {
		case "grilltemp":
			if grillResponse[grillTempHigh] == 1 {
				fmt.Fprintf(&writebuf, "\"grilltemp\" : %d ", int(grillResponse[grillTemp])+255)
			} else {
				fmt.Fprintf(&writebuf, "\"grilltemp\" : %v ", grillResponse[grillTemp])
			}
		case "grilltarget":
			if grillResponse[grillSetTempHigh] == 1 {
				fmt.Fprintf(&writebuf, "\"grilltarget\" : %d ", int(grillResponse[grillSetTemp])+255)
			} else {
				fmt.Fprintf(&writebuf, "\"grilltarget\" : %v ", grillResponse[grillSetTemp])
			}
		case "probetemp":
			if grillResponse[probeTempHigh] == 1 {
				fmt.Fprintf(&writebuf, "\"probetemp\" : %d ", int(grillResponse[probeTemp])+255)
			} else {
				fmt.Fprintf(&writebuf, "\"probetemp\" : %v ", grillResponse[probeTemp])
			}
		case "probetarget":
			if grillResponse[probeSetTempHigh] == 1 {
				fmt.Fprintf(&writebuf, "\"probetarget\" : %d ", int(grillResponse[probeSetTemp])+255)
			} else {
				fmt.Fprintf(&writebuf, "\"probetarget\" : %v ", grillResponse[probeSetTemp])
			}
		default:
			http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", http.StatusText(http.StatusBadRequest)), http.StatusBadRequest)
			return
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
		switch ps.ByName("name") {
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

func allTemp(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	grillResponse, err := getInfo()
	if err != nil {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
		return
	}
	var writebuf bytes.Buffer
	fmt.Fprint(&writebuf, "{ ")
	if grillResponse[grillTempHigh] == 1 {
		fmt.Fprintf(&writebuf, "\"grilltemp\" : %d , ", int(grillResponse[grillTemp])+255)
	} else {
		fmt.Fprintf(&writebuf, "\"grilltemp\" : %v , ", grillResponse[grillTemp])
	}
	if grillResponse[grillSetTempHigh] == 1 {
		fmt.Fprintf(&writebuf, "\"grilltarget\" : %d , ", int(grillResponse[grillSetTemp])+255)
	} else {
		fmt.Fprintf(&writebuf, "\"grilltarget\" : %v , ", grillResponse[grillSetTemp])
	}
	if grillResponse[probeTempHigh] == 1 {
		fmt.Fprintf(&writebuf, "\"probetemp\" : %d , ", int(grillResponse[probeTemp])+255)
	} else {
		fmt.Fprintf(&writebuf, "\"probetemp\" : %v , ", grillResponse[probeTemp])
	}
	if grillResponse[probeSetTempHigh] == 1 {
		fmt.Fprintf(&writebuf, "\"probetarget\" : %d ", int(grillResponse[probeSetTemp])+255)
	} else {
		fmt.Fprintf(&writebuf, "\"probetarget\" : %v ", grillResponse[probeSetTemp])
	}
	fmt.Fprint(&writebuf, " }")
	w.Write(writebuf.Bytes())
}

func logSrv(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var buf bytes.Buffer
	var f food
	fmt.Println("Message: Start Logging Food")

	// check for post method validate request
	if req.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), 405)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&f)
	if err != nil {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 400)
		return
	}

	//validate request
	if f.Food == "" || f.Weight == 0 {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", "Missing Required Values"), 400)
		return
	}
	if f.Interval == 0 {
		f.Interval = 5
	}

	err = log(&f)
	if err != nil {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
		return
	}
	fmt.Fprint(&buf, "{\"logging_started\": true }")
	w.Write(buf.Bytes())
}

func idSrv(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	grillResponse, err := grillID()
	if err != nil {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
		return
	}
	w.Write(bytes.Trim(grillResponse, "\x00"))
}

func fwSrv(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	grillResponse, err := grillFW()
	if err != nil {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
		return
	}
	w.Write(bytes.Trim(grillResponse, "\x00"))
}

func cmd(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var grillResponse []byte
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
		grillResponse, err = btoc(myGrill.Ssid, myGrill.Password)
		if err != nil {
			http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
			return
		}
	}

	w.Write(bytes.Trim(grillResponse, "\x00"))
}

func infoSrv(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	grillResponse, err := getInfo()
	if err != nil {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
		return
	}
	var writebuf bytes.Buffer
	// gmg support
	// UR[2 Byte Grill Temp][2 Byte food probe Temp][2 Byte Target Temp][skip 22 bytes][2 Byte target food probe][1byte on/off/fan][5 byte tail]
	fmt.Fprint(&writebuf, "{ ")
	if grillResponse[grillTempHigh] == 1 {
		fmt.Fprintf(&writebuf, "\"grilltemp\" : %d , ", int(grillResponse[grillTemp])+255)
	} else {
		fmt.Fprintf(&writebuf, "\"grilltemp\" : %v , ", grillResponse[grillTemp])
	}
	if grillResponse[grillSetTempHigh] == 1 {
		fmt.Fprintf(&writebuf, "\"grilltarget\" : %d , ", int(grillResponse[grillSetTemp])+255)
	} else {
		fmt.Fprintf(&writebuf, "\"grilltarget\" : %v , ", grillResponse[grillSetTemp])
	}
	if grillResponse[probeTempHigh] == 1 {
		fmt.Fprintf(&writebuf, "\"probetemp\" : %d , ", int(grillResponse[probeTemp])+255)
	} else {
		fmt.Fprintf(&writebuf, "\"probetemp\" : %v , ", grillResponse[probeTemp])
	}
	if grillResponse[probeSetTempHigh] == 1 {
		fmt.Fprintf(&writebuf, "\"probetarget\" : %d , ", int(grillResponse[probeSetTemp])+255)
	} else {
		fmt.Fprintf(&writebuf, "\"probetarget\" : %v , ", grillResponse[probeSetTemp])
	}
	fmt.Fprintf(&writebuf, "\"curveremaintime\" : %v , ", grillResponse[curveRemainTime])
	fmt.Fprintf(&writebuf, "\"warncode\" : \"%s\" , ", warnStates[int(grillResponse[warnCode])])
	fmt.Fprintf(&writebuf, "\"grillstate\" : \"%s\" , ", grillStates[int(grillResponse[grillState])])
	fmt.Fprintf(&writebuf, "\"firestate\" : \"%s\" , ", fireStates[int(grillResponse[fireState])])
	fmt.Fprintf(&writebuf, "\"filestatepercent\" : %v, ", grillResponse[fileStatePercent])
	fmt.Fprintf(&writebuf, "\"profileend\" : %v, ", grillResponse[profileEnd])
	fmt.Fprintf(&writebuf, "\"grilltype\" : %v", grillResponse[grillType])
	fmt.Fprint(&writebuf, " }")
	w.Write(writebuf.Bytes())
}

func powerSrv(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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

func historySrv(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))
	m, err := history(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("{ \"error\": \"%s\" }", err.Error()), 500)
		return
	}

	err = json.NewEncoder(w).Encode(&m)
}

/*
* This stuff is old, it could be used to switch the grill from ptp to wifi
	router.GET("/",
		func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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
