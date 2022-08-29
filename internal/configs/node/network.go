package nodeconfig

var (
	mainnetBootNodes = []string{
		"/ip4/107.155.65.131/tcp/28175/p2p/QmYAppojcNmbVu6Wb2CvQRwKns1TuKAP1S7KxpBxf8NS3F",
		"/ip4/89.233.108.44/tcp/28175/p2p/QmZKsSuz714uDhLcNLCurPLG9E1R8UBRJ9KaGXCBzK9jhe",
	}

	testnetBootNodes = []string{
		"/ip4/107.155.65.131/tcp/28175/p2p/QmYAppojcNmbVu6Wb2CvQRwKns1TuKAP1S7KxpBxf8NS3F",
		"/ip4/89.233.108.44/tcp/28175/p2p/QmZKsSuz714uDhLcNLCurPLG9E1R8UBRJ9KaGXCBzK9jhe",
	}

	pangaeaBootNodes = []string{
		"/ip4/52.40.84.2/tcp/9800/p2p/QmbPVwrqWsTYXq1RxGWcxx9SWaTUCfoo1wA6wmdbduWe29",
		"/ip4/54.86.126.90/tcp/9800/p2p/Qmdfjtk6hPoyrH1zVD9PEH4zfWLo38dP2mDvvKXfh3tnEv",
	}

	partnerBootNodes = []string{
		"/ip4/107.155.65.131/tcp/28175/p2p/QmYAppojcNmbVu6Wb2CvQRwKns1TuKAP1S7KxpBxf8NS3F",
		"/ip4/89.233.108.44/tcp/28175/p2p/QmZKsSuz714uDhLcNLCurPLG9E1R8UBRJ9KaGXCBzK9jhe",
	}

	stressBootNodes = []string{
		"/ip4/107.155.65.131/tcp/28175/p2p/QmYAppojcNmbVu6Wb2CvQRwKns1TuKAP1S7KxpBxf8NS3F",
		"/ip4/89.233.108.44/tcp/28175/p2p/QmZKsSuz714uDhLcNLCurPLG9E1R8UBRJ9KaGXCBzK9jhe",
	}

	devnetBootNodes = []string{}
)

const (
	mainnetDNSZone   = "b.timestope.net"
	testnetDNSZone   = "t.timestope.net"
	pangaeaDNSZone   = "p.timestope.net"
	partnerDNSZone   = "pa.timestope.net"
	stressnetDNSZone = "st.timestope.net"
)

const (
	// DefaultLocalListenIP is the IP used for local hosting
	DefaultLocalListenIP = "127.0.0.1"
	// DefaultPublicListenIP is the IP used for public hosting
	DefaultPublicListenIP = "0.0.0.0"
	// DefaultP2PPort is the key to be used for p2p communication
	DefaultP2PPort = 17000
	// DefaultLegacyDNSPort is the default legacy DNS port. The actual port used is DNSPort - 3000. This is a
	// very bad design. Refactored to DefaultDNSPort
	DefaultLegacyDNSPort = 17000
	// DefaultDNSPort is the default DNS port for both remote node and local server.
	DefaultDNSPort = 14000
	// DefaultRPCPort is the default rpc port. The actual port used is 17000+500
	DefaultRPCPort = 17500
	// DefaultAuthRPCPort is the default rpc auth port. The actual port used is 17000+501
	DefaultAuthRPCPort = 17501
	// DefaultRosettaPort is the default rosetta port. The actual port used is 17000+700
	DefaultRosettaPort = 17700
	// DefaultWSPort is the default port for web socket endpoint. The actual port used is
	DefaultWSPort = 17800
	// DefaultAuthWSPort is the default port for web socket auth endpoint. The actual port used is
	DefaultAuthWSPort = 17801
	// DefaultPrometheusPort is the default prometheus port. The actual port used is 17000+900
	DefaultPrometheusPort = 17900
	// DefaultP2PConcurrency is the default P2P concurrency, 0 means is set the default value of P2P Discovery, the actual value is 10
	DefaultP2PConcurrency = 0
	DefaultMaxConnPerIP   = 10
)

const (
	// DefaultRateLimit for RPC, the number of requests per second
	DefaultRPCRateLimit = 1000
)

const (
	// rpcHTTPPortOffset is the port offset for RPC HTTP requests
	rpcHTTPPortOffset = 500

	// rpcHTTPAuthPortOffset is the port offset for RPC Auth HTTP requests
	rpcHTTPAuthPortOffset = 501

	// rpcHTTPPortOffset is the port offset for rosetta HTTP requests
	rosettaHTTPPortOffset = 700

	// rpcWSPortOffSet is the port offset for RPC websocket requests
	rpcWSPortOffSet = 800

	// rpcWSAuthPortOffSet is the port offset for RPC Auth websocket requests
	rpcWSAuthPortOffSet = 801

	// prometheusHTTPPortOffset is the port offset for prometheus HTTP requests
	prometheusHTTPPortOffset = 900
)

// GetDefaultBootNodes get the default bootnode with the given network type
func GetDefaultBootNodes(networkType NetworkType) []string {
	switch networkType {
	case Mainnet:
		return mainnetBootNodes
	case Testnet:
		return testnetBootNodes
	case Pangaea:
		return pangaeaBootNodes
	case Partner:
		return partnerBootNodes
	case Stressnet:
		return stressBootNodes
	case Devnet:
		return devnetBootNodes
	}
	return nil
}

// GetDefaultDNSZone get the default DNS zone with the given network type
func GetDefaultDNSZone(networkType NetworkType) string {
	switch networkType {
	case Mainnet:
		return mainnetDNSZone
	case Testnet:
		return testnetDNSZone
	case Pangaea:
		return pangaeaDNSZone
	case Partner:
		return partnerDNSZone
	case Stressnet:
		return stressnetDNSZone
	}
	return ""
}

// GetDefaultDNSPort get the default DNS port for the given network type
func GetDefaultDNSPort(NetworkType) int {
	return DefaultDNSPort
}

// GetRPCHTTPPortFromBase return the rpc HTTP port from base port
func GetRPCHTTPPortFromBase(basePort int) int {
	return basePort + rpcHTTPPortOffset
}

// GetRPCAuthHTTPPortFromBase return the rpc HTTP port from base port
func GetRPCAuthHTTPPortFromBase(basePort int) int {
	return basePort + rpcHTTPAuthPortOffset
}

// GetRosettaHTTPPortFromBase return the rosetta HTTP port from base port
func GetRosettaHTTPPortFromBase(basePort int) int {
	return basePort + rosettaHTTPPortOffset
}

// GetWSPortFromBase return the Websocket port from the base port
func GetWSPortFromBase(basePort int) int {
	return basePort + rpcWSPortOffSet
}

// GetWSAuthPortFromBase return the Websocket port from the base auth port
func GetWSAuthPortFromBase(basePort int) int {
	return basePort + rpcWSAuthPortOffSet
}

// GetPrometheusHTTPPortFromBase return the prometheus HTTP port from base port
func GetPrometheusHTTPPortFromBase(basePort int) int {
	return basePort + prometheusHTTPPortOffset
}
