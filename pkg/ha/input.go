package ha

import (
	"encoding/json"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Input struct {
	component
	Icon    string   `json:"icon,omitempty"`
	Options []string `json:"options,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewInput(topic, Id, objectId string, options []string) (*Input, error) {
	self := new(Input)
	if err := self.Init(topic, "select", Id, objectId, "Input", true, true); err != nil {
		return nil, err
	}
	self.Icon = "mdi:import"
	self.Options = options

	// Return success
	return self, nil
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (self *Input) JSON() ([]byte, error) {
	return json.Marshal(self)
}
