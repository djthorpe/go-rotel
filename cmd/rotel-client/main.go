package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	rpc "github.com/djthorpe/gopi-rpc"
	rotel "github.com/djthorpe/rotel"
)

////////////////////////////////////////////////////////////////////////////////

func SetPower(stub rotel.RotelClient, value string) error {
	switch value {
	case "on":
		return stub.Set(rotel.RotelState{
			Power: rotel.ROTEL_POWER_ON,
		})
	case "off", "standby":
		return stub.Set(rotel.RotelState{
			Power: rotel.ROTEL_POWER_STANDBY,
		})
	case "toggle":
		return stub.Send(rotel.ROTEL_COMMAND_POWER_TOGGLE)
	default:
		return fmt.Errorf("-power value should be on, standby, off or toggle")
	}
}

func SetVolume(stub rotel.RotelClient, value rotel.Volume) error {
	if value == 0 || value > rotel.ROTEL_VOLUME_MAX {
		return fmt.Errorf("-volume value should be between 1 and %v", rotel.ROTEL_VOLUME_MAX)
	} else {
		return stub.Set(rotel.RotelState{
			Volume: rotel.Volume(value),
		})
	}
}

func SetSource(stub rotel.RotelClient, value string) error {
	values := make([]string, 0, rotel.ROTEL_SOURCE_MAX)
	for v := rotel.ROTEL_SOURCE_NONE + 1; v < rotel.ROTEL_SOURCE_MAX; v++ {
		str := strings.ToLower(strings.TrimPrefix(fmt.Sprint(v), "ROTEL_SOURCE_"))
		if strings.ToLower(value) == str {
			return stub.Set(rotel.RotelState{
				Source: rotel.Source(v),
			})
		} else {
			values = append(values, str)
		}
	}
	return fmt.Errorf("-source value should be one of: %v", strings.Join(values, ", "))
}

func SendCommand(stub rotel.RotelClient, value string) error {
	values := make([]string, 0, rotel.ROTEL_SOURCE_MAX)
	for v := rotel.ROTEL_COMMAND_NONE + 1; v <= rotel.ROTEL_COMMAND_MAX; v++ {
		str := strings.ToLower(strings.TrimPrefix(fmt.Sprint(v), "ROTEL_COMMAND_"))
		if strings.ToLower(value) == str {
			return stub.Send(v)
		} else {
			values = append(values, str)
		}
	}
	return fmt.Errorf("command should be one of: %v", strings.Join(values, ", "))
}

func EventLoop(stub gopi.Publisher, done <-chan struct{}) {
	evt := stub.Subscribe()
FOR_LOOP:
	for {
		select {
		case <-done:
			break FOR_LOOP
		case event := <-evt:
			if evt_ := event.(rotel.RotelEvent); evt_ != nil && evt_.Type() != rotel.EVENT_TYPE_NONE {
				fmt.Println(evt_.Type(), evt_.State())
			}
		}
	}
	stub.Unsubscribe(evt)
}

func Main(app *gopi.AppInstance, services []gopi.RPCServiceRecord, done chan<- struct{}) error {
	// Get the client
	if stub, err := app.ClientPool.NewClientEx("gopi.Rotel", services, gopi.RPC_FLAG_SERVICE_ANY); err != nil {
		return err
	} else if device, ok := stub.(rotel.RotelClient); device == nil || ok == false {
		return fmt.Errorf("Invalid rotel client")
	} else if err := device.Ping(); err != nil {
		return err
	} else {
		// Power
		if power, exists := app.AppFlags.GetString("power"); exists {
			if err := SetPower(device, power); err != nil {
				return err
			}
		}
		// Volume
		if volume, exists := app.AppFlags.GetUint("volume"); exists {
			if err := SetVolume(device, rotel.Volume(volume)); err != nil {
				return err
			}
		}
		// Source
		if source, exists := app.AppFlags.GetString("source"); exists {
			if err := SetSource(device, source); err != nil {
				return err
			}
		}

		// Commands
		for _, arg := range app.AppFlags.Args() {
			if err := SendCommand(device, arg); err != nil {
				return err
			}
		}

		// Watch
		ctx, cancel := context.WithCancel(context.Background())
		stop := make(chan struct{})
		go func() {
			if err := device.StreamEvents(ctx); err != nil && err != context.Canceled {
				app.Logger.Error("Error: %v", err)
			}
			stop <- gopi.DONE
		}()

		// Print events in background as the occur
		go EventLoop(device, stop)

		// Wait for CTRL+C then cancel
		app.Logger.Info("Press CTRL+C to exit")
		app.WaitForSignal()
		cancel()
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/rotel:client")

	config.AppFlags.FlagString("power", "", "Power switch (on, off or toggle)")
	config.AppFlags.FlagUint("volume", 55, "Set volume (1-96)")
	config.AppFlags.FlagString("source", "", "Set input source")

	// Run the command line tool
	os.Exit(rpc.Client(config, 200*time.Millisecond, Main))
}
