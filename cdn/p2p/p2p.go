package p2p

import (
	"context"
	"fmt"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	peerstore "github.com/libp2p/go-libp2p/core/peerstore"
)

// Node represents a libp2p node in the network.
type Node struct {
	Host      host.Host
	PeerStore peerstore.Peerstore
	Context   context.Context
	dht       *dht.IpfsDHT
	// Add more fields as needed.
}

// NewNode initializes a new libp2p node.
func NewNode(ctx context.Context) (*Node, error) {
	// Initialize the libp2p host and other configurations.
	h, err := NewHost(ctx)
	if err != nil {
		return nil, err
	}

	// Initialize the DHT for the node.
	dht, err := NewDHT(ctx, h, nil)
	if err != nil {
		return nil, err
	}

	// Create a new Node instance.
	node := &Node{
		Host:      h,
		PeerStore: h.Peerstore(),
		Context:   ctx,
		dht:       dht,
	}

	return node, nil
}

func NewHost(ctx context.Context) (host.Host, error) {
	// Generate a key pair and create a new host
	privateKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		return nil, err
	}
	h, err := libp2p.New(libp2p.Identity(privateKey))
	if err != nil {
		return nil, err
	}
	return h, nil
}

func NewDHT(ctx context.Context, host host.Host, bootstrapPeers []multiaddr.Multiaddr) (*dht.IpfsDHT, error) {
	options := []dht.Option{}
	if len(bootstrapPeers) == 0 {
		options = append(options, dht.Mode(dht.ModeServer))
	}

	kdht, err := dht.New(ctx, host, options...)
	if err != nil {
		return nil, err
	}

	if err = kdht.Bootstrap(ctx); err != nil {
		return nil, err
	}

	// Connect to bootstrap peers
	for _, addr := range bootstrapPeers {
		peerInfo, _ := peer.AddrInfoFromP2pAddr(addr)
		if err := host.Connect(ctx, *peerInfo); err != nil {
			fmt.Println("Failed to connect to bootstrap peer:", err)
		}
	}

	return kdht, nil
}
