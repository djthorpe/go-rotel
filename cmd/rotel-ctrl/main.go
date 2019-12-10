/*
	Rotel RS232 Control
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package main

import (
	"fmt"
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	rotel "github.com/djthorpe/rotel"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	// Wait for CTRL+C
	app.Logger.Info("Waiting for CTRL+C")
	app.WaitForSignal()

	// Success
	done <- gopi.DONE
	return nil
}

func EventLoop(app *gopi.AppInstance, start chan<- struct{}, done <-chan struct{}) error {
	if device, ok := app.ModuleInstance("mutablehome/rotel").(rotel.Rotel); ok == false {
		return fmt.Errorf("Invalid Rotel module")
	} else {
		evt := device.Subscribe()
		start <- gopi.DONE
	FOR_LOOP:
		for {
			select {
			case <-done:
				break FOR_LOOP
			case event := <-evt:
				fmt.Println(event)
			}
		}
		device.Unsubscribe(evt)
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("mutablehome/rotel")

	// Run the server and register all the services
	os.Exit(gopi.CommandLineTool2(config, Main, EventLoop))
}
