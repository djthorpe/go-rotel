package ha

import (
	"encoding/json"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Volume struct {
	component
	Icon string   `json:"icon,omitempty"`
	Min  *float32 `json:"min,omitempty"`
	Max  *float32 `json:"max,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewVolume(topic, Id, objectId string) (*Volume, error) {
	self := new(Volume)
	if err := self.Init(topic, "number", Id, objectId, "Volume", true, true); err != nil {
		return nil, err
	}
	self.Icon = "mdi:volume-high"

	// Return success
	return self, nil
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (self *Volume) SetRange(min, max float32) {
	self.Min = &min
	self.Max = &max
}

func (self *Volume) JSON() ([]byte, error) {
	return json.Marshal(self)
}
