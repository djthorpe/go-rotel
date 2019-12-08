/*
	Rotel RS232 Control
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package rotel

import (
	// Frameworks
	"fmt"

	gopi "github.com/djthorpe/gopi"
	rotel "github.com/djthorpe/rotel"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

type evt struct {
	source gopi.Driver
	typ    rotel.EventType
	power  rotel.Power
	input  rotel.Source
	volume rotel.Volume
}

////////////////////////////////////////////////////////////////////////////////
// EMIT EVENTS

func (this *driver) evtPower(value rotel.Power) {
	if this.power != value {
		this.power = value
		this.Emit(&evt{
			source: this,
			typ:    rotel.EVENT_TYPE_POWER,
			power:  this.power,
		})
	}
}

func (this *driver) evtInput(value rotel.Source) {
	if this.input != value {
		this.input = value
		this.Emit(&evt{
			source: this,
			typ:    rotel.EVENT_TYPE_INPUT,
			input:  value,
		})
	}
}

func (this *driver) evtVolume(value rotel.Volume) {
	if this.volume != value {
		this.volume = value
		this.Emit(&evt{
			source: this,
			typ:    rotel.EVENT_TYPE_VOLUME,
			volume: value,
		})
	}
}

////////////////////////////////////////////////////////////////////////////////
// EVENT IMPLEMENTATION

func (this *evt) Name() string {
	return "RotelEvent"
}

func (this *evt) Source() gopi.Driver {
	return this.source
}

func (this *evt) Type() rotel.EventType {
	return this.typ
}

func (this *evt) Power() rotel.Power {
	return this.power
}

func (this *evt) Input() rotel.Source {
	return this.input
}

func (this *evt) Volume() rotel.Volume {
	return this.volume
}

func (this *evt) String() string {
	switch this.typ {
	case rotel.EVENT_TYPE_POWER:
		return fmt.Sprintf("<rotel.Event>{ type=%v power=%v }", this.typ, this.power)
	case rotel.EVENT_TYPE_INPUT:
		return fmt.Sprintf("<rotel.Event>{ type=%v input=%v }", this.typ, this.input)
	case rotel.EVENT_TYPE_VOLUME:
		return fmt.Sprintf("<rotel.Event>{ type=%v volume=%v }", this.typ, this.volume)
	default:
		return fmt.Sprintf("<rotel.Event>{ type=%v }", this.typ)
	}
}
