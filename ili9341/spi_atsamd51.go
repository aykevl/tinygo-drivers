// +build atsamd51

package ili9341

import (
	"machine"
)

type spiDriver struct {
	bus machine.SPI
	dc  machine.Pin
	rst machine.Pin
	cs  machine.Pin
	rd  machine.Pin
}

func NewSpi(bus machine.SPI, dc, cs, rst machine.Pin) *Device {
	return &Device{
		dc:  dc,
		cs:  cs,
		rst: rst,
		rd:  machine.NoPin,
		driver: &spiDriver{
			bus: bus,
		},
	}
}

func (pd *spiDriver) configure(config *Config) {
}

//go:inline
func (pd *spiDriver) write8(b byte) {
	pd.bus.Tx([]byte{b}, nil)
}

//go:inline
func (pd *spiDriver) write16(data uint16) {
	pd.bus.Tx([]byte{byte(data >> 8), byte(data)}, nil)
}

//go:inline
func (pd *spiDriver) write16n(data uint16, n int) {
	for i := 0; i < n; i++ {
		pd.write16(data)
	}
}

//go:inline
func (pd *spiDriver) write16sl(data []uint16) {
	for i, c := 0, len(data); i < c; i++ {
		pd.write16(data[i])
	}
}
