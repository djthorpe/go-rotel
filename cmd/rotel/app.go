package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	// Package imports
	ha "github.com/djthorpe/go-rotel/pkg/ha"
	rotel "github.com/djthorpe/go-rotel/pkg/rotel"
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
	rotel  *rotel.Rotel
	ha     *ha.HA // Home assistant
	qos    int
	topic  string
	id     string

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

func NewApp(ctx context.Context, prefix, broker, credentials, id string, qos int, topic string, tty string) (*App, error) {
	self := new(App)

	// Broker configuration
	cfg := mosquitto.NewConfigWithBroker(broker).WithCallback(func(evt *mosquitto.Event) {
		if evt.Type == MOSQ_FLAG_EVENT_MESSAGE {
			if self.evtch != nil {
				self.evtch <- evt
			}
		}
	})
	if credentials := strings.TrimSpace(credentials); credentials != "" {
		userpass := strings.SplitN(credentials, ":", 2)
		cfg = cfg.WithCredentials(userpass[0], userpass[1])
	}

	// Connect to broker
	client, err := mosquitto.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("MQTT: %q: %w", broker, err)
	}

	// Home assistant
	ha, err := ha.New(topic, self.StateCallback)
	if err != nil {
		return nil, fmt.Errorf("Home Assistant: %q: %w", topic, err)
	}

	// Rotel amplifier
	rotel, err := rotel.NewWithConfig(rotel.Config{
		TTY: tty,
	})
	if err != nil {
		return nil, fmt.Errorf("Rotel: %q: %w", tty, err)
	}

	// Initialise logger
	self.Logger = log.New(os.Stderr, prefix+": ", log.LstdFlags)

	// Set app parameters
	self.qos = qos
	self.client = client
	self.ha = ha
	self.rotel = rotel
	self.id = fmt.Sprintf("%s_%s", "rotel", strings.TrimSpace(id))

	// Return success
	return self, nil
}

///////////////////////////////////////////////////////////////////////////////
// RUN

// runloop for the rotel app
func (self *App) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	var result error

	// Create channels for events and state changes
	self.evtch = make(chan *mosquitto.Event, 1)
	self.statech = make(chan StateChange, 2)
	rotelch := make(chan rotel.Event, 1)

	// Run rotel amplifier in background
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		if err := self.rotel.Run(ctx, rotelch); err != nil {
			result = errors.Join(result, err)
		}
	}(ctx)

	// Subscribe to the "status" topic to get online/offline messages
	self.topicStatusId = self.ha.TopicStatus()
	if _, err := self.client.Subscribe(self.topicStatusId, mosquitto.OptQoS(self.qos)); err != nil {
		return err
	}

	// Add a power button
	power, err := self.ha.AddPowerButton(self.id, "power")
	if err != nil {
		return err
	}
	if err := self.PublishComponent(power, true); err != nil {
		return err
	}

	// Add speaker A
	speakerA, err := self.ha.AddSpeaker(self.id, "speaker_a", "Speaker A")
	if err != nil {
		return err
	}
	if err := self.PublishComponent(speakerA, true); err != nil {
		return err
	}

	// Add speaker B
	speakerB, err := self.ha.AddSpeaker(self.id, "speaker_b", "Speaker B")
	if err != nil {
		return err
	}
	if err := self.PublishComponent(speakerB, true); err != nil {
		return err
	}

	// Add a volume slider
	volume, err := self.ha.AddVolume(self.id, "volume")
	if err != nil {
		return err
	}
	volume.(*ha.Volume).SetRange(rotel.VOLUME_MIN, rotel.VOLUME_MAX)
	if err := self.PublishComponent(volume, true); err != nil {
		return err
	}

	// Add tone sliders
	bass, err := self.ha.AddSlider(self.id, "bass", "Bass")
	if err != nil {
		return err
	}
	bass.(*ha.Slider).SetRange(rotel.TONE_MIN, rotel.TONE_MAX)
	if err := self.PublishComponent(bass, true); err != nil {
		return err
	}
	treble, err := self.ha.AddSlider(self.id, "treble", "Treble")
	if err != nil {
		return err
	}
	treble.(*ha.Slider).SetRange(rotel.TONE_MIN, rotel.TONE_MAX)
	if err := self.PublishComponent(treble, true); err != nil {
		return err
	}

	// Add input source
	source, err := self.ha.AddInput(self.id, "input", rotel.SOURCES)
	if err != nil {
		return err
	}
	if err := self.PublishComponent(source, true); err != nil {
		return err
	}

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
					// TODO: Get all latest rotel state
				}
			} else if err := self.ha.Command(evt.Topic, evt.Data); err != nil {
				log.Println("other event: ", evt)
			}
		case evt := <-self.statech:
			self.Logger.Println("publishing", string(evt.Data), "to", evt.Component.StateTopic())
			if _, err := self.client.Publish(evt.Component.StateTopic(), evt.Data); err != nil {
				return err
			}
			if evt.Component == power {
				if err := self.rotel.SetPower(string(evt.Data) == "ON"); err != nil {
					log.Println("error setting power:", err)
				}
			}
			if evt.Component == speakerA {
				if err := self.rotel.SetSpeaker(string(evt.Data) == "ON", "a"); err != nil {
					log.Println("error setting speaker A:", err)
				}
			}
			if evt.Component == speakerB {
				if err := self.rotel.SetSpeaker(string(evt.Data) == "ON", "b"); err != nil {
					log.Println("error setting speaker B:", err)
				}
			}
			if evt.Component == volume {
				if value, err := strconv.ParseUint(string(evt.Data), 10, 32); err != nil {
					log.Println("error parsing volume:", err)
				} else if err := self.rotel.SetVolume(uint(value)); err != nil {
					log.Println("error setting volume:", err)
				}
			}
			if evt.Component == source {
				if err := self.rotel.SetSource(string(evt.Data)); err != nil {
					log.Println("error setting source:", err)
				}
			}
			if evt.Component == bass {
				if value, err := strconv.ParseInt(string(evt.Data), 10, 32); err != nil {
					log.Println("error parsing bass:", err)
				} else if err := self.rotel.SetBass(int(value)); err != nil {
					log.Println("error setting bass:", err)
				}
			}
			if evt.Component == treble {
				if value, err := strconv.ParseInt(string(evt.Data), 10, 32); err != nil {
					log.Println("error parsing treble:", err)
				} else if err := self.rotel.SetTreble(int(value)); err != nil {
					log.Println("error setting treble:", err)
				}
			}
		case evt := <-rotelch:
			if evt.Err != nil {
				self.Logger.Println("rotel error", evt.Err)
			}
			if evt.Flag.Is(rotel.ROTEL_FLAG_MODEL) {
				self.Logger.Println("rotel model=", self.rotel.Model())
			}
			if evt.Flag.Is(rotel.ROTEL_FLAG_POWER) {
				if self.rotel.Power() {
					self.StateCallback(power, []byte("ON"))
				} else {
					self.StateCallback(power, []byte("OFF"))
				}
			}
			if evt.Flag.Is(rotel.ROTEL_FLAG_SPEAKER) {
				if self.rotel.SpeakerA() {
					self.StateCallback(speakerA, []byte("ON"))
				} else {
					self.StateCallback(speakerA, []byte("OFF"))
				}

				if self.rotel.SpeakerB() {
					self.StateCallback(speakerB, []byte("ON"))
				} else {
					self.StateCallback(speakerB, []byte("OFF"))
				}
			}
			if evt.Flag.Is(rotel.ROTEL_FLAG_VOLUME) {
				str := fmt.Sprintf("%d", self.rotel.Volume())
				self.StateCallback(volume, []byte(str))
			}
			if evt.Flag.Is(rotel.ROTEL_FLAG_SOURCE) {
				v := self.rotel.Source()
				self.StateCallback(source, []byte(v))
			}
			if evt.Flag.Is(rotel.ROTEL_FLAG_BASS) {
				str := fmt.Sprintf("%d", self.rotel.Bass())
				self.StateCallback(bass, []byte(str))
			}
			if evt.Flag.Is(rotel.ROTEL_FLAG_TREBLE) {
				str := fmt.Sprintf("%d", self.rotel.Treble())
				self.StateCallback(treble, []byte(str))
			}
		}
	}

	// Wait for rotel to finish
	wg.Wait()

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

	if component.SetState(string(data)) {
		self.Logger.Println("setting component state to", string(data), "for", component.StateTopic())
		payload := []byte(component.State())
		self.statech <- StateChange{component, payload}
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
	if self.ha != nil {
		str += fmt.Sprintf(" ha=%v", self.ha)
	}
	if self.rotel != nil {
		str += fmt.Sprintf(" rotel=%v", self.rotel)
	}
	str += fmt.Sprintf(" qos=%d", self.qos)
	if self.topic != "" {
		str += fmt.Sprintf(" topic=%v", self.topic)
	}
	return str + ">"
}
