package gateway

import "net"

// TODO(hacdias): copied from ../hostname.go

// I dont like this here. Needed to move it here so it is the same exact type.
type RequestContextKey string

func stripPort(hostname string) string {
	host, _, err := net.SplitHostPort(hostname)
	if err == nil {
		return host
	}
	return hostname
}
