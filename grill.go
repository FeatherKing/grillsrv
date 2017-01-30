package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// gmg support
// UR[2 Byte Grill Temp][2 Byte food probe Temp][2 Byte Target Temp]
// [skip 22 bytes][2 Byte target food probe][1byte on/off/fan][5 byte tail]
/* byte value map
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

// Meat is the id and slice of values
type Meat struct {
	Name string
	//ID     int
	Values []Value
}

// Value is the individual values to be added to meat
type Value struct {
	Time time.Time
	Temp int
}

// HistoryItem is an individual item
type HistoryItem struct {
	Name string
	ID   int
}

// AllItems is a slice of meats
type AllItems struct {
	Meats []Meat
}

//UR001!
func getInfo() ([]byte, error) {
	var buf bytes.Buffer
	fmt.Printf("%s    Request: Get All Info\n", time.Now().Format(time.RFC822))
	//fmt.Println("Request: Get All Info")
	fmt.Fprint(&buf, "UR001!")
	grillResponse, err := sendData(&buf)
	if err != nil {
		return nil, err
	}
	return grillResponse, nil
}

// UT###!
func setGrillTemp(temp int) ([]byte, error) {
	fmt.Printf("%s    Request: Set Grill Temp\n", time.Now().Format(time.RFC822))
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
	fmt.Printf("%s    Request: Set Probe Temp\n", time.Now().Format(time.RFC822))
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
	fmt.Printf("%s    Request: Turn Grill On\n", time.Now().Format(time.RFC822))
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
	fmt.Printf("%s    Request: Turn Grill Off\n", time.Now().Format(time.RFC822))
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
	fmt.Printf("%s    Request: Get Grill ID\n", time.Now().Format(time.RFC822))
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
	fmt.Printf("%s    Request: Get Grill FW\n", time.Now().Format(time.RFC822))
	fmt.Fprint(&buf, "UN!")
	grillResponse, err := sendData(&buf)
	if err != nil {
		return nil, err
	}
	return grillResponse, nil
}

//UH%c%c%s%c%s!
func btoc(ssid string, password string) ([]byte, error) {
	var buf bytes.Buffer
	fmt.Printf("%s    Request: Broadcast to Client Mode\n", time.Now().Format(time.RFC822))
	fmt.Fprintf(&buf, "UH%c%c%s%c%s!", 0, len(ssid), ssid, len(password), password)
	grillResponse, err := sendData(&buf)
	if err != nil {
		return nil, err
	}
	return grillResponse, nil
}

func sendData(b *bytes.Buffer) ([]byte, error) {
	barray := make([]byte, 1024)
	var err error
	var readBytes int
	retries := 5 // the grill doesnt always respond on the first try
	for i := 1; i <= retries; i++ {
		var conn net.Conn
		if i != 1 {
			fmt.Printf("Request Attempt %v\n", i)
		}
		if b.Len() == 0 && i == retries {
			return nil, errors.New("Nothing to Send to Grill")
		}
		// TODO this call using myGrill.GrillIP
		// forces this library to depend on main.go
		// this should be handled differently
		// we might be flipping from broadcast to client mode, this is a udp send
		if b.Len() > 6 {
			conn, err = net.DialTimeout("udp", fmt.Sprintf("%s", myGrill.GrillIP), 3*time.Second)
		} else {
			conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s", myGrill.GrillIP), 3*time.Second)
		}
		if err != nil {
			// TODO make this better
			fmt.Println("Error Connecting to Grill")
			time.Sleep(time.Second * 5)
			if i == retries {
				return nil, errors.New("1")
			}
			continue
		}
		if err != nil && i == retries {
			return nil, errors.New("Connection to Grill Failed")
		}
		timeout := time.Now().Add(time.Second)
		conn.SetReadDeadline(timeout) // sometimes the grill holds the conection forever
		//fmt.Println("Connected")

		defer conn.Close()
		//fmt.Println("Sending Data..")
		_, err = conn.Write(b.Bytes())
		if err != nil && i == retries {
			return nil, errors.New("Failure Sending Payload to Grill")
		}
		//fmt.Printf("Bytes Written: %v\n", ret)
		//b.Reset()

		//fmt.Println("Reading Data..")
		readBytes, err = bufio.NewReader(conn).Read(barray)
		if err != nil && i == retries {
			return nil, errors.New("Failed Reading Result From Grill")
		}
		if readBytes > 0 {
			break
		}
		// trim null of 1024 byte array
		//barray = bytes.Trim(barray, "\x00")

		// print what we got back
		/*
			fmt.Println(string(b.Bytes()))
			fmt.Println(string(barray))
			fmt.Println(barray)
			fmt.Println("Bytes Read:", status)
			fmt.Println("Read Buffer Size:", len(barray))
		*/
	}
	barray = barray[:36]
	return barray, nil
}

/////////////
// BEGIN database funcs
////////////

// TODO needed?
/*
func dbExists() (bool, error) {
	var err error
	db, err = sql.Open("sqlite3", "./grill.db")
	if err != nil {
		return false, errors.New("Error Opening Database")
	}
	defer db.Close()

	return true, nil
}
*/

// TODO create the tables
func createDB() error {
	db, err := sql.Open("sqlite3", "./grill.db")
	if err != nil {
		return errors.New("Error Connecting to Database")
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return errors.New("Error Enabling Foreign Keys")
	}
	sqlStmt := `
	CREATE TABLE if not exists item (
    id integer NOT NULL PRIMARY KEY,
    food text,
    weight real,
    starttime datetime,
    endtime datetime
		);
	CREATE TABLE if not exists log (
    item integer,
    logtime datetime,
    foodtemp integer,
    grilltemp integer,
		FOREIGN KEY (item) REFERENCES item(id)
		);
		`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return errors.New("Error Creating Database")
	}
	defer db.Close()
	return nil
}

func history(id int) (Meat, error) {
	var m Meat
	db, err := sql.Open("sqlite3", "./grill.db")
	if err != nil {
		return m, errors.New("Error Connecting to Database")
	}
	defer db.Close()
	var rows *sql.Rows
	var qerr error
	// get all items
	if id == 0 {
		rows, qerr = db.Query("SELECT id,foodtemp FROM item,log WHERE item.id = log.item")
		if qerr != nil {
			return m, errors.New("Query Failed")
		}
		// get single item
	} else {
		// get name from id
		name := db.QueryRow(`SELECT food
			FROM item
			WHERE id = $1`, id)
		if qerr != nil {
			return m, errors.New("Query Failed")
		}
		name.Scan(&m.Name)

		// get entries for id
		rows, qerr = db.Query(`SELECT logtime,foodtemp
			FROM log
			WHERE log.item = $1
			ORDER BY logtime`, id)
		if qerr != nil {
			return m, errors.New("Query Failed")
		}
		for rows.Next() {
			var v Value
			err = rows.Scan(&v.Time, &v.Temp)
			if err != nil {
				return m, errors.New("Scan Failed")
			}
			m.Values = append(m.Values, v)
		}
	}
	return m, nil
}

func log(f *food) error {
	fmt.Println("Message: Start Logging Food")

	// power on grill read OK from grill
	// check if grill is already on, it might be on from preheating
	/*
		checkPower, err := getInfo()
		if checkPower[grillState] != 1 {
			_, err = powerOn()
			if err != nil {
				return err
			}
		}
	*/
	// connect to persistent storage
	db, err := sql.Open("sqlite3", "./grill.db")
	if err != nil {
		return errors.New("Error Connecting to Database")
	}
	// kick off a go routine to log on interval
	go writeTemp(f, db)
	// inform client that logging was started
	return nil
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
			fmt.Printf("UPDATE item set endtime = '%s' where id = '%v'\n", time.Now().UTC().Format(time.RFC3339), lastInsertID)
			query := "UPDATE item SET endtime = $1 where id = $2"
			stmt, _ := db.Prepare(query)
			defer stmt.Close()
			_, nerr := stmt.Exec(time.Now().UTC().Format(time.RFC3339), lastInsertID)
			if err != nil {
				fmt.Println(nerr.Error())
			}
			break
		}

		// check for high temperature
		var grillInsert int
		var probeInsert int
		if grillResponse[grillTempHigh] == 1 {
			grillInsert = int(grillResponse[grillTemp]) + 256
		} else {
			grillInsert = int(grillResponse[grillTemp])
		}
		if grillResponse[probeTempHigh] == 1 {
			probeInsert = int(grillResponse[probeTemp]) + 256
		} else {
			probeInsert = int(grillResponse[probeTemp])
		}

		fmt.Printf("INSERT INTO log(item,logtime,foodtemp,grilltemp) VALUES('%v','%s','%v' '%v')\n",
			lastInsertID, time.Now().UTC().Format(time.RFC3339), probeInsert, grillInsert)
		query := `INSERT INTO log(item,logtime,foodtemp,grilltemp)
		VALUES($1,$2,$3,$4)`
		stmt, _ := db.Prepare(query)
		defer stmt.Close()
		result, err := stmt.Exec(lastInsertID, time.Now().UTC().Format(time.RFC3339),
			probeInsert, grillInsert)
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

// TODO this should probably return an error too
func historyItems() []HistoryItem {
	var hList []HistoryItem
	db, err := sql.Open("sqlite3", "./grill.db")
	if err != nil {
		return hList
	}
	defer db.Close()
	rows, err := db.Query(`SELECT id,food
			FROM item`)
	if err != nil {
		return hList
	}
	for rows.Next() {
		var h HistoryItem
		err = rows.Scan(&h.ID, &h.Name)
		if err != nil {
			return hList
		}
		hList = append(hList, h)
	}
	return hList
}
