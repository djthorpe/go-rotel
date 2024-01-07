package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	// Package imports
	ha "github.com/djthorpe/go-rotel/pkg/ha"
	mosquitto "github.com/mutablelogic/go-mosquitto/pkg/mosquitto"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-mosquitto"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	*log.Logger // Embedded logger

	client *mosquitto.Client // MQTT
	ha     *ha.HA            // Home assistant
	qos    int
	topic  string

	// Event channel
	evtch chan *mosquitto.Event

	// State change channel
	statech chan StateChange

	// Online/Offline messages
	topicStatusId string
}

type StateChange struct {
	Component ha.Component
	Data      []byte
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewApp(ctx context.Context, prefix, broker string, qos int, topic string) (*App, error) {
	self := new(App)

	// Connect to broker
	client, err := mosquitto.New(ctx, broker, func(evt *mosquitto.Event) {
		if evt.Type == MOSQ_FLAG_EVENT_MESSAGE {
			if self.evtch != nil {
				self.evtch <- evt
			}
		}
	})
	if err != nil {
		return nil, err
	}

	// Home assistant
	ha, err := ha.New(topic, self.StateCallback)
	if err != nil {
		return nil, err
	}

	// Initialise logger
	self.Logger = log.New(os.Stderr, prefix+": ", log.LstdFlags)

	// Set app parameters
	self.qos = qos
	self.client = client
	self.ha = ha

	// Return success
	return self, nil
}

///////////////////////////////////////////////////////////////////////////////
// RUN

// runloop for the rotel app
func (self *App) Run(ctx context.Context) error {
	var result error

	// Create channels for events and state changes
	self.evtch = make(chan *mosquitto.Event, 1)
	self.statech = make(chan StateChange, 1)

	// Subscribe to the "status" topic to get online/offline messages
	self.topicStatusId = self.ha.TopicStatus()
	if _, err := self.client.Subscribe(self.topicStatusId, mosquitto.OptQoS(self.qos)); err != nil {
		return err
	}

	// Add a power button
	power, err := self.ha.AddPowerButton("rotel_amp00_power", "rotel_amp00")
	if err != nil {
		return err
	}
	if err := self.PublishComponent(power, true); err != nil {
		return err
	}

	// Add a volume slider
	volume, err := self.ha.AddVolume("rotel_amp00_volume", "rotel_amp00")
	if err != nil {
		return err
	}
	volume.(*ha.Volume).SetRange(0, 50)
	if err := self.PublishComponent(volume, true); err != nil {
		return err
	}

	// Add input source
	input, err := self.ha.AddInput("rotel_amp00_input", "rotel_amp00", []string{
		"CD",
		"COAX1",
		"COAX2",
		"OPT1",
		"OPT2",
		"PHONO",
	})
	if err != nil {
		return err
	}
	if err := self.PublishComponent(input, true); err != nil {
		return err
	}

	// Change state every three seconds
	timer := time.NewTicker(5 * time.Second)

FOR_LOOP:
	for {
		select {
		case <-ctx.Done():
			break FOR_LOOP
		case evt := <-self.evtch:
			if evt.Topic == self.topicStatusId {
				if err := self.ha.SetStatus(string(evt.Data)); err != nil {
					log.Println("error setting status:", err)
				} else {
					log.Println("Home assistant status has changed:", self.ha)
				}
			} else if err := self.ha.Command(evt.Topic, evt.Data); err != nil {
				fmt.Println("other event=", evt)
			}
		case <-timer.C:
			if power.State() == "ON" {
				self.StateCallback(power, []byte("OFF"))
			} else {
				self.StateCallback(power, []byte("ON"))
			}
		case evt := <-self.statech:
			self.Logger.Println("publishing", string(evt.Data), "to", evt.Component.StateTopic())
			if _, err := self.client.Publish(evt.Component.StateTopic(), evt.Data); err != nil {
				return err
			}
		}
	}

	// Unpublish components
	for _, component := range self.ha.Components() {
		if err := self.PublishComponent(component, false); err != nil {
			result = errors.Join(result, err)
		}
	}

	// Close connection
	if err := self.client.Close(); err != nil {
		result = errors.Join(result, err)
	}

	// Close channels
	close(self.statech)
	close(self.evtch)

	// Return any errors
	return result
}

func (self *App) PublishComponent(component ha.Component, on bool) error {
	data, err := component.JSON()
	if err != nil {
		return err
	}
	if on {
		self.Logger.Println("publishing", string(data), "to", component.ConfigTopic())
		if _, err := self.client.Publish(component.ConfigTopic(), data, mosquitto.OptRetain()); err != nil {
			return err
		}
		if topic := component.CommandTopic(); topic != "" {
			if _, err := self.client.Subscribe(component.CommandTopic(), mosquitto.OptQoS(self.qos)); err != nil {
				return err
			}
		}
	} else {
		if _, err := self.client.Publish(component.ConfigTopic(), []byte{}, mosquitto.OptRetain()); err != nil {
			return err
		}
		if topic := component.CommandTopic(); topic != "" {
			if _, err := self.client.Unsubscribe(component.CommandTopic()); err != nil {
				return err
			}
		}
	}
	// Return success
	return nil
}

func (self *App) StateCallback(component ha.Component, data []byte) error {
	if component == nil || data == nil {
		return ErrBadParameter.Withf("invalid component or payload data")
	}

	self.Logger.Println("setting component state to", string(data), "for", component.StateTopic())

	if component.SetState(string(data)) {
		payload := []byte(component.State())
		self.Logger.Println("state change", component.StateTopic(), string(payload))
		self.statech <- StateChange{component, payload}
	} else {
		self.Logger.Println("state change ignored", component.StateTopic(), string(data))
	}

	// Return success
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (self *App) String() string {
	str := "<app"
	if self.client != nil {
		str += fmt.Sprintf(" client=%v", self.client)
	}
	str += fmt.Sprintf(" qos=%d", self.qos)
	if self.topic != "" {
		str += fmt.Sprintf(" topic=%v", self.topic)
	}
	return str + ">"
}
