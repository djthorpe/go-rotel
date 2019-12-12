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
		Power:   protoFromPower(state.Power),
		Volume:  protoFromVolume(state.Volume),
		Mute:    protoFromMute(state.Mute),
		Source:  protoFromSource(state.Source),
		Freq:    state.Freq,
		Bypass:  protoFromBypass(state.Bypass),
		Treble:  protoFromTone(state.Treble),
		Bass:    protoFromTone(state.Bass),
		Balance: protoFromBalance(state.Balance),
		Speaker: protoFromSpeaker(state.Speaker),
		Dimmer:  protoFromDimmer(state.Dimmer),
	}
}

func protoFromPower(value rotel.Power) pb.RotelState_Power {
	switch value {
	case rotel.ROTEL_POWER_ON:
		return pb.RotelState_POWER_ON
	case rotel.ROTEL_POWER_STANDBY:
		return pb.RotelState_POWER_STANDBY
	default:
		return pb.RotelState_POWER_NONE
	}
}

func protoFromVolume(value rotel.Volume) pb.RotelState_Volume {
	switch {
	case value > rotel.ROTEL_VOLUME_MIN && value < rotel.ROTEL_VOLUME_MAX:
		return pb.RotelState_Volume(value)
	default:
		return pb.RotelState_Volume(rotel.ROTEL_VOLUME_NONE)
	}
}

func protoFromSource(value rotel.Source) pb.RotelState_Source {
	switch {
	case value == rotel.ROTEL_SOURCE_OTHER:
		return pb.RotelState_SOURCE_OTHER
	case value > rotel.ROTEL_SOURCE_NONE && value <= rotel.ROTEL_SOURCE_MAX:
		return pb.RotelState_Source(value)
	default:
		return pb.RotelState_SOURCE_NONE
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

func protoFromMute(value rotel.Mute) pb.RotelState_Mute {
	switch value {
	case rotel.ROTEL_MUTE_ON:
		return pb.RotelState_MUTE_ON
	case rotel.ROTEL_MUTE_OFF:
		return pb.RotelState_MUTE_OFF
	default:
		return pb.RotelState_MUTE_NONE
	}
}

func protoFromBypass(value rotel.Bypass) pb.RotelState_Bypass {
	switch value {
	case rotel.ROTEL_BYPASS_ON:
		return pb.RotelState_BYPASS_ON
	case rotel.ROTEL_BYPASS_OFF:
		return pb.RotelState_BYPASS_OFF
	default:
		return pb.RotelState_BYPASS_NONE
	}
}

func protoFromTone(value rotel.Tone) pb.RotelState_Tone {
	switch {
	case value == rotel.ROTEL_TONE_OFF:
		return pb.RotelState_TONE_OFF
	case value >= rotel.ROTEL_TONE_MIN && value <= rotel.ROTEL_TONE_MAX:
		return pb.RotelState_Tone(value)
	default:
		return pb.RotelState_TONE_NONE
	}
}

func protoFromBalance(value rotel.Balance) pb.RotelState_Balance {
	switch {
	case value == rotel.ROTEL_BALANCE_OFF:
		return pb.RotelState_BALANCE_OFF
	case value >= rotel.ROTEL_BALANCE_LEFT_MAX && value <= rotel.ROTEL_BALANCE_RIGHT_MAX:
		return pb.RotelState_Balance(value)
	default:
		return pb.RotelState_BALANCE_NONE
	}
}

func protoFromSpeaker(value rotel.Speaker) pb.RotelState_Speaker {
	switch value {
	case rotel.ROTEL_SPEAKER_OFF:
		return pb.RotelState_SPEAKER_OFF
	case rotel.ROTEL_SPEAKER_A:
		return pb.RotelState_SPEAKER_A
	case rotel.ROTEL_SPEAKER_B:
		return pb.RotelState_SPEAKER_B
	case rotel.ROTEL_SPEAKER_ALL:
		return pb.RotelState_SPEAKER_ALL
	default:
		return pb.RotelState_SPEAKER_NONE
	}
}

func protoFromDimmer(value rotel.Dimmer) pb.RotelState_Dimmer {
	switch {
	case value == rotel.ROTEL_DIMMER_OFF:
		return pb.RotelState_DIMMER_OFF
	case value >= rotel.ROTEL_DIMMER_MIN && value <= rotel.ROTEL_DIMMER_MAX:
		return pb.RotelState_Dimmer(value)
	default:
		return pb.RotelState_DIMMER_NONE
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
			Mute:    protoToMute(proto.Mute),
			Source:  protoToSource(proto.Source),
			Freq:    proto.Freq,
			Bypass:  protoToBypass(proto.Bypass),
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
	default:
		return rotel.ROTEL_POWER_NONE
	}
}

func protoToVolume(value pb.RotelState_Volume) rotel.Volume {
	return rotel.Volume(value)
}

func protoToBypass(value pb.RotelState_Bypass) rotel.Bypass {
	return rotel.Bypass(value)
}

func protoToSource(value pb.RotelState_Source) rotel.Source {
	return rotel.Source(value)
}

func protoToMute(value pb.RotelState_Mute) rotel.Mute {
	return rotel.Mute(value)
}

func protoToSpeaker(value pb.RotelState_Speaker) rotel.Speaker {
	return rotel.Speaker(value)
}

func protoToTone(value pb.RotelState_Tone) rotel.Tone {
	return rotel.Tone(value)
}

func protoToDimmer(value pb.RotelState_Dimmer) rotel.Dimmer {
	return rotel.Dimmer(value)
}

func protoToBalance(value pb.RotelState_Balance) rotel.Balance {
	return rotel.Balance(value)
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
