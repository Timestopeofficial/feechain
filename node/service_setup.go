package node

import (
	"fmt"

	"github.com/Timestopeofficial/feechain/api/service"
	"github.com/Timestopeofficial/feechain/api/service/blockproposal"
	"github.com/Timestopeofficial/feechain/api/service/consensus"
	"github.com/Timestopeofficial/feechain/api/service/explorer"
)

// RegisterValidatorServices register the validator services.
func (node *Node) RegisterValidatorServices() {
	// Register consensus service.
	node.serviceManager.Register(
		service.Consensus,
		consensus.New(node.BlockChannel, node.Consensus, node.startConsensus),
	)
	// Register new block service.
	node.serviceManager.Register(
		service.BlockProposal,
		blockproposal.New(node.Consensus.ReadySignal, node.Consensus.CommitSigChannel, node.WaitForConsensusReadyV2),
	)
}

// RegisterExplorerServices register the explorer services
func (node *Node) RegisterExplorerServices() {
	// Register explorer service.
	node.serviceManager.Register(
		service.SupportExplorer, explorer.New(&node.SelfPeer, node.Blockchain(), node),
	)
}

// RegisterService register a service to the node service manager
func (node *Node) RegisterService(st service.Type, s service.Service) {
	node.serviceManager.Register(st, s)
}

// StartServices runs registered services.
func (node *Node) StartServices() error {
	return node.serviceManager.StartServices()
}

// StopServices runs registered services.
func (node *Node) StopServices() error {
	return node.serviceManager.StopServices()
}

func (node *Node) networkInfoDHTPath() string {
	return fmt.Sprintf(".dht-%s-%s-c%s",
		node.SelfPeer.IP,
		node.SelfPeer.Port,
		node.chainConfig.ChainID,
	)
}
