package ha

import (
	"fmt"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Component interface {
	// Return the unique id for the component
	Id() string

	// Return config topic
	ConfigTopic() string

	// Return command topic, or empty string if not writable
	CommandTopic() string

	// Return status topic, or empty string if not readable
	StateTopic() string

	// Set the state of the component as a string, return true if changed
	SetState(string) bool

	// Return the state of the component as a string
	State() string

	// Return JSON representation of component configuration
	JSON() ([]byte, error)
}

type component struct {
	Id_           string `json:"unique_id,omitempty"`
	ObjectId      string `json:"object_id,omitempty"`
	Name          string `json:"name,omitempty"`
	ConfigTopic_  string `json:"-"`
	CommandTopic_ string `json:"command_topic,omitempty"`
	StateTopic_   string `json:"state_topic,omitempty"`

	// The state of the component
	state string `json:"-"`
}

///////////////////////////////////////////////////////////////////////////////
// INIT

func (c *component) Init(topic, integration string, uniqueId, objectId, name string, readable, writable bool) error {
	c.Id_ = uniqueId
	c.ObjectId = objectId
	c.Name = name
	c.ConfigTopic_ = topicName(topic, integration, objectId, "config")
	if readable {
		c.StateTopic_ = topicName(topic, integration, objectId, "state")
	}
	if writable {
		c.CommandTopic_ = topicName(topic, integration, objectId, "command")
	}

	// Return success
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (self *component) Id() string {
	return self.Id_
}

// Return config topic, or empty string if not writable
func (self *component) ConfigTopic() string {
	return self.ConfigTopic_
}

// Return command topic, or empty string if not writable
func (self *component) CommandTopic() string {
	return self.CommandTopic_
}

// Return state topic, or empty string if not readable
func (self *component) StateTopic() string {
	return self.StateTopic_
}

// Set the state of the component as a string, return true if changed
func (self *component) SetState(v string) bool {
	v = strings.TrimSpace(v)
	fmt.Printf("SetState: old %q new %q\n", self.state, v)
	if self.state != v {
		self.state = v
		return true
	}
	return false
}

// Return the state of the component as a string
func (self *component) State() string {
	return self.state
}
