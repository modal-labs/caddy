package caddyhttp

import (
	"strings"
	"testing"
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
	}
	std := cfg.std()
	if std.MaxConcurrentStreams != 1024 ||
		std.MaxReceiveBufferPerConnection != 64<<20 ||
		std.MaxReceiveBufferPerStream != 2<<20 {
		t.Fatalf("std conversion mismatch: %+v", std)
	}
}
