# Rotel Amplifier Control

This repository implements control for Rotel Audio Equipment with RS232 port,
and has been tested only with the A12 integrated amplifer. There are two ways
you can use this software:

  1. Implement a Rotel control service, which can be accessed remotely through gRPC calls from a client;
  2. Embed the Rotel module within your own software.

The following sections describe requirements, installation and how to use this software in each circumstance.

## Requirements

This software should compile and run under Go 1.12 and has the following dependencies:

  * The `gopi` application framework;
  * Google's Protocol Buffers and gRPC;
  * An RS232 communication mechanism. You can often use a USB-to-RS232
    cable in order to achieve this;
  * The software has been tested on Linux and is likely to work on
    Macintosh as necessary.

## Installation

On a Debian system, use the following commands to satisfy dependencies:

```bash
bash# sudo apt install protobuf-compiler
bash# sudo apt install libprotobuf-dev
bash# go get -u github.com/golang/protobuf/protoc-gen-go
bash# git clone git@github.com:djthorpe/rotel.git
bash# cd rotel
bash# make
```

On Macintosh with Homebrew installed, try this instead:

```bash
bash# brew install protobuf
bash# go get -u github.com/golang/protobuf/protoc-gen-go
bash# git clone git@github.com:djthorpe/rotel.git
bash# cd rotel
bash# make
```

You should end up with two binaries, `rotel-service` and `rotel-client`. These are a microservice for control of your Rotel amplifier and a client to interact with the service over the network.

## Running the Rotel Service

Your service will need access to the RS232 device (usually `/dev/ttyUSB0` if you are using a USB cable). On Debian, you can add the `dialout` group to your user in order to read and write to the device:

```bash
bash# usermod -a -G dialout $USER
```

Here are the command-line options for the `rotel-service` which you can invoke with the `-help` flag:

```bash
Syntax: rotel-service <flags>...

  -rotel.baudrate uint
    	RS232 speed (default 115200)
  -rotel.tty string
    	RS232 device (default "/dev/ttyUSB0")
  -rpc.port uint
    	Server Port
  -rpc.sslcert string
    	SSL Certificate Path
  -rpc.sslkey string
    	SSL Key Path
  -verbose
    	Verbose logging
  -debug
    	Set debugging mode on
```

For example, fire up the service so communication is performed unencrypted on port 8080:

```bash
bash# rotel-service -rpc.port 8080 -verbose
```

You can verify the service is working by adding the `-debug` flag when invoking it. You can then see all the communications between your amplifer and the service as it reads the current amplifier state.

## Using the example client

TODO

## Incorporating the Rotel module in your own code

Instead of using the serive, you may want to include the software in
your own code. The Rotel interface is as follows:

```go
type Rotel interface {
	gopi.Driver
	gopi.Publisher

	// Information
	Model() string

	// Get and set state
	Get() RotelState
	Set(RotelState) error

	// Send Command
	Send(Command) error
}
```

The main features are:

  * The `Model()` method is used to determine the amplifier that is connected, written as a string;
  * The `Get()` method returns the current state of the amplifier (Power, Volume, and so forth);
  * The `Set()` method can be called with a `RotelState` object to change current amplifier state;
  * The `Send()` method can be called to send commands, such as Play, Pause, Volume Up, and so forth;
  * The interface conforms to the `gopi.Publisher` interface, which provides methods to subscribe to events of type `RotelEvent`. Events are generated on amplifier state change.

An example application is provided in the `cmd/rotel-ctrl` folder and youcan build this with the following command:

```bash
bash# make rotel-ctrl
```

This application has a `Main` function to send commands:

```go

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
  // Send command to toggle mute on/off
  device := app.ModuleInstance("mutablehome/rotel").(rotel.Rotel)
	if err := device.Send(rotel.ROTEL_COMMAND_MUTE_TOGGLE); err != nil {
    return err
  }

	// Wait for CTRL+C
	app.Logger.Info("Waiting for CTRL+C")
	app.WaitForSignal()

	// Success
	done <- gopi.DONE
	return nil
}

```

There is also a background `EventLoop` which listens for state changes:

```go

func EventLoop(app *gopi.AppInstance, ..., done <-chan struct{}) error {
  device := app.ModuleInstance("mutablehome/rotel").(rotel.Rotel)
	evt := device.Subscribe()
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
	return nil
}
```

If you run the code and change the volume or other settings on the amplifier itself, you will see the changes reflected on your console. In order to use
the Rotel module, the interface is imported as `github.com/djthorpe/rotel`
into your code and you also need to anonymously import the module as follows:

```go
import (
  // Interfaces
  "github.com/djthorpe/rotel"

  // Modules
  _ "github.com/djthorpe/rotel/sys/rotel"
)
```

## Filing bugs and feature requests

Please don't hesitate to file bugs and feature requests on Github, here:

| https://github.com/djthorpe/rotel/issues

Even better, you can help contributing by sending pull requests to the Repository including some information with the pull request about the change you want to make. Currently this code works with the A12 amplifier across the RS232 link, but it shouldn't be hard to adapt to the A14 and using the Ethernet connection.



