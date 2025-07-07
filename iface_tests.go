package wgiface

import (
	"context"
	"net"
	"testing"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// The following tests runs only inside the container test

func Test_Interface(t *testing.T) {
	ctx := context.Background()

	i, err := New(ctx, "wgtest",
		WithAddresses([]net.IPNet{
			{
				IP:   net.IPv4(10, 23, 45, 67),
				Mask: net.CIDRMask(24, 32),
			},
			{
				IP:   net.ParseIP("fc00::"),
				Mask: net.CIDRMask(7, 128),
			},
		}),
		WithListenPort(51821),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = i.Start(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = i.Stop(context.Background())
		if err != nil {
			t.Fatal(err)
		}
	}()

	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	err = i.UpdatePeers(wgtypes.PeerConfig{
		PublicKey: privateKey.PublicKey(),
		AllowedIPs: []net.IPNet{{
			IP:   net.IPv4(10, 0, 0, 1),
			Mask: net.CIDRMask(32, 32),
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
}
