package driver

import context "context"

// Driver describes a driver and the current car location.
type Driver struct {
	DriverID string
	Location string
}

// Interface exposed by the Driver service.
type Interface interface {
	FindNearest(ctx context.Context, location string) ([]Driver, error)
}
