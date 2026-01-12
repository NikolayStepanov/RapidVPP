package vpp

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/NikolayStepanov/RapidVPP/pkg/logger"
	"go.fd.io/govpp"
	"go.fd.io/govpp/api"
	"go.fd.io/govpp/binapi/memclnt"
	"go.fd.io/govpp/core"
	"go.uber.org/zap"
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

func (c *Client) NewStream(ctx context.Context) (api.Stream, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, fmt.Errorf("vpp client closed")
	}

	if c.conn == nil {
		return nil, fmt.Errorf("vpp not connected")
	}

	stream, err := c.conn.NewStream(ctx,
		core.WithRequestSize(50),
		core.WithReplySize(50),
		core.WithReplyTimeout(10*time.Second))
	if err != nil {
		return nil, fmt.Errorf("create stream failed: %w", err)
	}

	return stream, nil
}

func (c *Client) Do(ctx context.Context, fn func(stream api.Stream) error) error {
	stream, err := c.NewStream(ctx)
	if err != nil {
		return err
	}
	defer stream.Close()

	return fn(stream)
}

func (c *Client) DoWithTimeout(timeout time.Duration, fn func(stream api.Stream) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return c.Do(ctx, fn)
}

func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}

	c.closed = true
	if c.conn != nil {
		c.conn.Disconnect()
	}
}

func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return !c.closed && c.conn != nil
}

func (c *Client) SendMultiRequest(ctx context.Context, request api.Message) ([]api.Message, error) {
	stream, err := c.NewStream(ctx)
	if err != nil {
		return nil, fmt.Errorf("create stream: %w", err)
	}
	defer stream.Close()

	if err := stream.SendMsg(request); err != nil {
		logger.Error("send request failed", zap.Error(err))
		return nil, fmt.Errorf("send request: %w", err)
	}
	err = stream.SendMsg(&memclnt.ControlPing{})
	if err != nil {
		logger.Error("send ping failed", zap.Error(err))
		return nil, err
	}

	var (
		messages []api.Message
		errors   []error
	)

	for {
		message, err := stream.RecvMsg()
		if err != nil {
			logger.Error("recv message failed", zap.Error(err))
			return nil, fmt.Errorf("receive message: %w", err)
		}
		switch message.(type) {
		case *memclnt.ControlPingReply:
			if len(errors) > 0 {
				return messages, fmt.Errorf("completed with %d errors: %v", len(errors), errors)
			}
			return messages, nil
		default:
			// TODO: Retval check via reflection is suboptimal
			if retval := getRetval(message); retval != 0 {
				errors = append(errors, fmt.Errorf("retval=%d for %T", retval, message))
				continue
			}
			messages = append(messages, message)
		}
	}
}

func getRetval(msg api.Message) int32 {
	if m, ok := msg.(interface{ Retval() int32 }); ok {
		return m.Retval()
	}
	v := reflect.ValueOf(msg)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		f := v.FieldByName("Retval")
		if f.IsValid() && f.Kind() == reflect.Int32 {
			return int32(f.Int())
		}
	}
	return 0
}

func (c *Client) SendMultiRequestWithTimeout(timeout time.Duration, request api.Message) ([]api.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c.SendMultiRequest(ctx, request)
}

func Dump[T any](
	ctx context.Context,
	client *Client,
	request api.Message,
	converter func(api.Message) (T, bool),
) ([]T, error) {
	var results []T

	messages, err := client.SendMultiRequest(ctx, request)
	if err != nil {
		logger.Debug("dump request failed", zap.Error(err))
		return nil, err
	}
	for _, message := range messages {
		if item, ok := converter(message); ok {
			logger.Debug("dump message received", zap.Any("message", item))
			results = append(results, item)
		}
	}

	return results, nil
}

func DumpWithTimeout[T any](
	client *Client,
	timeout time.Duration,
	request api.Message,
	converter func(api.Message) (T, bool),
) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return Dump(ctx, client, request, converter)
}

func DoRequest[Req, Resp api.Message](client *Client, ctx context.Context, req Req) (Resp, error) {
	var emptyResp Resp
	var reply Resp

	err := client.Do(ctx, func(stream api.Stream) error {
		if err := stream.SendMsg(req); err != nil {
			return fmt.Errorf("send request: %w", err)
		}

		msg, err := stream.RecvMsg()
		if err != nil {
			return fmt.Errorf("receive reply: %w", err)
		}

		var ok bool
		reply, ok = msg.(Resp)
		if !ok {
			return fmt.Errorf("unexpected message type: %T, expected %T", msg, reply)
		}
		if retval := getRetval(reply); retval != 0 {
			return fmt.Errorf("retval=%d for %T", retval, reply)
		}
		return nil
	})

	if err != nil {
		return emptyResp, err
	}

	return reply, nil
}
