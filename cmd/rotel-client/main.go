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

func Main(app *gopi.AppInstance, services []gopi.RPCServiceRecord, done chan<- struct{}) error {
	// Get the client
	if stub, err := app.ClientPool.NewClientEx("gopi.Rotel", services, gopi.RPC_FLAG_SERVICE_ANY); err != nil {
		return err
	} else {
		rotel := stub.(rotel.RotelClient)
		fmt.Println(rotel)
		rotel.Ping()
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/rotel:client")

	// Run the command line tool
	os.Exit(rpc.Client(config, 200*time.Millisecond, Main))
}
