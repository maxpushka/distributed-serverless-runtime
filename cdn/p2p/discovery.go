package p2p

import (
	"context"
	"fmt"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"time"
)

// DiscoverPeers handles peer discovery in the network.
func (n *Node) DiscoverPeers() error {
	// Implement peer discovery logic.
	return nil
}

func Discover(ctx context.Context, h host.Host, dht *dht.IpfsDHT, rendezvous string) {
	rt := routing.NewRoutingDiscovery(dht)

	// Advertise presence in the DHT
	go func() {
		for {
			ttl, err := rt.Advertise(ctx, rendezvous)
			if err != nil {
				fmt.Println("Advertise failed:", err)
				continue
			}
			fmt.Println("Advertise success with TTL:", ttl)
			time.Sleep(3 * time.Second) // Adjust as needed
		}
	}()

	// Periodically find peers
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for range ticker.C {
		peers, err := rt.FindPeers(ctx, rendezvous)
		if err != nil {
			fmt.Println("Failed to find peers:", err)
			continue
		}

		fmt.Printf("FindPeers succeeded: found peers: %d\n", len(peers))
		for p := range peers {
			// Connect to discovered peers (if not already connected)
			if h.Network().Connectedness(p.ID) != network.Connected {
				_, err = h.Network().DialPeer(ctx, p.ID)
				if err != nil {
					fmt.Println("Failed to connect to peer:", err)
				} else {
					fmt.Println("Connected to peer:", p.ID)
				}
			}
		}
	}
}
