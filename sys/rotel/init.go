/*
	Rotel RS232 Control
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package rotel

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	gopi.RegisterModule(gopi.Module{
		Name: "mutablehome/rotel",
		Type: gopi.MODULE_TYPE_OTHER,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagString("rotel.tty", "/dev/ttyUSB0", "RS232 device")
			config.AppFlags.FlagUint("rotel.baudrate",BAUD_RATE_DEFAULT, "RS232 speed")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			tty, _ := app.AppFlags.GetString("rotel.tty")
			baudrate,  _ := app.AppFlags.GetUint("rotel.baudrate")
			return gopi.Open(Rotel{
				TTY:     tty,
				BaudRate: baudrate,
			}, app.Logger)
		},
	})
}