package caddyhttp

import (
	"strings"
	"testing"
	"time"

	"github.com/caddyserver/caddy/v2"
)

func TestHTTP2ConfigValidate(t *testing.T) {
	for _, tc := range []struct {
		name    string
		cfg     HTTP2Config
		wantErr string
	}{
		{
			name: "zero values are valid",
			cfg:  HTTP2Config{},
		},
		{
			name: "typical values are valid",
			cfg: HTTP2Config{
				MaxConcurrentStreams:          1024,
				MaxReceiveBufferPerConnection: 64 << 20,
				MaxReceiveBufferPerStream:     2 << 20,
			},
		},
		{
			name:    "negative streams",
			cfg:     HTTP2Config{MaxConcurrentStreams: -1},
			wantErr: "max_concurrent_streams",
		},
		{
			name: "connection buffer below initial window size",
			// x/net silently reverts values below 65535 to its default,
			// so they must be rejected here.
			cfg:     HTTP2Config{MaxReceiveBufferPerConnection: 65534},
			wantErr: "max_receive_buffer_per_connection",
		},
		{
			name:    "negative stream buffer",
			cfg:     HTTP2Config{MaxReceiveBufferPerStream: -1},
			wantErr: "max_receive_buffer_per_stream",
		},
		{
			name: "ping and write timeouts are valid",
			cfg: HTTP2Config{
				SendPingTimeout:  caddy.Duration(30 * time.Second),
				PingTimeout:      caddy.Duration(15 * time.Second),
				WriteByteTimeout: caddy.Duration(30 * time.Second),
			},
		},
		{
			name:    "negative send ping timeout",
			cfg:     HTTP2Config{SendPingTimeout: caddy.Duration(-time.Second)},
			wantErr: "send_ping_timeout",
		},
		{
			name:    "negative ping timeout",
			cfg:     HTTP2Config{PingTimeout: caddy.Duration(-time.Second)},
			wantErr: "ping_timeout",
		},
		{
			name:    "negative write byte timeout",
			cfg:     HTTP2Config{WriteByteTimeout: caddy.Duration(-time.Second)},
			wantErr: "write_byte_timeout",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.validate()
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestHTTP2ConfigStd(t *testing.T) {
	cfg := HTTP2Config{
		MaxConcurrentStreams:          1024,
		MaxReceiveBufferPerConnection: 64 << 20,
		MaxReceiveBufferPerStream:     2 << 20,
		SendPingTimeout:               caddy.Duration(30 * time.Second),
		PingTimeout:                   caddy.Duration(15 * time.Second),
		WriteByteTimeout:              caddy.Duration(30 * time.Second),
	}
	std := cfg.std()
	if std.MaxConcurrentStreams != 1024 ||
		std.MaxReceiveBufferPerConnection != 64<<20 ||
		std.MaxReceiveBufferPerStream != 2<<20 ||
		std.SendPingTimeout != 30*time.Second ||
		std.PingTimeout != 15*time.Second ||
		std.WriteByteTimeout != 30*time.Second {
		t.Fatalf("std conversion mismatch: %+v", std)
	}
}
