/*
	Rotel RS232 Control
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

// Rotel RS232 Control
package rotel

import (
	// Frameworks
	"context"
	"fmt"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	Command   uint16
	EventType uint16
	Power     uint16
	Volume    uint16
	Source    uint16
	Mute      uint16
	Tone      int16
	Balance   int16
	Dimmer    uint16
)

type Speaker struct{ A, B bool }

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

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

type RotelEvent interface {
	gopi.Event

	Type() EventType
	State() RotelState
}

type RotelState struct {
	Power
	Volume
	Mute
	Source
	Freq   string
	Bypass bool
	Treble Tone
	Bass   Tone
	Balance
	Speaker
	Dimmer
}

type RotelClient interface {
	gopi.RPCClient
	gopi.Publisher

	// Ping remote service
	Ping() error

	// Get and set state
	Get() (RotelState, error)
	Set(RotelState) error

	// Send command
	Send(Command) error

	// Stream state changes
	StreamEvents(ctx context.Context) error
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ROTEL_POWER_NONE Power = 0
	ROTEL_POWER_ON   Power = iota
	ROTEL_POWER_STANDBY
	ROTEL_POWER_TOGGLE
	ROTEL_POWER_MAX = ROTEL_POWER_TOGGLE
)

const (
	ROTEL_MUTE_NONE Mute = 0
	ROTEL_MUTE_ON   Mute = iota
	ROTEL_MUTE_OFF
	ROTEL_MUTE_TOGGLE
	ROTEL_MUTE_OTHER
	ROTEL_MUTE_MAX = ROTEL_MUTE_TOGGLE
)

const (
	ROTEL_SOURCE_NONE Source = 0
	ROTEL_SOURCE_CD   Source = iota
	ROTEL_SOURCE_COAX1
	ROTEL_SOURCE_COAX2
	ROTEL_SOURCE_OPT1
	ROTEL_SOURCE_OPT2
	ROTEL_SOURCE_AUX1
	ROTEL_SOURCE_AUX2
	ROTEL_SOURCE_TUNER
	ROTEL_SOURCE_PHONO
	ROTEL_SOURCE_USB
	ROTEL_SOURCE_BLUETOOTH
	ROTEL_SOURCE_PC_USB
	ROTEL_SOURCE_OTHER
	ROTEL_SOURCE_MAX = ROTEL_SOURCE_OTHER
)

const (
	ROTEL_VOLUME_NONE Volume = 0
	ROTEL_VOLUME_MAX  Volume = 96
)

const (
	ROTEL_TONE_NONE  Tone = 0
	ROTEL_TONE_MIN   Tone = -100
	ROTEL_TONE_MAX   Tone = 100
	ROTEL_TONE_OTHER Tone = ROTEL_TONE_MAX + 1
)

const (
	ROTEL_BALANCE_NONE      Balance = 0
	ROTEL_BALANCE_LEFT_MAX  Balance = -15
	ROTEL_BALANCE_RIGHT_MAX Balance = 15
	ROTEL_BALANCE_OTHER     Balance = ROTEL_BALANCE_RIGHT_MAX + 1
)

const (
	ROTEL_DIMMER_NONE  Dimmer = 0
	ROTEL_DIMMER_MAX   Dimmer = 9
	ROTEL_DIMMER_OTHER Dimmer = ROTEL_DIMMER_MAX + 1
)

const (
	ROTEL_COMMAND_NONE Command = 0
	ROTEL_COMMAND_PLAY Command = iota
	ROTEL_COMMAND_STOP
	ROTEL_COMMAND_PAUSE
	ROTEL_COMMAND_TRACK_NEXT
	ROTEL_COMMAND_TRACK_PREV
	ROTEL_COMMAND_MUTE_OFF
	ROTEL_COMMAND_MUTE_ON
	ROTEL_COMMAND_MUTE_TOGGLE
	ROTEL_COMMAND_VOL_UP
	ROTEL_COMMAND_VOL_DOWN
	ROTEL_COMMAND_BYPASS_OFF
	ROTEL_COMMAND_BYPASS_ON
	ROTEL_COMMAND_BASS_UP
	ROTEL_COMMAND_BASS_DOWN
	ROTEL_COMMAND_BASS_RESET
	ROTEL_COMMAND_TREBLE_UP
	ROTEL_COMMAND_TREBLE_DOWN
	ROTEL_COMMAND_TREBLE_RESET
	ROTEL_COMMAND_BALANCE_LEFT
	ROTEL_COMMAND_BALANCE_RIGHT
	ROTEL_COMMAND_BALANCE_RESET
	ROTEL_COMMAND_SPEAKER_A_TOGGLE
	ROTEL_COMMAND_SPEAKER_B_TOGGLE
	ROTEL_COMMAND_SPEAKER_A_ON
	ROTEL_COMMAND_SPEAKER_A_OFF
	ROTEL_COMMAND_SPEAKER_B_ON
	ROTEL_COMMAND_SPEAKER_B_OFF
	ROTEL_COMMAND_DIMMER_TOGGLE
	ROTEL_COMMAND_RS232_UPDATE_ON
	ROTEL_COMMAND_RS232_UPDATE_OFF
	ROTEL_COMMAND_MAX = ROTEL_COMMAND_RS232_UPDATE_OFF
)

const (
	EVENT_TYPE_NONE  EventType = 0
	EVENT_TYPE_POWER EventType = iota
	EVENT_TYPE_VOLUME
	EVENT_TYPE_SOURCE
	EVENT_TYPE_MUTE
	EVENT_TYPE_FREQ
	EVENT_TYPE_BYPASS
	EVENT_TYPE_BASS
	EVENT_TYPE_TREBLE
	EVENT_TYPE_BALANCE
	EVENT_TYPE_SPEAKER
	EVENT_TYPE_DIMMER
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (p Power) String() string {
	switch p {
	case ROTEL_POWER_NONE:
		return "ROTEL_POWER_NONE"
	case ROTEL_POWER_ON:
		return "ROTEL_POWER_ON"
	case ROTEL_POWER_STANDBY:
		return "ROTEL_POWER_STANDBY"
	case ROTEL_POWER_TOGGLE:
		return "ROTEL_POWER_TOGGLE"
	default:
		return "[?? Invalid Power value]"
	}
}

func (m Mute) String() string {
	switch m {
	case ROTEL_MUTE_NONE:
		return "ROTEL_MUTE_NONE"
	case ROTEL_MUTE_ON:
		return "ROTEL_MUTE_ON"
	case ROTEL_MUTE_OFF:
		return "ROTEL_MUTE_OFF"
	case ROTEL_MUTE_TOGGLE:
		return "ROTEL_MUTE_TOGGLE"
	case ROTEL_MUTE_OTHER:
		return "ROTEL_MUTE_OTHER"
	default:
		return "[?? Invalid Mute value]"
	}
}

func (v Volume) String() string {
	if v == ROTEL_VOLUME_NONE {
		return "ROTEL_VOLUME_NONE"
	} else if v == ROTEL_VOLUME_MAX {
		return "ROTEL_VOLUME_MAX"
	} else if v < ROTEL_VOLUME_MAX {
		return fmt.Sprintf("ROTEL_VOLUME_%d", v)
	} else {
		return "[?? Invalid Volume value]"
	}
}

func (t Tone) String() string {
	switch {
	case t == ROTEL_TONE_NONE:
		return "ROTEL_TONE_NONE"
	case t == ROTEL_TONE_MAX:
		return "ROTEL_TONE_MAX"
	case t == ROTEL_TONE_MIN:
		return "ROTEL_TONE_MIN"
	case t == ROTEL_TONE_OTHER:
		return "ROTEL_TONE_OTHER"
	case t >= ROTEL_TONE_MIN && t < ROTEL_TONE_NONE:
		return fmt.Sprintf("ROTEL_TONE_MINUS_%d", -t)
	case t <= ROTEL_TONE_MAX && t > ROTEL_TONE_NONE:
		return fmt.Sprintf("ROTEL_TONE_PLUS_%d", t)
	default:
		return "[?? Invalid Tone value]"
	}
}

func (b Balance) String() string {
	switch {
	case b == ROTEL_BALANCE_NONE:
		return "ROTEL_BALANCE_NONE"
	case b == ROTEL_BALANCE_LEFT_MAX:
		return "ROTEL_BALANCE_LEFT_MAX"
	case b == ROTEL_BALANCE_RIGHT_MAX:
		return "ROTEL_BALANCE_RIGHT_MAX"
	case b >= ROTEL_BALANCE_LEFT_MAX && b < ROTEL_BALANCE_NONE:
		return fmt.Sprintf("ROTEL_BALANCE_LEFT_%d", -b)
	case b <= ROTEL_BALANCE_RIGHT_MAX && b > ROTEL_BALANCE_NONE:
		return fmt.Sprintf("ROTEL_BALANCE_RIGHT_%d", b)
	default:
		return "[?? Invalid Balance value]"
	}
}

func (s Source) String() string {
	switch s {
	case ROTEL_SOURCE_NONE:
		return "ROTEL_SOURCE_NONE"
	case ROTEL_SOURCE_CD:
		return "ROTEL_SOURCE_CD"
	case ROTEL_SOURCE_COAX1:
		return "ROTEL_SOURCE_COAX1"
	case ROTEL_SOURCE_COAX2:
		return "ROTEL_SOURCE_COAX2"
	case ROTEL_SOURCE_OPT1:
		return "ROTEL_SOURCE_OPT1"
	case ROTEL_SOURCE_OPT2:
		return "ROTEL_SOURCE_OPT2"
	case ROTEL_SOURCE_AUX1:
		return "ROTEL_SOURCE_AUX1"
	case ROTEL_SOURCE_AUX2:
		return "ROTEL_SOURCE_AUX2"
	case ROTEL_SOURCE_TUNER:
		return "ROTEL_SOURCE_TUNER"
	case ROTEL_SOURCE_PHONO:
		return "ROTEL_SOURCE_PHONO"
	case ROTEL_SOURCE_USB:
		return "ROTEL_SOURCE_USB"
	case ROTEL_SOURCE_BLUETOOTH:
		return "ROTEL_SOURCE_BLUETOOTH"
	case ROTEL_SOURCE_PC_USB:
		return "ROTEL_SOURCE_PC_USB"
	case ROTEL_SOURCE_OTHER:
		return "ROTEL_SOURCE_OTHER"
	default:
		return "[?? Invalid Source value]"
	}
}

func (s Speaker) String() string {
	switch {
	case s.A == true && s.B == false:
		return "ROTEL_SPEAKER_A"
	case s.B == true && s.A == false:
		return "ROTEL_SPEAKER_B"
	case s.B == true && s.A == true:
		return "ROTEL_SPEAKER_BOTH"
	default:
		return "ROTEL_SPEAKER_NONE"
	}
}

func (d Dimmer) String() string {
	switch {
	case d == ROTEL_DIMMER_NONE:
		return "ROTEL_DIMMER_NONE"
	case d == ROTEL_DIMMER_MAX:
		return "ROTEL_DIMMER_MAX"
	case d == ROTEL_DIMMER_OTHER:
		return "ROTEL_DIMMER_OTHER"
	case d > ROTEL_DIMMER_NONE && d <= ROTEL_DIMMER_MAX:
		return fmt.Sprintf("ROTEL_DIMMER_%d", d)
	default:
		return "[?? Invalid Dimmer value]"
	}
}

func (e EventType) String() string {
	switch e {
	case EVENT_TYPE_NONE:
		return "EVENT_TYPE_NONE"
	case EVENT_TYPE_POWER:
		return "EVENT_TYPE_POWER"
	case EVENT_TYPE_VOLUME:
		return "EVENT_TYPE_VOLUME"
	case EVENT_TYPE_SOURCE:
		return "EVENT_TYPE_SOURCE"
	case EVENT_TYPE_MUTE:
		return "EVENT_TYPE_MUTE"
	case EVENT_TYPE_FREQ:
		return "EVENT_TYPE_FREQ"
	case EVENT_TYPE_BYPASS:
		return "EVENT_TYPE_BYPASS"
	case EVENT_TYPE_BASS:
		return "EVENT_TYPE_BASS"
	case EVENT_TYPE_TREBLE:
		return "EVENT_TYPE_TREBLE"
	case EVENT_TYPE_BALANCE:
		return "EVENT_TYPE_BALANCE"
	case EVENT_TYPE_SPEAKER:
		return "EVENT_TYPE_SPEAKER"
	case EVENT_TYPE_DIMMER:
		return "EVENT_TYPE_DIMMER"
	default:
		return "[?? Invalid EventType value]"
	}
}

func (c Command) String() string {
	switch c {
	case ROTEL_COMMAND_NONE:
		return "ROTEL_COMMAND_NONE"
	case ROTEL_COMMAND_PLAY:
		return "ROTEL_COMMAND_PLAY"
	case ROTEL_COMMAND_STOP:
		return "ROTEL_COMMAND_STOP"
	case ROTEL_COMMAND_PAUSE:
		return "ROTEL_COMMAND_PAUSE"
	case ROTEL_COMMAND_TRACK_NEXT:
		return "ROTEL_COMMAND_TRACK_NEXT"
	case ROTEL_COMMAND_TRACK_PREV:
		return "ROTEL_COMMAND_TRACK_PREV"
	case ROTEL_COMMAND_MUTE_OFF:
		return "ROTEL_COMMAND_MUTE_OFF"
	case ROTEL_COMMAND_MUTE_ON:
		return "ROTEL_COMMAND_MUTE_ON"
	case ROTEL_COMMAND_MUTE_TOGGLE:
		return "ROTEL_COMMAND_MUTE_TOGGLE"
	case ROTEL_COMMAND_VOL_UP:
		return "ROTEL_COMMAND_VOL_UP"
	case ROTEL_COMMAND_VOL_DOWN:
		return "ROTEL_COMMAND_VOL_DOWN"
	case ROTEL_COMMAND_BYPASS_OFF:
		return "ROTEL_COMMAND_BYPASS_OFF"
	case ROTEL_COMMAND_BYPASS_ON:
		return "ROTEL_COMMAND_BYPASS_ON"
	case ROTEL_COMMAND_BASS_UP:
		return "ROTEL_COMMAND_BASS_UP"
	case ROTEL_COMMAND_BASS_DOWN:
		return "ROTEL_COMMAND_BASS_DOWN"
	case ROTEL_COMMAND_BASS_RESET:
		return "ROTEL_COMMAND_BASS_RESET"
	case ROTEL_COMMAND_TREBLE_UP:
		return "ROTEL_COMMAND_TREBLE_UP"
	case ROTEL_COMMAND_TREBLE_DOWN:
		return "ROTEL_COMMAND_TREBLE_DOWN"
	case ROTEL_COMMAND_TREBLE_RESET:
		return "ROTEL_COMMAND_TREBLE_RESET"
	case ROTEL_COMMAND_BALANCE_LEFT:
		return "ROTEL_COMMAND_BALANCE_LEFT"
	case ROTEL_COMMAND_BALANCE_RIGHT:
		return "ROTEL_COMMAND_BALANCE_RIGHT"
	case ROTEL_COMMAND_BALANCE_RESET:
		return "ROTEL_COMMAND_BALANCE_RESET"
	case ROTEL_COMMAND_SPEAKER_A_TOGGLE:
		return "ROTEL_COMMAND_SPEAKER_A_TOGGLE"
	case ROTEL_COMMAND_SPEAKER_B_TOGGLE:
		return "ROTEL_COMMAND_SPEAKER_B_TOGGLE"
	case ROTEL_COMMAND_SPEAKER_A_ON:
		return "ROTEL_COMMAND_SPEAKER_A_ON"
	case ROTEL_COMMAND_SPEAKER_A_OFF:
		return "ROTEL_COMMAND_SPEAKER_A_OFF"
	case ROTEL_COMMAND_SPEAKER_B_ON:
		return "ROTEL_COMMAND_SPEAKER_B_ON"
	case ROTEL_COMMAND_SPEAKER_B_OFF:
		return "ROTEL_COMMAND_SPEAKER_B_OFF"
	case ROTEL_COMMAND_DIMMER_TOGGLE:
		return "ROTEL_COMMAND_DIMMER_TOGGLE"
	case ROTEL_COMMAND_RS232_UPDATE_ON:
		return "ROTEL_COMMAND_RS232_UPDATE_ON"
	case ROTEL_COMMAND_RS232_UPDATE_OFF:
		return "ROTEL_COMMAND_RS232_UPDATE_OFF"
	default:
		return "[?? Invalid Command value]"
	}
}
