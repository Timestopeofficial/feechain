package version

// Service check the latest version
type Service struct {
	CheckVersion func()
}

// New creates the a new version checking service
func New(CheckVersion func()) *Service {
	return &Service{CheckVersion: CheckVersion}
}

// Start start the service
func (s *Service) Start() error {
	s.CheckVersion()
	return nil
}

// Stop stop the service
func (s *Service) Stop() error {
	return nil
}
