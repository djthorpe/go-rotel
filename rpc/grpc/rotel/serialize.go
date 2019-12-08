package rotel

import (
	rotel "github.com/djthorpe/rotel"
	pb "github.com/djthorpe/rotel/rpc/protobuf/rotel"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type evt struct {
	r *pb.RotelEvent
}

////////////////////////////////////////////////////////////////////////////////
// TO PROTOBUF

func protoFromEvent(evt rotel.RotelEvent) *pb.RotelEvent {
	return &pb.RotelEvent{}
}

func protoFromState(power rotel.Power, volume rotel.Volume, input rotel.Source) *pb.RotelState {
	return &pb.RotelState{
		Power:  protoFromPower(power),
		Volume: protoFromVolume(volume),
		Input:  protoFromSource(input),
	}
}

func protoFromPower(value rotel.Power) pb.RotelState_Power {
	switch value {
	case rotel.ROTEL_POWER_ON:
		return pb.RotelState_POWER_ON
	case rotel.ROTEL_POWER_STANDY:
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

////////////////////////////////////////////////////////////////////////////////
// FROM PROTOBUF

func protoToState(proto *pb.RotelState) rotel.RotelState {
	if proto == nil {
		return rotel.RotelState{}
	} else {
		return rotel.RotelState{
			Power:  protoToPower(proto.Power),
			Volume: protoToVolume(proto.Volume),
			Source: protoToSource(proto.Input),
		}
	}
}

func protoToPower(value pb.RotelState_Power) rotel.Power {
	switch value {
	case pb.RotelState_POWER_ON:
		return rotel.ROTEL_POWER_ON
	case pb.RotelState_POWER_STANDBY:
		return rotel.ROTEL_POWER_STANDY
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
