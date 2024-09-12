package metadata

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

type MetadataExtractor struct{}

func NewMetadataExtractor() *MetadataExtractor {
	return &MetadataExtractor{}
}

func (me *MetadataExtractor) Extract(ctx context.Context) Metadata {
	md := Metadata{}

	if grpcMd, ok := metadata.FromIncomingContext(ctx); ok {
		md.UserAgent = me.extractUserAgent(grpcMd)
		md.ClientIP = me.extractClientIP(grpcMd)
	}

	// If ClientIP is not found in metadata, try to get it from peer info
	if md.ClientIP == "" {
		if p, ok := peer.FromContext(ctx); ok {
			md.ClientIP = p.Addr.String()
		}
	}

	return md
}

func (me *MetadataExtractor) extractUserAgent(md metadata.MD) string {
	if ua := md.Get("User-Agent"); len(ua) > 0 {
		return ua[0]
	}
	if ua := md.Get("user-agent"); len(ua) > 0 {
		return ua[0]
	}
	return ""
}

func (me *MetadataExtractor) extractClientIP(md metadata.MD) string {
	if xff := md.Get("X-Forwarded-For"); len(xff) > 0 {
		ips := strings.Split(xff[0], ",")
		return strings.TrimSpace(ips[0])
	}
	if xrip := md.Get("X-Real-IP"); len(xrip) > 0 {
		return xrip[0]
	}
	if ra := md.Get("RemoteAddr"); len(ra) > 0 {
		return strings.Split(ra[0], ":")[0]
	}
	return ""
}
