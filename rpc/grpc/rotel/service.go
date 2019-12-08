package rotel

import (
	"context"
	"fmt"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	grpc "github.com/djthorpe/gopi-rpc/sys/grpc"
	event "github.com/djthorpe/gopi/util/event"
	rotel "github.com/djthorpe/rotel"

	// Protocol buffers
	pb "github.com/djthorpe/rotel/rpc/protobuf/rotel"
	empty "github.com/golang/protobuf/ptypes/empty"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Service struct {
	Server gopi.RPCServer
	Rotel  rotel.Rotel
}

type service struct {
	log   gopi.Logger
	rotel rotel.Rotel

	event.Publisher
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open the server
func (config Service) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<grpc.service.rotel>Open{ %+v }", config)

	if config.Server == nil || config.Rotel == nil {
		return nil, gopi.ErrBadParameter
	}

	this := new(service)
	this.log = log
	this.rotel = config.Rotel

	// Register service with GRPC server
	pb.RegisterRotelServer(config.Server.(grpc.GRPCServer).GRPCServer(), this)

	// Success
	return this, nil
}

func (this *service) Close() error {
	this.log.Debug("<grpc.service.rotel>Close{ %v }", this.rotel)

	// Close event channel
	this.Publisher.Close()

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RPCService IMPLEMENTATION

func (this *service) CancelRequests() error {
	this.log.Debug("<grpc.service.rotel>CancelRequests{}")

	// Put empty event onto the channel to indicate any on-going
	// requests should be ended
	this.Emit(event.NullEvent)

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *service) String() string {
	return fmt.Sprintf("grpc.service.rotel{ %v }", this.rotel)
}

////////////////////////////////////////////////////////////////////////////////
// RPC Methods

func (this *service) Ping(context.Context, *empty.Empty) (*empty.Empty, error) {
	this.log.Debug("<grpc.service.rotel.Ping>{ }")
	return &empty.Empty{}, nil
}

func (this *service) Query(context.Context, *empty.Empty) (*pb.RotelState, error) {
	this.log.Debug("<grpc.service.rotel.Query>{ }")
	return protoFromState(
		this.rotel.Model(),
		this.rotel.Power(),
		this.rotel.Volume(),
		this.rotel.Input(),
	), nil
}

// Stream events
func (this *service) StreamEvents(_ *empty.Empty, stream pb.Rotel_StreamEventsServer) error {
	this.log.Debug2("<grpc.service.rotel.StreamEvents>{}")

	// Subscribe to channel for incoming events, and continue until cancel request is received, send
	// empty events occasionally to ensure the channel is still alive
	events := this.rotel.Subscribe()
	cancel := this.Subscribe()
	ticker := time.NewTicker(time.Second)

FOR_LOOP:
	for {
		select {
		case evt := <-events:
			if evt == nil {
				break FOR_LOOP
			} else if evt_, ok := evt.(rotel.RotelEvent); ok {
				if err := stream.Send(protoFromEvent(evt_)); err != nil {
					this.log.Warn("StreamEvents: %v", err)
					break FOR_LOOP
				}
			} else {
				this.log.Warn("StreamEvents: Ignoring event: %v", evt)
			}
		case <-ticker.C:
			if err := stream.Send(&pb.RotelEvent{}); err != nil {
				this.log.Warn("StreamEvents: %v", err)
				break FOR_LOOP
			}
		case <-cancel:
			break FOR_LOOP
		}
	}

	// Stop ticker, unsubscribe from events
	ticker.Stop()
	this.rotel.Unsubscribe(events)
	this.Unsubscribe(cancel)

	this.log.Debug2("StreamEvents: Ended")

	// Return success
	return nil
}
