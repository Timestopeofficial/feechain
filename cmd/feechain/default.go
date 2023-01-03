package main

import (
	feechainconfig "github.com/Timestopeofficial/feechain/internal/configs/feechain"
	nodeconfig "github.com/Timestopeofficial/feechain/internal/configs/node"
)

const tomlConfigVersion = "2.5.1" // bump from 2.5.0 for AccountSlots

const (
	defNetworkType = nodeconfig.Mainnet
)

var defaultConfig = feechainconfig.FeechainConfig{
	Version: tomlConfigVersion,
	General: feechainconfig.GeneralConfig{
		NodeType:         "validator",
		NoStaking:        false,
		ShardID:          -1,
		IsArchival:       false,
		IsBeaconArchival: false,
		IsOffline:        false,
		DataDir:          "./",
	},
	Network: getDefaultNetworkConfig(defNetworkType),
	P2P: feechainconfig.P2pConfig{
		Port:            nodeconfig.DefaultP2PPort,
		IP:              nodeconfig.DefaultPublicListenIP,
		KeyFile:         "./.nkey",
		DiscConcurrency: nodeconfig.DefaultP2PConcurrency,
		MaxConnsPerIP:   nodeconfig.DefaultMaxConnPerIP,
	},
	HTTP: feechainconfig.HttpConfig{
		Enabled:        true,
		RosettaEnabled: false,
		IP:             "127.0.0.1",
		Port:           nodeconfig.DefaultRPCPort,
		AuthPort:       nodeconfig.DefaultAuthRPCPort,
		RosettaPort:    nodeconfig.DefaultRosettaPort,
	},
	WS: feechainconfig.WsConfig{
		Enabled:  true,
		IP:       "127.0.0.1",
		Port:     nodeconfig.DefaultWSPort,
		AuthPort: nodeconfig.DefaultAuthWSPort,
	},
	RPCOpt: feechainconfig.RpcOptConfig{
		DebugEnabled:      false,
		RateLimterEnabled: true,
		RequestsPerSecond: nodeconfig.DefaultRPCRateLimit,
	},
	BLSKeys: feechainconfig.BlsConfig{
		KeyDir:   "./.blskeys",
		KeyFiles: []string{},
		MaxKeys:  10,

		PassEnabled:      true,
		PassSrcType:      blsPassTypeAuto,
		PassFile:         "",
		SavePassphrase:   false,
		KMSEnabled:       false,
		KMSConfigSrcType: kmsConfigTypeShared,
		KMSConfigFile:    "",
	},
	TxPool: feechainconfig.TxPoolConfig{
		BlacklistFile:  "./blacklist.txt",
		RosettaFixFile: "",
		AccountSlots:   16,
	},
	Sync: getDefaultSyncConfig(defNetworkType),
	Pprof: feechainconfig.PprofConfig{
		Enabled:            false,
		ListenAddr:         "127.0.0.1:6060",
		Folder:             "./profiles",
		ProfileNames:       []string{},
		ProfileIntervals:   []int{600},
		ProfileDebugValues: []int{0},
	},
	Log: feechainconfig.LogConfig{
		Folder:       "./logs",
		FileName:     "feechain.log",
		RotateSize:   100,
		RotateCount:  0,
		RotateMaxAge: 0,
		Verbosity:    3,
		VerbosePrints: feechainconfig.LogVerbosePrints{
			Config: true,
		},
	},
	DNSSync: getDefaultDNSSyncConfig(defNetworkType),
	ShardData: feechainconfig.ShardDataConfig{
		EnableShardData: false,
		DiskCount:       8,
		ShardCount:      4,
		CacheTime:       10,
		CacheSize:       512,
	},
}

var defaultSysConfig = feechainconfig.SysConfig{
	NtpServer: "1.pool.ntp.org",
}

var defaultDevnetConfig = feechainconfig.DevnetConfig{
	NumShards:   2,
	ShardSize:   10,
	FchNodeSize: 10,
}

var defaultRevertConfig = feechainconfig.RevertConfig{
	RevertBeacon: false,
	RevertBefore: 0,
	RevertTo:     0,
}

var defaultLogContext = feechainconfig.LogContext{
	IP:   "127.0.0.1",
	Port: 17000,
}

var defaultConsensusConfig = feechainconfig.ConsensusConfig{
	MinPeers:     6,
	AggregateSig: true,
}

var defaultPrometheusConfig = feechainconfig.PrometheusConfig{
	Enabled:    true,
	IP:         "0.0.0.0",
	Port:       17900,
	EnablePush: false,
	Gateway:    "https://gateway.timestope.net",
}

var (
	defaultMainnetSyncConfig = feechainconfig.SyncConfig{
		Enabled:        false,
		Downloader:     false,
		Concurrency:    4,
		MinPeers:       4,
		InitStreams:    4,
		DiscSoftLowCap: 8,
		DiscHardLowCap: 6,
		DiscHighCap:    128,
		DiscBatch:      8,
	}

	defaultTestNetSyncConfig = feechainconfig.SyncConfig{
		Enabled:        true,
		Downloader:     false,
		Concurrency:    2,
		MinPeers:       2,
		InitStreams:    2,
		DiscSoftLowCap: 2,
		DiscHardLowCap: 2,
		DiscHighCap:    1024,
		DiscBatch:      3,
	}

	defaultLocalNetSyncConfig = feechainconfig.SyncConfig{
		Enabled:        true,
		Downloader:     true,
		Concurrency:    2,
		MinPeers:       2,
		InitStreams:    2,
		DiscSoftLowCap: 2,
		DiscHardLowCap: 2,
		DiscHighCap:    1024,
		DiscBatch:      3,
	}

	defaultElseSyncConfig = feechainconfig.SyncConfig{
		Enabled:        true,
		Downloader:     true,
		Concurrency:    4,
		MinPeers:       4,
		InitStreams:    4,
		DiscSoftLowCap: 4,
		DiscHardLowCap: 4,
		DiscHighCap:    1024,
		DiscBatch:      8,
	}
)

const (
	defaultBroadcastInvalidTx = false
)

func getDefaultFchConfigCopy(nt nodeconfig.NetworkType) feechainconfig.FeechainConfig {
	config := defaultConfig

	config.Network = getDefaultNetworkConfig(nt)
	if nt == nodeconfig.Devnet {
		devnet := getDefaultDevnetConfigCopy()
		config.Devnet = &devnet
	}
	config.Sync = getDefaultSyncConfig(nt)
	config.DNSSync = getDefaultDNSSyncConfig(nt)

	return config
}

func getDefaultSysConfigCopy() feechainconfig.SysConfig {
	config := defaultSysConfig
	return config
}

func getDefaultDevnetConfigCopy() feechainconfig.DevnetConfig {
	config := defaultDevnetConfig
	return config
}

func getDefaultRevertConfigCopy() feechainconfig.RevertConfig {
	config := defaultRevertConfig
	return config
}

func getDefaultLogContextCopy() feechainconfig.LogContext {
	config := defaultLogContext
	return config
}

func getDefaultConsensusConfigCopy() feechainconfig.ConsensusConfig {
	config := defaultConsensusConfig
	return config
}

func getDefaultPrometheusConfigCopy() feechainconfig.PrometheusConfig {
	config := defaultPrometheusConfig
	return config
}

const (
	nodeTypeValidator = "validator"
	nodeTypeExplorer  = "explorer"
)

const (
	blsPassTypeAuto   = "auto"
	blsPassTypeFile   = "file"
	blsPassTypePrompt = "prompt"

	kmsConfigTypeShared = "shared"
	kmsConfigTypePrompt = "prompt"
	kmsConfigTypeFile   = "file"

	legacyBLSPassTypeDefault = "default"
	legacyBLSPassTypeStdin   = "stdin"
	legacyBLSPassTypeDynamic = "no-prompt"
	legacyBLSPassTypePrompt  = "prompt"
	legacyBLSPassTypeStatic  = "file"
	legacyBLSPassTypeNone    = "none"

	legacyBLSKmsTypeDefault = "default"
	legacyBLSKmsTypePrompt  = "prompt"
	legacyBLSKmsTypeFile    = "file"
	legacyBLSKmsTypeNone    = "none"
)
