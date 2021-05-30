// +build !digispark,!arduino,!qtpy

package main

import "machine"

// Replace neo and led in the code below to match the pin
// that you are using if different.
var neo machine.Pin = machine.NEOPIXELS
var led = machine.LED
