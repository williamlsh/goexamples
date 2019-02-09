package userip

import (
	"context"
	"net"
	"net/http"
)

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

// userIPkey is the context key for the user IP address.  Its value of zero is
// arbitrary.  If this package defined other context keys, they would have
// different integer values.
const userIPKey key = 0

// FromRequest extracts a userIP value from an http.Request
func FromRequest(r *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, err
	}
	return []byte(ip), nil
}

// NewContext returns a new Context that carries a provided userIP value.
func NewContext(ctx context.Context, userIP net.IP) context.Context {
	return context.WithValue(ctx, userIPKey, userIP)
}

// FromContext extracts a userIP from a Context.
func FromContext(ctx context.Context) (net.IP, bool) {
	userIP, ok := ctx.Value(userIPKey).(net.IP)
	return userIP, ok
}
