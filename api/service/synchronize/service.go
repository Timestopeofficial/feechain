package synchronize

import (
	"github.com/Timestopeofficial/feechain/core"
	"github.com/Timestopeofficial/feechain/fch/downloader"
	"github.com/Timestopeofficial/feechain/p2p"
)

// Service is simply a adapter of Downloaders, which support block synchronization
type Service struct {
	Downloaders *downloader.Downloaders
}

// NewService creates the a new downloader service
func NewService(host p2p.Host, bcs []*core.BlockChain, config downloader.Config) *Service {
	return &Service{
		Downloaders: downloader.NewDownloaders(host, bcs, config),
	}
}

// Start start the service
func (s *Service) Start() error {
	s.Downloaders.Start()
	return nil
}

// Stop stop the service
func (s *Service) Stop() error {
	s.Downloaders.Close()
	return nil
}
