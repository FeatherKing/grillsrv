package main

import (
	"bytes"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type GrillRespLength struct {
	command    []byte
	resultSize int
}

var tests = []GrillRespLength{
	{[]byte("UR001!"), 36},
	{[]byte("UL!"), 36},
	{[]byte("UN!"), 36},
}

// checks that the return value is 36 bytes long
func TestSendDataLength(t *testing.T) {
	loadConfig()
	for _, item := range tests {
		var buf bytes.Buffer
		fmt.Fprint(&buf, "UR001!")
		v, err := sendData(&buf)
		if err != nil {
			t.Error("Error", err.Error())
		}
		if len(v) != item.resultSize {
			t.Error(
				"Command Sent", item.command,
				"Expected Return Size", item.resultSize,
				"Received", len(v),
			)
		}
	}
}

// TODO test network Connection
// TODO test an entire 'getinfo' request and inspect the return to be valid
// TODO test loading the config, config could add new parameters over time
