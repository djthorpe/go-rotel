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
