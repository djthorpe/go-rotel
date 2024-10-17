package rotel

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	// Packages

	term "github.com/pkg/term"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

// Ref: https://www.rotel.com/sites/default/files/product/rs232/A12-A14%20Protocol.pdf

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	TTY     string        `yaml:"tty"`
	Baud    uint          `yaml:"baud"`
	Timeout time.Duration `yaml:"timeout"`
}

type Rotel struct {
	state
	fd  *term.Term // TTY file handle
	buf *strings.Builder
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DEFAULT_TTY         = "/dev/ttyUSB0"
	DEFAULT_TTY_BAUD    = 115200
	DEFAULT_TTY_TIMEOUT = 100 * time.Millisecond
	deltaUpdate         = 500 * time.Millisecond
	VOLUME_MIN          = 1
	VOLUME_MAX          = 96
	TONE_MIN            = -10 // Bass and treble
	TONE_MAX            = 10  // Bass and treble
)

var (
	SOURCES = []string{
		"pc_usb", "cd", "coax1", "coax2", "opt1", "opt2", "aux1", "aux2", "tuner", "phono", "usb", "bluetooth",
	}
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewWithConfig(cfg Config) (*Rotel, error) {
	self := new(Rotel)

	// Set tty from config
	if cfg.TTY == "" {
		cfg.TTY = DEFAULT_TTY
	}
	if cfg.Baud == 0 {
		cfg.Baud = DEFAULT_TTY_BAUD
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = DEFAULT_TTY_TIMEOUT
	}

	// Check parameters
	if _, err := os.Stat(cfg.TTY); os.IsNotExist(err) {
		return nil, ErrBadParameter.With("tty: ", strconv.Quote(cfg.TTY))
	} else if err != nil {
		return nil, err
	}

	// Open term
	if fd, err := term.Open(cfg.TTY, term.Speed(int(cfg.Baud)), term.RawMode); err != nil {
		return nil, err
	} else {
		self.fd = fd
		self.buf = new(strings.Builder)
	}

	// Set term read timeout
	if err := self.fd.SetReadTimeout(cfg.Timeout); err != nil {
		defer self.fd.Close()
		return nil, err
	}

	// Return success
	return self, nil
}

func (self *Rotel) Run(ctx context.Context, ch chan<- Event) error {
	// Update rotel status every 100ms
	timer := time.NewTimer(100 * time.Millisecond)
	defer timer.Stop()

	// Loop handling messages until done
FOR_LOOP:
	for {
		select {
		case <-ctx.Done():
			break FOR_LOOP
		case <-timer.C:
			if cmd := self.state.Update(false); cmd != "" {
				if err := self.writetty(cmd); err != nil {
					send(ch, Event{ROTEL_FLAG_NONE, fmt.Errorf("writetty: %w", err)})
				}
			}
			timer.Reset(time.Millisecond * 500)
		default:
			if err := self.readtty(ch); err != nil {
				send(ch, Event{ROTEL_FLAG_NONE, fmt.Errorf("readtty: %w", err)})
			}
		}
	}

	// Close RS232 connection
	var result error
	if self.fd != nil {
		if err := self.fd.Close(); err != nil {
			result = errors.Join(result, err)
		}
	}

	// Clear resources
	self.fd = nil
	self.buf = nil

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (self *Rotel) SetPower(state bool) error {
	if state {
		return self.writetty("power_on!")
	} else {
		return self.writetty("power_off!")
	}
}

func (self *Rotel) SetSpeaker(state bool, speaker string) error {
	// Cannot set value when power is off
	if !self.Power() {
		return ErrOutOfOrder.With("SetSpeaker")
	}

	if state {
		return self.writetty("speaker_" + speaker + "_on!")
	} else {
		return self.writetty("speaker_" + speaker + "_off!")
	}
}

func (self *Rotel) SetSource(value string) error {
	// Cannot set value when power is off
	if !self.Power() {
		return ErrOutOfOrder.With("SetSource")
	}

	// Check parameter and send command
	switch value {
	case "pc_usb":
		return self.writetty("pcusb!")
	case "cd", "coax1", "coax2", "opt1", "opt2", "aux1", "aux2", "tuner", "phono", "usb", "bluetooth":
		return self.writetty(value + "!")
	default:
		return ErrBadParameter.Withf("invalid source: %q", value)
	}
}

func (self *Rotel) SetVolume(value uint) error {
	// Cannot set value when power is off
	if self.Power() == false {
		return ErrOutOfOrder.With("SetVolume")
	}

	// Check parameter and send command
	if value < 1 || value > 96 {
		return ErrBadParameter.With("invalid volume: %d", value)
	} else {
		return self.writetty(fmt.Sprintf("vol_%02d!", value))
	}
}

func (self *Rotel) SetBass(value int) error {
	// Cannot set value when power is off
	if self.Power() == false {
		return ErrOutOfOrder.With("SetBass")
	}

	// Check parameter and send command
	if value < TONE_MIN || value > TONE_MAX {
		return ErrBadParameter.With("SetBass")
	} else if value == 0 {
		return self.writetty("bass_000!")
	} else if value < 0 {
		return self.writetty(fmt.Sprint("bass_", value, "!"))
	} else {
		return self.writetty(fmt.Sprint("bass_+", value, "!"))
	}
}

func (self *Rotel) SetTreble(value int) error {
	// Cannot set value when power is off
	if self.Power() == false {
		return ErrOutOfOrder.With("SetTreble")
	}

	// Check parameter and send command
	if value < TONE_MIN || value > TONE_MAX {
		return ErrBadParameter.With("SetTreble")
	} else if value == 0 {
		return self.writetty("treble_000!")
	} else if value < 0 {
		return self.writetty(fmt.Sprint("treble_", value, "!"))
	} else {
		return self.writetty(fmt.Sprint("treble_+", value, "!"))
	}
}

/*
func (this *Manager) SetMute(state bool) error {
	// Cannot set value when power is off
	if this.Power() == false {
		return ErrOutOfOrder.With("SetMute")
	}

	// Check parameter and send command
	if state {
		return this.writetty("mute_on!")
	} else {
		return this.writetty("mute_off!")
	}
}

func (this *Manager) SetBypass(state bool) error {
	// Cannot set value when power is off
	if this.Power() == false {
		return ErrOutOfOrder.With("SetBypass")
	}

	// Check parameter and send command
	if state {
		return this.writetty("bypass_on!")
	} else {
		return this.writetty("bypass_off!")
	}
}

func (this *Manager) SetBalance(loc string) error {
	// Cannot set value when power is off
	if this.Power() == false {
		return ErrOutOfOrder.With("SetBalance")
	}

	// Check parameter and send command
	switch loc {
	case "0":
		return this.writetty("balance_000!")
	case "L", "R":
		return this.writetty(fmt.Sprintf("balance_%v!", strings.ToLower(loc)))
	default:
		return ErrBadParameter.With("SetBalance")
	}
}

func (this *Manager) SetDimmer(value uint) error {
	// Cannot set value when power is off
	if this.Power() == false {
		return ErrOutOfOrder.With("SetDimmer")
	}

	// Check parameter and send command
	if value > 6 {
		return ErrBadParameter.With("SetDimmer")
	} else {
		return this.writetty(fmt.Sprint("dimmer_", value, "!"))
	}
}

func (this *Manager) Play() error {
	// Cannot perform action when power is off
	if this.Power() == false {
		return ErrOutOfOrder.With("Play")
	}

	// Send command
	return this.writetty("play!")
}

func (this *Manager) Stop() error {
	// Cannot perform action when power is off
	if this.Power() == false {
		return ErrOutOfOrder.With("Stop")
	}

	// Send command
	return this.writetty("stop!")
}

func (this *Manager) Pause() error {
	// Cannot perform action when power is off
	if this.Power() == false {
		return ErrOutOfOrder.With("Pause")
	}

	// Send command
	return this.writetty("pause!")
}

func (this *Manager) NextTrack() error {
	// Cannot perform action when power is off
	if this.Power() == false {
		return ErrOutOfOrder.With("NextTrack")
	}

	// Send command
	return this.writetty("trkf!")
}

func (this *Manager) PrevTrack() error {
	// Cannot perform action when power is off
	if this.Power() == false {
		return ErrOutOfOrder.With("PrevTrack")
	}

	// Send command
	return this.writetty("trkb!")
}
*/

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (self *Rotel) String() string {
	str := "<rotel"
	if self.fd != nil {
		str += fmt.Sprintf(" tty=%q", self.fd)
	}
	//str += fmt.Sprint(" ", this.State.String())
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHDOS

func (self *Rotel) readtty(ch chan<- Event) error {
	var result error
	var flags Flag

	// Append data to the buffer and parse any parameters
	buf := make([]byte, 1024)
	if n, err := self.fd.Read(buf); err == io.EOF {
		return nil
	} else if err != nil {
		return err
	} else if _, err := self.buf.Write(buf[:n]); err != nil {
		return err
	} else if fields := strings.Split(self.buf.String(), "$"); len(fields) > 0 {
		// Parse each field and update state
		for _, param := range fields[0 : len(fields)-1] {
			if flag, err := self.state.Set(param); err != nil {
				result = errors.Join(result, fmt.Errorf("%q: %w", param, err))
			} else {
				flags |= flag
			}
		}
		// Reset buffer with any remaining data not parsed
		self.buf.Reset()
		self.buf.WriteString(fields[len(fields)-1])
	}

	// If any flags set, then emit an event
	if flags != ROTEL_FLAG_NONE {
		if err := send(ch, Event{flags, nil}); err != nil {
			result = errors.Join(result, err)
		}
	}

	// Return any errors
	return result
}

func (self *Rotel) writetty(cmd string) error {
	_, err := self.fd.Write([]byte(cmd))
	return err
}

func send(ch chan<- Event, evt Event) error {
	select {
	case ch <- evt:
		return nil
	default:
		return ErrChannelBlocked
	}
}
