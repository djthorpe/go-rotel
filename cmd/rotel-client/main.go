package main

import (
	"fmt"
	"os"
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
			Power: rotel.ROTEL_POWER_STANDY,
		})
	case "", "toggle":
		return stub.Set(rotel.RotelState{
			Power: rotel.ROTEL_POWER_TOGGLE,
		})
	default:
		return fmt.Errorf("-power value should be on, standby or toggle")
	}
}

func Main(app *gopi.AppInstance, services []gopi.RPCServiceRecord, done chan<- struct{}) error {
	// Get the client
	if stub, err := app.ClientPool.NewClientEx("gopi.Rotel", services, gopi.RPC_FLAG_SERVICE_ANY); err != nil {
		return err
	} else if rotel, ok := stub.(rotel.RotelClient); rotel == nil || ok == false {
		return fmt.Errorf("Invalid rotel client")
	} else if err := rotel.Ping(); err != nil {
		return err
	} else {
		// Power
		if power, exists := app.AppFlags.GetString("power"); exists {
			if err := SetPower(rotel, power); err != nil {
				return err
			}
		}
		if state, err := rotel.Get(); err != nil {
			return err
		} else {
			fmt.Println(state)
		}
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/rotel:client")

	config.AppFlags.FlagString("power", "", "on, off or toggle")

	// Run the command line tool
	os.Exit(rpc.Client(config, 200*time.Millisecond, Main))
}
