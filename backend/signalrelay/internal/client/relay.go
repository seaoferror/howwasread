package client

import (
	"backend/common/proto"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// ErrServerUnavailable means the conversation server did not answer, its connection is dropped
var ErrServerUnavailable = errors.New("conversation server unavailable")

// RelayClient forwards WebRTC signals to the conversation server that holds the receivers' connections
type RelayClient interface {
	// Do sends req to the conversation server at ip.
	// When the server is unreachable the error wraps ErrServerUnavailable and the next Do dials again.
	Do(ctx context.Context, ip string, req *proto.RelaySignalRequest) error
}

type relayClient struct {
	dialOpts []grpc.DialOption
	mu       sync.Mutex
	conns    map[string]*grpc.ClientConn
}

// NewRelayClient dials conversation servers on port 50051 without TLS, opts are appended to the dial options
func NewRelayClient(opts ...grpc.DialOption) RelayClient {
	return &relayClient{
		dialOpts: append([]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, opts...),
		conns:    make(map[string]*grpc.ClientConn),
	}
}

func (r *relayClient) Do(ctx context.Context, ip string, req *proto.RelaySignalRequest) error {
	conn, err := r.conn(ip)
	if err != nil {
		return err
	}
	_, err = proto.NewSignalServiceClient(conn).RelaySignal(ctx, req)
	if code := status.Code(err); code == codes.Unavailable || code == codes.DeadlineExceeded {
		r.drop(ip, conn)
		return fmt.Errorf("%w: %w", ErrServerUnavailable, err)
	}
	return err
}

// conn returns the cached connection to ip, dialing it on first use
func (r *relayClient) conn(ip string) (*grpc.ClientConn, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if conn, ok := r.conns[ip]; ok {
		return conn, nil
	}
	conn, err := grpc.NewClient(ip+":50051", r.dialOpts...)
	if err != nil {
		return nil, err
	}
	r.conns[ip] = conn
	return conn, nil
}

func (r *relayClient) drop(ip string, conn *grpc.ClientConn) {
	r.mu.Lock()
	if r.conns[ip] == conn {
		delete(r.conns, ip)
	}
	r.mu.Unlock()
	if err := conn.Close(); err != nil {
		slog.Error("fail to close grpc client connection", "err", err)
	}
}
