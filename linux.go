//go:build linux

package wgiface

import (
	"context"
	"fmt"
	"os/exec"
)

func (i *Interface) createInter(ctx context.Context) error {
	_, err := execCmd(exec.CommandContext(ctx, "ip", "link", "add", i.name, "type", "wireguard"))
	return err
}

func (i *Interface) addInterAddr(ctx context.Context) error {
	for _, ipnet := range i.addr {
		if ipnet.IP.To4() != nil {
			if _, err := execCmd(exec.CommandContext(ctx, "ip", "-4", "address", "add", ipnet.String(), "dev", i.name)); err != nil {
				return err
			}
		} else if ipnet.IP.To16() != nil {
			if _, err := execCmd(exec.CommandContext(ctx, "ip", "-6", "address", "add", ipnet.String(), "dev", i.name)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (i *Interface) setMTU(ctx context.Context) error {
	if i.mtu == nil {
		return nil
	}
	_, err := execCmd(exec.CommandContext(ctx, "ip", "link", "set", "mtu", fmt.Sprint(*i.mtu), "dev", i.name))
	return err
}

func (i *Interface) startInter(ctx context.Context) error {
	if _, err := execCmd(exec.CommandContext(ctx, "ip", "link", "set", "up", i.name)); err != nil {
		return fmt.Errorf("failed to set mtu and start for wireguard interface: %w", err)
	}
	return nil
}

func (i *Interface) stopInter(ctx context.Context) error {
	_, err := execCmd(exec.CommandContext(ctx, "ip", "link", "delete", i.name))
	return err
}
