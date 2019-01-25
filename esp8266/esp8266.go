// Package esp8266 implements TCP/UDP communication over serial
// with a separate Wifi ESP8266 board using the Espressif AT command set
// across a UART interface.
//
// In order to use this driver, the ESP8266 must be flashed with firmware
// supporting the AT command set. Many ESP8266 chips already have this firmware
// installed by default. You will need to install this firmware if you have an
// ESP8266 that has been flashed with NodeMCU (Lua) or Arduino firmware.
//
// Datasheet:
// https://www.espressif.com/sites/default/files/documentation/0a-esp8266ex_datasheet_en.pdf
//
// AT command set:
// https://www.espressif.com/sites/default/files/documentation/4a-esp8266_at_instruction_set_en.pdf
//
package esp8266

import (
	"machine"
	"strconv"
	"strings"
	"time"
)

// Device wraps UART connection to the ESP8266.
type Device struct {
	bus machine.UART

	// command responses that come back from the ESP8266
	response []byte

	// data received from a TCP/UDP connection forwarded by the ESP8266
	socketdata    []byte
	socketdataLen int
}

// New returns a new esp8266-wifi driver. Pass in a fully configured UART bus.
func New(b machine.UART) Device {
	return Device{bus: b, response: make([]byte, 512), socketdata: make([]byte, 1024)}
}

// Configure sets up the device for communication.
func (d Device) Configure() {
}

// Connected checks if there is communication with the ESP8266.
func (d *Device) Connected() bool {
	d.Execute(Test)

	// handle response here, should include "OK"
	r := d.Response()
	if strings.Contains(string(r), "OK") {
		return true
	}
	return false
}

// Write raw bytes to the UART.
func (d *Device) Write(b []byte) (n int, err error) {
	return d.bus.Write(b)
}

// Read raw bytes from the UART.
func (d *Device) Read(b []byte) (n int, err error) {
	return d.bus.Read(b)
}

// ReadSocket returns the data that has already been read in from the responses.
func (d *Device) ReadSocket(b []byte) (n int, err error) {
	// make sure no data in buffer
	d.Response()

	count := len(b)
	if len(b) > d.socketdataLen {
		count = d.socketdataLen
	}

	for i := 0; i < count; i++ {
		b[i] = d.socketdata[i]
	}

	d.socketdataLen = 0
	return count, nil
}

// Response gets the next response bytes from the ESP8266.
func (d *Device) Response() []byte {
	var i, retries int

	header := make([]byte, 2)
	for {
		for d.bus.Buffered() > 0 {
			// get the first 2 bytes
			header[0], _ = d.bus.ReadByte()
			header[1], _ = d.bus.ReadByte()

			if d.isLeadingCRLF(header) {
				// skip it
				header[0], _ = d.bus.ReadByte()
				header[1], _ = d.bus.ReadByte()
			}

			if d.isIPD(header) {
				// is socket data packet
				d.parseIPD()
			} else {
				// no, so put into response
				d.response[i] = header[0]
				i++
				d.response[i] = header[1]
				i++
			}

			// read the rest of normal command response
			for d.bus.Buffered() > 0 {
				data, _ := d.bus.ReadByte()
				d.response[i] = data
				i++
			}
		}
		retries++
		if retries > 2 {
			break
		}

		// pause to make sure is no more data to be read
		time.Sleep(10 * time.Millisecond)
	}
	return d.response[:i]
}

func (d *Device) isLeadingCRLF(b []byte) bool {
	if len(b) < 2 {
		return false
	}
	if b[0] == 13 && b[1] == 10 {
		return true
	}
	return false
}

func (d *Device) isIPD(b []byte) bool {
	if len(b) < 2 {
		return false
	}
	if b[0] == '+' && b[1] == 'I' {
		return true
	}
	return false
}

func (d *Device) parseIPD() bool {
	data, _ := d.bus.ReadByte()
	if data != 'P' {
		// error
		return false
	}
	data, _ = d.bus.ReadByte()
	if data != 'D' {
		// error
		return false
	}
	data, _ = d.bus.ReadByte()
	if data != ',' {
		// error
		return false
	}

	// get the expected data length
	// skip remaining header up to the ":"
	buf := []byte{}
	data, _ = d.bus.ReadByte()
	for data != ':' {
		// put into the buffer with int value here
		buf = append(buf, data)

		// read next value
		data, _ = d.bus.ReadByte()
	}

	val := string(buf)
	count, err := strconv.Atoi(val)
	if err != nil {
		// not expected data here. what to do?
		return false
	}

	// load up the socket data
	// only read the expected amount of data
	for m := 0; m < count; m++ {
		data, _ = d.bus.ReadByte()
		d.socketdata[d.socketdataLen] = data
		d.socketdataLen++
	}

	return true
}
