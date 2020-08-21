package systemd

import (
	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/pkg/errors"
)

func StartUnit(name string) error {
	bus, err := dbus.New()
	if err != nil {
		return err
	}
	defer bus.Close()

	_, err = bus.StartUnit(name, "fail", nil)
	return err
}

func StopUnit(name string) error {
	bus, err := dbus.New()
	if err != nil {
		return err
	}
	defer bus.Close()

	_, err = bus.StopUnit(name, "fail", nil)
	return err
}

func UnitStatus(name string) (*dbus.UnitStatus, error) {
	bus, err := dbus.New()
	if err != nil {
		return nil, err
	}
	defer bus.Close()

	services, err := bus.ListUnitsByNames([]string{name})
	if err != nil {
		return nil, err
	}
	if len(services) != 1 {
		return nil, errors.Errorf("cannot find unit: %q", name)
	}
	return &services[0], nil
}

const (
	ActiveStateActive   = "active"
	ActiveStateInactive = "inactive"
	SubStateRunning     = "running"
	SubStateDead        = "dead"
)

func UnitIsReady(status *dbus.UnitStatus) bool {
	return status.ActiveState == ActiveStateActive && status.SubState == SubStateRunning
}
