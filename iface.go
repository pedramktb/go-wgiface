package wgiface

import (
	"context"
	"errors"
	"fmt"
	"net"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Interface struct {
	wgClient   *wgctrl.Client
	name       string
	privateKey *wgtypes.Key
	addr       []net.IPNet
	listenPort *int
	mtu        *int
	fwmark     *int
	peers      []wgtypes.PeerConfig
}
type Option func(*Interface)

func WithPrivateKey(key wgtypes.Key) Option {
	return func(i *Interface) {
		i.privateKey = &key
	}
}

func WithAddresses(addr []net.IPNet) Option {
	return func(i *Interface) {
		i.addr = addr[:]
	}
}

func WithListenPort(listenPort uint16) Option {
	return func(i *Interface) {
		port := int(listenPort)
		i.listenPort = &port
	}
}

func WithMTU(mtu uint16) Option {
	return func(i *Interface) {
		m := int(mtu)
		i.mtu = &m
	}
}

func WithFWMark(fwmark int) Option {
	return func(i *Interface) {
		i.fwmark = &fwmark
	}
}

func WithPeers(peers []wgtypes.PeerConfig) Option {
	return func(i *Interface) {
		i.peers = peers[:]
	}
}

func New(ctx context.Context, name string, options ...Option) (*Interface, error) {
	iface := &Interface{name: name}
	for i := range options {
		options[i](iface)
	}
	for i := range iface.addr {
		if iface.addr[i].IP.To16() == nil {
			return nil, net.InvalidAddrError(fmt.Sprintf("invalid address: %v", iface.addr[i]))
		}
	}
	return iface, nil
}

func (i *Interface) Start(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			err = errors.Join(err, i.Stop(ctx))
		}
	}()

	i.wgClient, err = wgctrl.New()
	if err != nil {
		return fmt.Errorf("failed to create the wireguard client: %w", err)
	}

	err = i.createInter(ctx)
	if err != nil {
		return fmt.Errorf("failed to create the wireguard interface: %w", err)
	}

	err = i.addInterAddr(ctx)
	if err != nil {
		return fmt.Errorf("failed to add addresses to wireguard interface: %w", err)
	}

	err = i.setMTU(ctx)
	if err != nil {
		return fmt.Errorf("failed to set mtu for the wireguard interface: %w", err)
	}

	err = i.startInter(ctx)
	if err != nil {
		return fmt.Errorf("failed to start the wireguard interface: %w", err)
	}

	err = i.wgClient.ConfigureDevice(i.name, wgtypes.Config{
		PrivateKey:   i.privateKey,
		ListenPort:   i.listenPort,
		ReplacePeers: true,
	})

	return nil
}

func (i *Interface) Stop(ctx context.Context) error {
	if err := i.stopInter(ctx); err != nil {
		return fmt.Errorf("failed to stop the wireguard interface: %w", err)
	}
	if err := i.wgClient.Close(); err != nil {
		return fmt.Errorf("failed to create the wireguard client: %w", err)
	}
	return nil
}

func (i *Interface) Restart(ctx context.Context) error {
	if err := i.stopInter(ctx); err != nil {
		return fmt.Errorf("failed to stop the wireguard interface: %w", err)
	}
	if err := i.startInter(ctx); err != nil {
		return fmt.Errorf("failed to start the wireguard interface: %w", err)
	}
	return nil
}

// PublicKey returns the public key of a Wireguard interface
func (i *Interface) PublicKey() wgtypes.Key {
	return i.privateKey.PublicKey()
}

// UpdatePeers updates the peer configuration for a Wireguard interface.
func (i *Interface) UpdatePeers(peers ...wgtypes.PeerConfig) error {
	return i.wgClient.ConfigureDevice(i.name, wgtypes.Config{
		Peers:        peers,
		ReplacePeers: false,
	})
}

// SetPeers replaces the peer configuration for a Wireguard interface.
func (i *Interface) SetPeers(peers ...wgtypes.PeerConfig) error {
	return i.wgClient.ConfigureDevice(i.name, wgtypes.Config{
		Peers:        peers,
		ReplacePeers: true,
	})
}
