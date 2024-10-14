package main

import (
	"flag"
	"fmt"
	"os"

	// Package imports
	rotel "github.com/djthorpe/go-rotel/pkg/rotel"
	version "github.com/djthorpe/go-rotel/pkg/version"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Args struct {
	*flag.FlagSet

	// Flags
	Topic       string
	Broker      string
	Credentials string
	Qos         int
	TTY         string
	Version     bool
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultBroker      = "localhost:1833"
	defaultTopic       = "homeassistant"
	defaultCredentials = ""
)

var (
	ErrHelp = flag.ErrHelp
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewArgs(name string, args []string) (*Args, error) {
	self := &Args{
		FlagSet: flag.NewFlagSet(name, flag.ContinueOnError),
	}

	// Register flags
	self.registerFlags()

	// Parse flags
	if err := self.Parse(args); err != nil {
		return nil, err
	}
	// No arguments are allowed
	if self.NArg() > 0 {
		return nil, ErrBadParameter.Withf("unexpected argument: %q", self.Arg(0))
	}
	// Print version and exit
	if self.Version {
		version.Print(os.Stdout)
		return nil, ErrHelp
	}

	return self, nil
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (self *Args) String() string {
	str := "<flags"
	if self.Name() != "" {
		str += fmt.Sprintf(" name=%q", self.Name())
	}
	if self.Broker != "" {
		str += fmt.Sprintf(" mqtt=%q", self.Broker)
	}
	if self.Topic != "" {
		str += fmt.Sprintf(" topic=%q", self.Topic)
	}
	str += fmt.Sprintf(" qos=%d", self.Qos)
	if self.TTY != "" {
		str += fmt.Sprintf(" tty=%q", self.TTY)
	}
	str += fmt.Sprintf(" version=%v", self.Version)
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (self *Args) registerFlags() {
	self.StringVar(&self.Broker, "mqtt", defaultBroker, "MQTT broker address")
	self.StringVar(&self.Credentials, "credentials", defaultCredentials, "MQTT credentails (user:password)")
	self.StringVar(&self.Topic, "topic", defaultTopic, "Topic for messages")
	self.IntVar(&self.Qos, "qos", 0, "MQTT quality of service")
	self.StringVar(&self.TTY, "tty", rotel.DEFAULT_TTY, "TTY for Rotel device")
	self.BoolVar(&self.Version, "version", false, "Print version and exit")
}
