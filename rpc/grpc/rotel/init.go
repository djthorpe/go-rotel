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
	rotel "github.com/djthorpe/rotel"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Service "rpc/rotel:service"
	gopi.RegisterModule(gopi.Module{
		Name:     "rpc/rotel:service",
		Type:     gopi.MODULE_TYPE_SERVICE,
		Requires: []string{"rpc/server", "mutablehome/rotel"},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(Service{
				Server: app.ModuleInstance("rpc/server").(gopi.RPCServer),
				Rotel:  app.ModuleInstance("mutablehome/rotel").(rotel.Rotel),
			}, app.Logger)
		},
	})

	// Stub "rpc/rotel:client"
	gopi.RegisterModule(gopi.Module{
		Name:     "rpc/rotel:client",
		Type:     gopi.MODULE_TYPE_CLIENT,
		Requires: []string{"rpc/clientpool"},
		Run: func(app *gopi.AppInstance, _ gopi.Driver) error {
			if clientpool := app.ModuleInstance("rpc/clientpool").(gopi.RPCClientPool); clientpool == nil {
				return gopi.ErrAppError
			} else {
				clientpool.RegisterClient("gopi.Rotel", NewRotelClient)
				return nil
			}
		},
	})

}
