package rotel

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi"
	rotel "github.com/djthorpe/rotel"

	// Protocol Buffers
	pb "github.com/djthorpe/rotel/rpc/protobuf/rotel"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type evt struct {
	r *pb.RotelEvent
	c gopi.RPCClientConn
}

////////////////////////////////////////////////////////////////////////////////
// EVENT IMPLEMENTATION

func (this *evt) Source() gopi.Driver {
	if this != nil {
		return this.c
	} else {
		return nil
	}
}

func (this *evt) Name() string {
	return "RotelEvent"
}

func (this *evt) Type() rotel.EventType {
	if this != nil {
		return rotel.EventType(this.r.Type)
	} else {
		return rotel.EVENT_TYPE_NONE
	}
}

func (this *evt) State() rotel.RotelState {
	if this != nil {
		return protoToState(this.r.State)
	} else {
		return rotel.RotelState{}
	}
}

////////////////////////////////////////////////////////////////////////////////
// TO PROTOBUF

func protoFromEvent(evt rotel.RotelEvent) *pb.RotelEvent {
	if evt == nil {
		return nil
	} else {
		return &pb.RotelEvent{
			Type:  pb.RotelEvent_Type(evt.Type()),
			State: protoFromState(evt.State()),
		}
	}
}

func protoFromState(state rotel.RotelState) *pb.RotelState {
	return &pb.RotelState{
		Power:  protoFromPower(state.Power),
		Volume: protoFromVolume(state.Volume),
		Source: protoFromSource(state.Source),
	}
}

func protoFromPower(value rotel.Power) pb.RotelState_Power {
	switch value {
	case rotel.ROTEL_POWER_ON:
		return pb.RotelState_POWER_ON
	case rotel.ROTEL_POWER_STANDBY:
		return pb.RotelState_POWER_STANDBY
	case rotel.ROTEL_POWER_TOGGLE:
		return pb.RotelState_POWER_TOGGLE
	default:
		return pb.RotelState_POWER_NONE
	}
}

func protoFromVolume(value rotel.Volume) pb.RotelState_Volume {
	switch {
	case value == rotel.ROTEL_VOLUME_NONE:
		return pb.RotelState_VOLUME_NONE
	case value < rotel.ROTEL_VOLUME_MAX:
		return pb.RotelState_Volume(value)
	default:
		return pb.RotelState_Volume(rotel.ROTEL_VOLUME_MAX)
	}
}

func protoFromSource(value rotel.Source) pb.RotelState_Source {
	switch {
	case value == rotel.ROTEL_SOURCE_OTHER:
		return pb.RotelState_INPUT_NONE
	case value <= rotel.ROTEL_SOURCE_MAX:
		return pb.RotelState_Source(value)
	default:
		return pb.RotelState_INPUT_NONE
	}
}

func protoFromCommand(value rotel.Command) *pb.RotelCommand {
	switch {
	case value > rotel.ROTEL_COMMAND_NONE && value <= rotel.ROTEL_COMMAND_MAX:
		return &pb.RotelCommand{
			Command: pb.RotelCommand_Command(value),
		}
	default:
		return &pb.RotelCommand{}
	}
}

////////////////////////////////////////////////////////////////////////////////
// FROM PROTOBUF

func protoToState(proto *pb.RotelState) rotel.RotelState {
	if proto == nil {
		return rotel.RotelState{}
	} else {
		return rotel.RotelState{
			Power:   protoToPower(proto.Power),
			Volume:  protoToVolume(proto.Volume),
			Source:  protoToSource(proto.Source),
			Mute:    protoToMute(proto.Mute),
			Freq:    proto.Freq,
			Bypass:  proto.Bypass,
			Bass:    protoToTone(proto.Bass),
			Treble:  protoToTone(proto.Treble),
			Balance: protoToBalance(proto.Balance),
			Speaker: protoToSpeaker(proto.Speaker),
			Dimmer:  protoToDimmer(proto.Dimmer),
		}
	}
}

func protoToPower(value pb.RotelState_Power) rotel.Power {
	switch value {
	case pb.RotelState_POWER_ON:
		return rotel.ROTEL_POWER_ON
	case pb.RotelState_POWER_STANDBY:
		return rotel.ROTEL_POWER_STANDBY
	case pb.RotelState_POWER_TOGGLE:
		return rotel.ROTEL_POWER_TOGGLE
	default:
		return rotel.ROTEL_POWER_NONE
	}
}

func protoToVolume(value pb.RotelState_Volume) rotel.Volume {
	return rotel.Volume(value)
}

func protoToSource(value pb.RotelState_Source) rotel.Source {
	return rotel.Source(value)
}

func protoToMute(value pb.RotelState_Mute) rotel.Mute {
	switch value {
	case pb.RotelState_MUTE_ON:
		return rotel.ROTEL_MUTE_ON
	case pb.RotelState_MUTE_OFF:
		return rotel.ROTEL_MUTE_OFF
	default:
		return rotel.ROTEL_MUTE_NONE
	}
}

func protoToTone(value pb.RotelState_Tone) rotel.Tone {
	value_ := rotel.Tone(value)
	switch {
	case value_ == rotel.ROTEL_TONE_NONE:
		return value_
	case value_ < rotel.ROTEL_TONE_MAX:
		return value_
	case value_ == rotel.ROTEL_TONE_MAX:
		return value_
	default:
		return rotel.ROTEL_TONE_NONE
	}
}

func protoToSpeaker(value pb.RotelState_Speaker) rotel.Speaker {
	switch value {
	case pb.RotelState_SPEAKER_NONE:
		return rotel.Speaker{true, true}
	case pb.RotelState_SPEAKER_A:
		return rotel.Speaker{A: true}
	case pb.RotelState_SPEAKER_B:
		return rotel.Speaker{B: true}
	default:
		return rotel.Speaker{false, false}
	}
}

func protoToDimmer(value pb.RotelState_Dimmer) rotel.Dimmer {
	value_ := rotel.Dimmer(value)
	switch {
	case value_ == rotel.ROTEL_DIMMER_NONE:
		return value_
	case value_ < rotel.ROTEL_DIMMER_MAX:
		return value_
	case value_ == rotel.ROTEL_DIMMER_MAX:
		return value_
	default:
		return rotel.ROTEL_DIMMER_NONE
	}
}

func protoToBalance(value pb.RotelState_Balance) rotel.Balance {
	value_ := rotel.Balance(value)
	switch {
	case value_ == rotel.ROTEL_BALANCE_NONE:
		return value_
	case value_ == rotel.ROTEL_BALANCE_LEFT_MAX:
		return value_
	case value_ == rotel.ROTEL_BALANCE_RIGHT_MAX:
		return value_
	case value_ > rotel.ROTEL_BALANCE_LEFT_MAX:
		return value_
	case value_ < rotel.ROTEL_BALANCE_RIGHT_MAX:
		return value_
	default:
		return rotel.ROTEL_BALANCE_NONE
	}
}

func protoToCommand(value pb.RotelCommand_Command) rotel.Command {
	return rotel.Command(value)
}

func protoToEvent(value *pb.RotelEvent, conn gopi.RPCClientConn) rotel.RotelEvent {
	if value == nil {
		return nil
	} else {
		return &evt{value, conn}
	}
}
