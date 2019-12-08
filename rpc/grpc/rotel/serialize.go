package rotel

import (
	rotel "github.com/djthorpe/rotel"
	pb "github.com/djthorpe/rotel/rpc/protobuf/rotel"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type evt struct {
	r *pb.Event
}

////////////////////////////////////////////////////////////////////////////////
// TO PROTOBUF

func protoFromEvent(evt rotel.RotelEvent) *pb.Event {
	return &pb.Event{}
}
