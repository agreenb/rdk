// Package client provides a multi-faceted approach for connecting to a server.
package client

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/edaniels/golog"

	"go.viam.com/core/rpc"
	"go.viam.com/core/rpc/dialer"
	rpcwebrtc "go.viam.com/core/rpc/webrtc"
)

// DialOptions are extra dial time options.
type DialOptions struct {
	// Insecure determines if the RPC connection is TLS based.
	Insecure bool

	// Signaling server specifies the signaling server to
	// contact on behalf of this client for WebRTC communications.
	SignalingServer string
}

// Dial attempts to make the most convenient connection to the given address. It first tries a direct
// connection if the address is an IP. It next tries to connect to the local version of the host followed
// by a WebRTC brokered connection.
func Dial(ctx context.Context, address string, opts DialOptions, logger golog.Logger) (dialer.ClientConn, error) {
	var host string
	var port string
	if strings.ContainsRune(address, ':') {
		var err error
		host, port, err = net.SplitHostPort(address)
		if err != nil {
			return nil, err
		}
	} else {
		host = address
	}

	if addr := net.ParseIP(host); addr == nil {
		localHost := fmt.Sprintf("local.%s", host)
		if _, err := lookupHost(ctx, localHost); err == nil {
			localAddress := localHost
			if port != "" {
				localAddress = fmt.Sprintf("%s:%s", localHost, port)
			}
			// TODO(erd): This needs to authenticate the server so we don't have a confused
			// deputy.
			if conn, err := dialer.DialDirectGRPC(ctx, localAddress, true); err == nil {
				logger.Debugw("connected directly via local host", "address", localAddress)
				return conn, nil
			} else if ctx.Err() != nil {
				return nil, ctx.Err()
			}
		} else if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	if opts.SignalingServer != "" {
		webrtcAddress := rpc.HostURI(opts.SignalingServer, address)
		conn, err := rpcwebrtc.Dial(ctx, webrtcAddress, opts.Insecure, logger)
		if err != nil && !errors.Is(err, rpcwebrtc.ErrNoSignaler) {
			return nil, err
		}
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		logger.Debug("connected via WebRTC")
		return conn, nil
	}

	conn, err := dialer.DialDirectGRPC(ctx, address, opts.Insecure)
	if err != nil {
		return nil, err
	}
	logger.Debugw("connected directly", "address", address)
	return conn, nil

}

func lookupHost(ctx context.Context, host string) (addrs []string, err error) {
	if ctxResolver := dialer.ContextResolver(ctx); ctxResolver != nil {
		return ctxResolver.LookupHost(ctx, host)
	}
	return net.DefaultResolver.LookupHost(ctx, host)
}
