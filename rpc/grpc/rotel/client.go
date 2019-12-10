/*
	Rotel RS232 Control
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package rotel

import (
	"context"
	"fmt"
	"io"
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

type Client struct {
	pb.RotelClient
	conn gopi.RPCClientConn
	event.Publisher
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewRotelClient(conn gopi.RPCClientConn) gopi.RPCClient {
	return &Client{pb.NewRotelClient(conn.(grpc.GRPCClientConn).GRPCConn()), conn, event.Publisher{}}
}

func (this *Client) NewContext(timeout time.Duration) context.Context {
	if timeout == 0 {
		timeout = this.conn.Timeout()
	}
	if timeout == 0 {
		return context.Background()
	} else {
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		return ctx
	}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Client) Conn() gopi.RPCClientConn {
	return this.conn
}

////////////////////////////////////////////////////////////////////////////////
// CALLS

func (this *Client) Ping() error {
	this.conn.Lock()
	defer this.conn.Unlock()

	// Perform ping
	if _, err := this.RotelClient.Ping(this.NewContext(0), &empty.Empty{}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) Get() (rotel.RotelState, error) {
	this.conn.Lock()
	defer this.conn.Unlock()

	// Get state
	if state, err := this.RotelClient.Get(this.NewContext(0), &empty.Empty{}); err != nil {
		return rotel.RotelState{}, err
	} else {
		return protoToState(state), nil
	}

}

func (this *Client) Set(state rotel.RotelState) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	// Set state
	if _, err := this.RotelClient.Set(this.NewContext(0), protoFromState(state)); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) Send(command rotel.Command) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	if _, err := this.RotelClient.Send(this.NewContext(0), protoFromCommand(command)); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Client) StreamEvents(ctx context.Context) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	stream, err := this.RotelClient.StreamEvents(ctx, &empty.Empty{})
	if err != nil {
		return err
	}

	// Errors channel receives errors from recv
	errors := make(chan error)
	defer close(errors)

	// Receive messages in the background
	go func() {
	FOR_LOOP:
		for {
			if evt_, err := stream.Recv(); err == io.EOF {
				break FOR_LOOP
			} else if err != nil {
				errors <- err
				break FOR_LOOP
			} else if evt := protoToEvent(evt_, this.conn); evt != nil {
				this.Emit(evt)
			}
		}
		fmt.Println("StreamEvents: ENDED GOROUTINE")
	}()

	// Continue until error or io.EOF is returned
	for {
		select {
		case err := <-errors:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Client) String() string {
	return fmt.Sprintf("<rpc.service.rotel.Client>{ conn=%v }", this.conn)
}
