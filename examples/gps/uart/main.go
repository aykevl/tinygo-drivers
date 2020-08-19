package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/gps"
)

func main() {
	println("GPS UART Example")
	machine.UART1.Configure(machine.UARTConfig{BaudRate: 9600})
	ublox := gps.NewUART(&machine.UART1)
	parser := gps.NewParser()
	var fix gps.Fix
	for {
		s, err := ublox.NextSentence()
		if err != nil {
			println(err)
			continue
		}

		fix, err = parser.Parse(s)
		if err != nil {
			println(err)
			continue
		}
		if fix.Valid {
			print(fix.Time.Format("15:04:05"))
			print(", lat=")
			print(fix.Latitude)
			print(", long=")
			print(fix.Longitude)
			print(", altitude:=", fix.Altitude)
			print(", satellites=", fix.Satellites)
			println()
		} else {
			println("No fix")
		}
		time.Sleep(200 * time.Millisecond)
	}
}
