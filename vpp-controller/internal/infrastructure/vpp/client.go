package vpp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.fd.io/govpp"
	"go.fd.io/govpp/api"
	"go.fd.io/govpp/core"
)

type Client struct {
	conn   *core.Connection
	mu     sync.RWMutex
	closed bool
}

func NewClient(socketPath string) (*Client, error) {
	conn, connEv, err := govpp.AsyncConnect(socketPath,
		3,
		time.Second)
	if err != nil {
		return nil, fmt.Errorf("async connect failed: %w", err)
	}

	select {
	case e := <-connEv:
		if e.State != core.Connected {
			conn.Disconnect()
			return nil, fmt.Errorf("connection failed: %v", e.Error)
		}
	case <-time.After(10 * time.Second):
		conn.Disconnect()
		return nil, fmt.Errorf("connection timeout")
	}

	return &Client{
		conn: conn,
	}, nil
}

func (v *Client) NewStream(ctx context.Context) (api.Stream, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.closed {
		return nil, fmt.Errorf("vpp client closed")
	}

	if v.conn == nil {
		return nil, fmt.Errorf("vpp not connected")
	}

	stream, err := v.conn.NewStream(ctx,
		core.WithRequestSize(50),
		core.WithReplySize(50),
		core.WithReplyTimeout(5*time.Second))
	if err != nil {
		return nil, fmt.Errorf("create stream failed: %w", err)
	}

	return stream, nil
}

func (v *Client) Do(ctx context.Context, fn func(stream api.Stream) error) error {
	stream, err := v.NewStream(ctx)
	if err != nil {
		return err
	}
	defer stream.Close()

	return fn(stream)
}

func (v *Client) DoWithTimeout(timeout time.Duration, fn func(stream api.Stream) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return v.Do(ctx, fn)
}

func (v *Client) Close() {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.closed {
		return
	}

	v.closed = true
	if v.conn != nil {
		v.conn.Disconnect()
	}
}

func (v *Client) IsConnected() bool {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return !v.closed && v.conn != nil
}
