package ha

import (
	"encoding/json"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Slider struct {
	component
	Icon string   `json:"icon,omitempty"`
	Min  *float32 `json:"min,omitempty"`
	Max  *float32 `json:"max,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewSlider(topic, Id, objectId, name string) (*Slider, error) {
	self := new(Slider)
	if err := self.Init(topic, "number", Id, objectId, name, true, true); err != nil {
		return nil, err
	}
	self.Icon = "mdi:tune-variant"

	// Return success
	return self, nil
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (self *Slider) SetRange(min, max float32) {
	self.Min = &min
	self.Max = &max
}

func (self *Slider) JSON() ([]byte, error) {
	return json.Marshal(self)
}
