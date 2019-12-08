/*
	Rotel RS232 Control
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package rotel

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	event "github.com/djthorpe/gopi/util/event"
	rotel "github.com/djthorpe/rotel"
	term "github.com/pkg/term"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Rotel struct {
	TTY      string
	BaudRate uint
}

type driver struct {
	log gopi.Logger
	tty string
	fd  *term.Term
	buf string

	model  string
	power  rotel.Power
	input  rotel.Source
	volume rotel.Volume

	event.Tasks
	event.Publisher
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

const (
	BAUD_RATE_DEFAULT = 115200
	READ_TIMEOUT      = 100 * time.Millisecond
)

var (
	reModel  = regexp.MustCompile("^model=(\\w+)$")
	rePower  = regexp.MustCompile("^power=(on|standby)$")
	reVolume = regexp.MustCompile("^volume=(\\d+)$")
	reMute   = regexp.MustCompile("^mute=(on|off)$")
	reSource = regexp.MustCompile("^source=(\\w+)$")
	reFreq   = regexp.MustCompile("^freq=(.+)$")
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config Rotel) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("<rotel.Open>{ config=%+v }", config)

	this := new(driver)
	this.log = logger

	// Check TTY parameter
	if config.TTY == "" {
		return nil, gopi.ErrBadParameter
	} else if _, err := os.Stat(config.TTY); os.IsNotExist(err) {
		return nil, fmt.Errorf("TTY: %w", err)
	} else {
		this.tty = config.TTY
	}

	// Open term
	if fd, err := term.Open(this.tty, term.Speed(int(config.BaudRate)), term.RawMode); err != nil {
		return nil, fmt.Errorf("TTY: %w", err)
	} else {
		this.fd = fd
	}

	// Set term read timeout
	if err := this.fd.SetReadTimeout(READ_TIMEOUT); err != nil {
		defer this.fd.Close()
		return nil, fmt.Errorf("TTY: %w", err)
	}

	// Start background thread
	this.Tasks.Start(this.ReadTask)

	// Success
	return this, nil
}

func (this *driver) Close() error {
	this.log.Debug("<rotel.Close>{ tty=%v }", strconv.Quote(this.tty))

	// Remove subscribers
	this.Publisher.Close()

	// Stop background tasks
	if err := this.Tasks.Close(); err != nil {
		return err
	}

	// Close RS232 connection
	if this.fd != nil {
		if err := this.fd.Close(); err != nil {
			return err
		} else {
			this.fd = nil
		}
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *driver) String() string {
	if this.fd == nil {
		return "<rotel>{ nil }"
	} else {
		return fmt.Sprintf("<rotel>{ tty=%v model=%v power=%v input=%v volume=%v }", strconv.Quote(this.tty), strconv.Quote(this.model), this.power, this.input, this.volume)
	}
}

////////////////////////////////////////////////////////////////////////////////
// GET PARAMETERS

func (this *driver) Model() string {
	return this.model
}

func (this *driver) Power() rotel.Power {
	return this.power
}

func (this *driver) Volume() rotel.Volume {
	return this.volume
}

func (this *driver) Input() rotel.Source {
	return this.input
}

////////////////////////////////////////////////////////////////////////////////
// SET PARAMETERS

// SetPower to on or standby
func (this *driver) SetPower(value rotel.Power) error {
	this.log.Debug2("<rotel.SetPower>{ %v }", value)

	switch value {
	case rotel.ROTEL_POWER_ON:
		return this.write("power_on")
	case rotel.ROTEL_POWER_STANDBY:
		return this.write("power_off")
	case rotel.ROTEL_POWER_TOGGLE:
		return this.write("power_toggle")
	default:
		return gopi.ErrBadParameter
	}
	return gopi.ErrNotImplemented

}

// SetVolume between 0 and 96
func (this *driver) SetVolume(value rotel.Volume) error {
	this.log.Debug2("<rotel.SetVolume>{ %v }", value)

	if value > rotel.ROTEL_VOLUME_MAX {
		return gopi.ErrBadParameter
	} else {
		return this.write(fmt.Sprintf("vol_%d", value))
	}
}

// SetInput source of audio
func (this *driver) SetInput(value rotel.Source) error {
	this.log.Debug2("<rotel.SetInput>{ %v }", value)
	if str := sourceToString(value); str != "pc_usb" && str != "" {
		return this.write(str)
	} else if str == "pc_usb" {
		return this.write("pcusb")
	} else {
		return gopi.ErrBadParameter
	}
}

// SendCommand
func (this *driver) SendCommand(value rotel.Command) error {
	this.log.Debug2("<rotel.SendCommand>{ %v }", value)

	str := strings.TrimPrefix(fmt.Sprint(value), "ROTEL_COMMAND_")

	switch value {
	case rotel.ROTEL_COMMAND_PLAY, rotel.ROTEL_COMMAND_STOP, rotel.ROTEL_COMMAND_PAUSE:
		return this.write(strings.ToLower(str))
	case rotel.ROTEL_COMMAND_TRACK_NEXT:
		return this.write("trkf")
	case rotel.ROTEL_COMMAND_TRACK_PREV:
		return this.write("trkb")
	case rotel.ROTEL_COMMAND_MUTE_ON, rotel.ROTEL_COMMAND_MUTE_OFF:
		return this.write(strings.ToLower(str))
	case rotel.ROTEL_COMMAND_MUTE_TOGGLE:
		return this.write("mute")
	case rotel.ROTEL_COMMAND_BYPASS_OFF, rotel.ROTEL_COMMAND_BYPASS_ON:
		return this.write(strings.ToLower(str))
	case rotel.ROTEL_COMMAND_BASS_UP, rotel.ROTEL_COMMAND_TREBLE_UP, rotel.ROTEL_COMMAND_BASS_DOWN, rotel.ROTEL_COMMAND_TREBLE_DOWN:
		return this.write(strings.ToLower(str))
	case rotel.ROTEL_COMMAND_BASS_RESET:
		return this.write("bass_000")
	case rotel.ROTEL_COMMAND_TREBLE_RESET:
		return this.write("treble_000")
	case rotel.ROTEL_COMMAND_BALANCE_LEFT:
		return this.write("balance_l")
	case rotel.ROTEL_COMMAND_BALANCE_RIGHT:
		return this.write("balance_r")
	case rotel.ROTEL_COMMAND_BALANCE_RESET:
		return this.write("balance_000")
	case rotel.ROTEL_COMMAND_SPEAKER_A_TOGGLE:
		return this.write("speaker_a")
	case rotel.ROTEL_COMMAND_SPEAKER_B_TOGGLE:
		return this.write("speaker_b")
	case rotel.ROTEL_COMMAND_SPEAKER_A_ON, rotel.ROTEL_COMMAND_SPEAKER_A_OFF, rotel.ROTEL_COMMAND_SPEAKER_B_ON, rotel.ROTEL_COMMAND_SPEAKER_B_OFF:
		return this.write(strings.ToLower(str))
	case rotel.ROTEL_COMMAND_DIMMER_TOGGLE:
		return this.write("dimmer")
	case
		rotel.ROTEL_COMMAND_DIMMER_0, rotel.ROTEL_COMMAND_DIMMER_1, rotel.ROTEL_COMMAND_DIMMER_2,
		rotel.ROTEL_COMMAND_DIMMER_3, rotel.ROTEL_COMMAND_DIMMER_4, rotel.ROTEL_COMMAND_DIMMER_5,
		rotel.ROTEL_COMMAND_DIMMER_6:
		return this.write(strings.ToLower(str))
	case rotel.ROTEL_COMMAND_RS232_UPDATE_ON, rotel.ROTEL_COMMAND_RS232_UPDATE_OFF:
		return this.write(strings.ToLower(str))
	default:
		return gopi.ErrBadParameter
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *driver) write(command string) error {
	this.log.Debug2("<rotel.Write>{ %v! }", command)
	_, err := this.fd.Write([]byte(command + "!"))
	return err
}

func (this *driver) read(command string) error {
	this.log.Debug2("<rotel.Read>{ %v? }", command)
	_, err := this.fd.Write([]byte(command + "?"))
	return err
}

func (this *driver) parse(commands []string) error {
	this.log.Debug2("<rotel.Parse>{ %v }", strconv.Quote(strings.Join(commands, ",")))
	for _, command := range commands {
		if value := reModel.FindStringSubmatch(command); len(value) > 1 {
			this.model = value[1]
		} else if value := rePower.FindStringSubmatch(command); len(value) > 1 {
			switch value[1] {
			case "on":
				this.evtPower(rotel.ROTEL_POWER_ON)
			case "standby":
				this.evtPower(rotel.ROTEL_POWER_STANDBY)
			default:
				this.evtPower(rotel.ROTEL_POWER_OTHER)
			}
		} else if value := reSource.FindStringSubmatch(command); len(value) > 1 {
			if source := stringToSource(value[1]); source != rotel.ROTEL_SOURCE_NONE {
				this.evtInput(source)
			} else {
				this.evtInput(rotel.ROTEL_SOURCE_OTHER)
			}
		} else if value := reVolume.FindStringSubmatch(command); len(value) > 1 {
			if v, err := strconv.ParseUint(value[1], 10, 32); err == nil && v >= 0 && v <= uint64(rotel.ROTEL_VOLUME_MAX) {
				this.evtVolume(rotel.Volume(v))
			} else {
				return fmt.Errorf("Cannot parse: %v", strconv.Quote(command))
			}
		} else if value := reFreq.FindStringSubmatch(command); len(value) > 1 {
			// Do nothing with this
		} else {
			return fmt.Errorf("Cannot parse: %v", strconv.Quote(command))
		}
	}

	// Success
	return nil
}

func (this *driver) retrieveparams() error {
	if this.fd == nil {
		// If no file descriptor, do nothing
		return nil
	} else if this.model == "" {
		return this.read("model")
	} else if this.power == rotel.ROTEL_POWER_NONE {
		return this.read("power")
	} else if this.power == rotel.ROTEL_POWER_ON && this.input == rotel.ROTEL_SOURCE_NONE {
		return this.read("source")
	} else if this.power == rotel.ROTEL_POWER_ON && this.volume == rotel.ROTEL_VOLUME_NONE {
		return this.read("volume")
	}

	// Nothing to do
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BACKGROUND TASKS

func (this *driver) ReadTask(start chan<- event.Signal, stop <-chan event.Signal) error {
	start <- gopi.DONE
	timer := time.NewTicker(5 * time.Second)

FOR_LOOP:
	for {
		select {
		case <-stop:
			break FOR_LOOP
		case <-timer.C:
			if err := this.retrieveparams(); err != nil {
				this.log.Warn("ReadTask: %v", err)
			}
		default:
			for {
				if n, err := this.fd.Available(); err != nil {
					this.log.Error("ReadTask: %v", err)
					break FOR_LOOP
				} else if n > 0 {
					buf := make([]byte, n)
					if _, err := this.fd.Read(buf); err != nil {
						this.log.Warn("ReadTask: %v", err)
					} else {
						// Append buffer and extract first field
						this.buf += string(buf)
						if fields := strings.Split(this.buf, "$"); len(fields) > 0 {
							this.buf = fields[len(fields)-1]
							if err := this.parse(fields[0 : len(fields)-1]); err != nil {
								this.log.Warn("ReadTask: %v", err)
							}
						}
					}
				} else {
					break
				}
			}
		}
	}

	// Stop timer
	timer.Stop()

	// Success
	return nil
}
