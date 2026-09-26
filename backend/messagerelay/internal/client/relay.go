package client

import (
	"backend/common/proto"
	"context"
	"log/slog"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RelayClient forwards messages to the chat server that holds the receivers' websocket connections
type RelayClient interface {
	// Do sends req to the chat server at ip and returns the receivers it could not deliver to.
	// A failed call drops the connection, the next Do to that ip dials again.
	Do(ctx context.Context, ip string, req *proto.RelayMessagingRequest) (pushToIds [][]byte, err error)
}

type relayClient struct {
	dialOpts []grpc.DialOption
	mu       sync.Mutex
	conns    map[string]*grpc.ClientConn
}

// NewRelayClient dials chat servers on port 50051 without TLS, opts are appended to the dial options
func NewRelayClient(opts ...grpc.DialOption) RelayClient {
	return &relayClient{
		dialOpts: append([]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, opts...),
		conns:    make(map[string]*grpc.ClientConn),
	}
}

func (r *relayClient) Do(ctx context.Context, ip string, req *proto.RelayMessagingRequest) ([][]byte, error) {
	conn, err := r.conn(ip)
	if err != nil {
		return nil, err
	}
	res, err := proto.NewMessagingServiceClient(conn).RelayMessaging(ctx, req)
	if err != nil {
		r.drop(ip, conn)
		return nil, err
	}
	return res.GetPushToIds(), nil
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
