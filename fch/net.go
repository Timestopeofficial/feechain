package fch

import (
	nodeconfig "github.com/Timestopeofficial/feechain/internal/configs/node"
	commonRPC "github.com/Timestopeofficial/feechain/rpc/common"
	"github.com/Timestopeofficial/feechain/staking/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

// GetCurrentUtilityMetrics ..
func (fch *Feechain) GetCurrentUtilityMetrics() (*network.UtilityMetric, error) {
	return network.NewUtilityMetricSnapshot(fch.BlockChain)
}

// GetPeerInfo returns the peer info to the node, including blocked peer, connected peer, number of peers
func (fch *Feechain) GetPeerInfo() commonRPC.NodePeerInfo {

	topics := fch.NodeAPI.ListTopic()
	p := make([]commonRPC.P, len(topics))

	for i, t := range topics {
		topicPeer := fch.NodeAPI.ListPeer(t)
		p[i].Topic = t
		p[i].Peers = make([]peer.ID, len(topicPeer))
		copy(p[i].Peers, topicPeer)
	}

	return commonRPC.NodePeerInfo{
		PeerID:       nodeconfig.GetPeerID(),
		BlockedPeers: fch.NodeAPI.ListBlockedPeer(),
		P:            p,
	}
}
