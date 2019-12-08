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
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	grpc "github.com/djthorpe/gopi-rpc/sys/grpc"
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
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewRotelClient(conn gopi.RPCClientConn) gopi.RPCClient {
	return &Client{pb.NewRotelClient(conn.(grpc.GRPCClientConn).GRPCConn()), conn}
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

func (this *Client) Set(rotel.RotelState) error {
	this.conn.Lock()
	defer this.conn.Unlock()

	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Client) String() string {
	return fmt.Sprintf("<rpc.service.rotel.Client>{ conn=%v }", this.conn)
}
