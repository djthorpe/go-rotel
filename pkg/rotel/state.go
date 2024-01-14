package rotel

import (
	"fmt"
	"regexp"
	"strconv"

	// Modules
	. "github.com/djthorpe/go-errors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type state struct {
	model         string
	power         string
	update        string // rs232 update
	volume, mute  string
	bass, treble  string
	balance       []string
	source        string
	freq          string
	bypass        string
	speaker       string
	dimmer        string
	volume_update bool
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

var (
	commands = []struct {
		re *regexp.Regexp
		fn func(this *state, args []string) (Flag, error)
	}{
		{regexp.MustCompile("^model=(\\w+)$"), SetModel},
		{regexp.MustCompile("^power=(on|standby)$"), SetPower},
		{regexp.MustCompile("^volume=(\\d+)$"), SetVolume},
		{regexp.MustCompile("^update_mode=(auto|manual)$"), SetUpdateMode},
		{regexp.MustCompile("^bass=([\\+\\-]?\\d+)$"), SetBass},
		{regexp.MustCompile("^treble=([\\+\\-]?\\d+)$"), SetTreble},
		{regexp.MustCompile("^balance=([LR]?)(\\d+)$"), SetBalance},
		{regexp.MustCompile("^mute=(on|off)$"), SetMute},
		{regexp.MustCompile("^source=(\\w+)$"), SetSource},
		{regexp.MustCompile("^freq=(.+)$"), SetFreq},
		{regexp.MustCompile("^bypass=(on|off)$"), SetBypass},
		{regexp.MustCompile("^speaker=(a|b|a_b|off)$"), SetSpeaker},
		{regexp.MustCompile("^dimmer=(\\d+)$"), SetDimmer},
	}
)

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *state) Model() string {
	return this.model
}

func (this *state) Power() bool {
	return this.power == "on"
}

func (this *state) Volume() uint {
	if this.power == "on" {
		if vol, err := strconv.ParseUint(this.volume, 0, 32); err == nil {
			return uint(vol)
		}
	}
	return 0
}

func (this *state) Bass() int {
	if this.power == "on" {
		if bass, err := strconv.ParseInt(this.bass, 0, 32); err == nil {
			return int(bass)
		}
	}
	return 0
}

func (this *state) Treble() int {
	if this.power == "on" {
		if treble, err := strconv.ParseInt(this.treble, 0, 32); err == nil {
			return int(treble)
		}
	}
	return 0
}

func (this *state) Balance() (string, uint) {
	if this.power == "on" && this.balance != nil {
		if scalar, err := strconv.ParseUint(this.balance[1], 0, 32); err == nil {
			return this.balance[0], uint(scalar)
		}
	}
	return "", 0
}

func (this *state) Dimmer() uint {
	if this.power == "on" {
		if dimmer, err := strconv.ParseUint(this.dimmer, 0, 32); err == nil {
			return uint(dimmer)
		}
	}
	return 0
}

func (this *state) Muted() bool {
	if this.power == "on" && this.mute == "on" {
		return true
	} else {
		return false
	}
}

func (this *state) Bypass() bool {
	if this.power == "on" && this.bypass == "on" {
		return true
	} else {
		return false
	}
}

func (this *state) Source() string {
	if this.power == "on" {
		return this.source
	} else {
		return ""
	}
}

func (this *state) Freq() string {
	if this.power == "on" && this.freq != "off" {
		return this.freq
	} else {
		return ""
	}
}

func (this *state) Speakers() string {
	if this.power == "on" {
		return this.speaker
	}
	return ""
}

func (this *state) SpeakerA() bool {
	if this.power == "on" {
		switch this.speaker {
		case "a":
			return true
		case "b":
			return false
		case "a_b":
			return true
		}
	}
	return false
}

func (this *state) SpeakerB() bool {
	if this.power == "on" {
		switch this.speaker {
		case "a":
			return false
		case "b":
			return true
		case "a_b":
			return true
		}
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Update returns a query to get state of an unknown value
func (this *state) Update() string {
	switch {
	case this.model == "":
		return "model?"
	case this.power == "":
		return "power?"
	case this.power != "on": // When power is off, don't read other values
		return ""
	case this.volume == "" || this.volume == "0" || this.volume_update == true:
		return "volume?"
	case this.source == "":
		return "source?"
	case this.update == "":
		return "rs232_update_on!"
	case this.freq == "":
		return "freq?"
	case this.bypass == "":
		return "bypass?"
	case this.speaker == "":
		return "speaker?"
	case this.mute == "":
		return "mute?"
	case this.bass == "":
		return "bass?"
	case this.treble == "":
		return "treble?"
	case this.balance == nil:
		return "balance?"
	case this.dimmer == "":
		return "dimmer?"
	}

	// By default, no state needs read
	return ""
}

// Set sets state from data coming from amp
func (this *state) Set(param string) (Flag, error) {
	for _, command := range commands {
		if args := command.re.FindStringSubmatch(param); len(args) != 0 {
			return command.fn(this, args[1:])
		}
	}
	// Cannot match command
	return 0, ErrUnexpectedResponse.With(strconv.Quote(param))
}

func SetModel(this *state, args []string) (Flag, error) {
	if args[0] == "" {
		return 0, ErrBadParameter.With("SetModel")
	} else if this.model != args[0] {
		this.model = args[0]
		return ROTEL_FLAG_MODEL, nil
	}
	return 0, nil
}

func SetPower(this *state, args []string) (Flag, error) {
	if args[0] == "" {
		return 0, ErrBadParameter.With("SetPower")
	} else if this.power == args[0] {
		return 0, nil
	}
	this.power = args[0]

	// If the power is switched on, then update the volume
	if this.power == "on" {
		this.volume_update = true
	}

	// Return the power changed flag
	return ROTEL_FLAG_POWER, nil
}

func SetUpdateMode(this *state, args []string) (Flag, error) {
	if args[0] == "" {
		return 0, ErrBadParameter.With("SetUpdateMode")
	}
	this.update = args[0]
	return 0, nil
}

func SetVolume(this *state, args []string) (Flag, error) {
	this.volume_update = false
	if volume, err := strconv.ParseUint(args[0], 10, 32); err != nil {
		return 0, err
	} else if volume_ := fmt.Sprint(volume); volume_ != this.volume {
		this.volume = volume_
		return ROTEL_FLAG_VOLUME, nil
	}
	return 0, nil
}

func SetBass(this *state, args []string) (Flag, error) {
	if bass, err := strconv.ParseInt(args[0], 10, 32); err != nil {
		return 0, err
	} else if bass_ := fmt.Sprint(bass); bass_ != this.bass {
		this.bass = bass_
		return ROTEL_FLAG_BASS, nil
	}
	return 0, nil
}

func SetTreble(this *state, args []string) (Flag, error) {
	if treble, err := strconv.ParseInt(args[0], 10, 32); err != nil {
		return 0, err
	} else if treble_ := fmt.Sprint(treble); treble_ != this.treble {
		this.treble = treble_
		return ROTEL_FLAG_TREBLE, nil
	}
	return 0, nil
}

func SetBalance(this *state, args []string) (Flag, error) {
	if scalar, err := strconv.ParseUint(args[1], 10, 32); err != nil {
		return 0, err
	} else {
		scalar_ := fmt.Sprint(scalar)
		if this.balance == nil || scalar_ != this.balance[1] || args[0] != this.balance[0] {
			this.balance = []string{args[0], fmt.Sprint(scalar)}
			return ROTEL_FLAG_BALANCE, nil
		}
	}
	return 0, nil
}

func SetMute(this *state, args []string) (Flag, error) {
	if args[0] != this.mute {
		this.mute = args[0]
		return ROTEL_FLAG_MUTE, nil
	}
	return 0, nil
}

func SetSource(this *state, args []string) (Flag, error) {
	if args[0] != this.source {
		this.source = args[0]
		return ROTEL_FLAG_SOURCE, nil
	}
	return 0, nil
}

func SetFreq(this *state, args []string) (Flag, error) {
	if args[0] != this.freq {
		this.freq = args[0]
		return ROTEL_FLAG_FREQ, nil
	}
	return 0, nil
}

func SetBypass(this *state, args []string) (Flag, error) {
	if args[0] != this.bypass {
		this.bypass = args[0]
		return ROTEL_FLAG_BYPASS, nil
	}
	return 0, nil
}

func SetSpeaker(this *state, args []string) (Flag, error) {
	if args[0] != this.speaker {
		this.speaker = args[0]
		return ROTEL_FLAG_SPEAKER, nil
	}
	return 0, nil
}

func SetDimmer(this *state, args []string) (Flag, error) {
	if dimmer, err := strconv.ParseUint(args[0], 10, 32); err != nil {
		return 0, err
	} else if dimmer_ := fmt.Sprint(dimmer); this.dimmer != dimmer_ {
		this.dimmer = dimmer_
		return ROTEL_FLAG_DIMMER, nil
	}
	return 0, nil
}
