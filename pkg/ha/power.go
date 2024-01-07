package ha

import (
	"encoding/json"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type PowerButton struct {
	component
	Icon string `json:"icon,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewPowerButton(topic, Id, objectId string) (Component, error) {
	self := new(PowerButton)
	if err := self.Init(topic, "switch", Id, objectId, "Power", true, true); err != nil {
		return nil, err
	}
	self.Icon = "mdi:power"

	// Return success
	return self, nil
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (self *PowerButton) JSON() ([]byte, error) {
	return json.Marshal(self)
}
