package cdn

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"serverless/cdn/storage"
	"serverless/cdn/utils"
	"time"
)

const DiscoveryInterval = time.Minute * 30

const DiscoveryServiceTag = "pubsub-cdn"

const UploadTopic = "upload"

type CDNHandler struct {
	ctx   context.Context
	ps    *pubsub.PubSub
	topic *pubsub.Topic
	sub   *pubsub.Subscription
	self  peer.ID
}

func (handler *CDNHandler) Upload(id string, file []byte) error {
	checksum := utils.ComputeChecksum(file)

	message := CDNMessage{
		FileId:   id,
		Checksum: checksum,
		Content:  string(file),
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return handler.topic.Publish(handler.ctx, data)
}

type CDNMessage struct {
	FileId   string
	Checksum string
	Content  string
}

func InitCDNHandler(ctx context.Context, storage *storage.StorageCDN, log bool) (*CDNHandler, error) {
	h, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		return nil, err
	}

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		return nil, err
	}

	if err := setupDiscovery(h); err != nil {
		return nil, err
	}

	topic, err := ps.Join(UploadTopic)
	if err != nil {
		return nil, err
	}

	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	if log {
		fmt.Println("CDN Handler Subscribed")
	}

	cdnHandler := &CDNHandler{
		ctx:   ctx,
		ps:    ps,
		topic: topic,
		sub:   sub,
		self:  h.ID(),
	}

	go cdnHandler.listen(storage)

	return cdnHandler, nil
}

func (handler *CDNHandler) listen(storage *storage.StorageCDN) {
	for {
		msg, err := handler.sub.Next(handler.ctx)
		if err != nil {
			panic(err)
		}
		if msg.ReceivedFrom == handler.self {
			continue
		}

		message := new(CDNMessage)
		err = json.Unmarshal(msg.Data, message)
		if err != nil {
			continue
		}

		err = storage.StoreFile(message.FileId, []byte(message.Content))
		if err != nil {
			panic(err)
		}
	}
}

type discoveryNotifee struct {
	h host.Host
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("discovered new peer %s\n", pi.ID)
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pi.ID, err)
	}
}

func setupDiscovery(h host.Host) error {
	s := mdns.NewMdnsService(h, DiscoveryServiceTag, &discoveryNotifee{h: h})
	return s.Start()
}
