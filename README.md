# wgiface

wgiface is a go package that provides easy wireguard interface initialization by wrapping `golang.zx2c4.com/wireguard/wgctrl` and abstracting os calls (such as `ip link add ...` on linux) for interface creation, configuration and management. It also provides methods to manage the lifecycle and peers of the interface, while using the original types from `golang.zx2c4.com/wireguard/wgctrl/wgtypes`.

This package currently supports the following platforms:
- Linux

## Examples
```go
wg, err := wgiface.New(ctx, "wg",
    wgiface.WithAddresses([]net.IPNet{
        {
            IP:   net.IPv4(10, 0, 0, 0),
            Mask: net.CIDRMask(24, 32),
        },
        {
            IP:   net.ParseIP("fc00::"),
            Mask: net.CIDRMask(120, 128),
        },
    }),
    wgiface.WithListenPort(51821),
)
// ...
err = i.UpdatePeers(wgtypes.PeerConfig{
    PublicKey: peerPrivateKey.PublicKey(),
    AllowedIPs: []net.IPNet{{
        IP:   net.IPv4(10, 0, 0, 1),
        Mask: net.CIDRMask(32, 32),
    }},
})
```