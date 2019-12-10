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
	"strings"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	rotel "github.com/djthorpe/rotel"
)

////////////////////////////////////////////////////////////////////////////////

func SendCommand(device rotel.Rotel, cmd string) error {
	values := make([]string, 0, rotel.ROTEL_SOURCE_MAX)
	for v := rotel.ROTEL_COMMAND_NONE + 1; v <= rotel.ROTEL_COMMAND_MAX; v++ {
		str := strings.ToLower(strings.TrimPrefix(fmt.Sprint(v), "ROTEL_COMMAND_"))
		if strings.ToLower(cmd) == str {
			return device.Send(v)
		} else {
			values = append(values, str)
		}
	}
	return fmt.Errorf("command should be one of: %v", strings.Join(values, ", "))
}

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	if device, ok := app.ModuleInstance("mutablehome/rotel").(rotel.Rotel); ok == false {
		return fmt.Errorf("Invalid Rotel module")
	} else {
		for _, arg := range app.AppFlags.Args() {
			if err := SendCommand(device, arg); err != nil {
				return err
			}
		}
	}

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
