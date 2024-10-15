package ha

import (
	"encoding/json"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Speaker struct {
	component
	Icon string `json:"icon,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewSpeaker(topic, Id, objectId string, speakerName string) (*Speaker, error) {
	self := new(Speaker)
	if err := self.Init(topic, "switch", Id, objectId, speakerName, true, true); err != nil {
		return nil, err
	}
	self.Icon = "mdi:speaker"

	// Return success
	return self, nil
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (self *Speaker) JSON() ([]byte, error) {
	return json.Marshal(self)
}
