package metadata

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func TestMetadataExtractor(t *testing.T) {
	t.Parallel()

	extractor := NewMetadataExtractor()
	assert.NotNil(t, extractor, "NewMetadataExtractor should return a non-nil extractor")
}

func TestMetadataExtractor_Extract(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupContext   func() context.Context
		expectedResult Metadata
	}{
		{
			name: "Extract from gRPC metadata",
			setupContext: func() context.Context {
				ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
					"User-Agent": "test-agent",
					"X-Real-IP":  "192.168.1.1",
				}))
				return ctx
			},
			expectedResult: Metadata{
				UserAgent: "test-agent",
				ClientIP:  "192.168.1.1",
			},
		},
		{
			name: "Extract from peer info",
			setupContext: func() context.Context {
				ctx := peer.NewContext(context.Background(), &peer.Peer{
					Addr: &net.TCPAddr{
						IP:   net.ParseIP("10.0.0.1"),
						Port: 1234,
					},
				})
				return ctx
			},
			expectedResult: Metadata{
				UserAgent: "",
				ClientIP:  "10.0.0.1:1234",
			},
		},
		{
			name: "Extract with X-Forwarded-For",
			setupContext: func() context.Context {
				ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
					"X-Forwarded-For": "203.0.113.195, 70.41.3.18, 150.172.238.178",
				}))
				return ctx
			},
			expectedResult: Metadata{
				UserAgent: "",
				ClientIP:  "203.0.113.195",
			},
		},
		{
			name: "Extract with lowercase user-agent",
			setupContext: func() context.Context {
				ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
					"user-agent": "lower-case-agent",
				}))
				return ctx
			},
			expectedResult: Metadata{
				UserAgent: "lower-case-agent",
				ClientIP:  "",
			},
		},
		{
			name: "Extract with RemoteAddr",
			setupContext: func() context.Context {
				ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
					"RemoteAddr": "192.168.0.1:5678",
				}))
				return ctx
			},
			expectedResult: Metadata{
				UserAgent: "",
				ClientIP:  "192.168.0.1",
			},
		},
		{
			name: "Extract with empty context",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedResult: Metadata{
				UserAgent: "",
				ClientIP:  "",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			extractor := NewMetadataExtractor()
			ctx := tt.setupContext()
			result := extractor.Extract(ctx)

			assert.Equal(t, tt.expectedResult.UserAgent, result.UserAgent, "UserAgent mismatch")
			assert.Equal(t, tt.expectedResult.ClientIP, result.ClientIP, "ClientIP mismatch")
		})
	}
}
