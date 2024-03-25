package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	feechainconfig "github.com/Timestopeofficial/feechain/internal/configs/feechain"

	"github.com/spf13/cobra"

	"github.com/Timestopeofficial/feechain/internal/cli"
	nodeconfig "github.com/Timestopeofficial/feechain/internal/configs/node"
)

var (
	trueBool = true
)

func TestFeechainFlags(t *testing.T) {
	tests := []struct {
		argStr    string
		expConfig feechainconfig.FeechainConfig
	}{
		{
			// running staking command from legacy node.sh
			argStr: "--bootnodes /ip4/100.26.90.187/tcp/9874/p2p/Qmdfjtk6hPoyrH1zVD9PEH4zfWLo38dP2mDvvKXfh3tnEv," +
				"/ip4/54.213.43.194/tcp/9874/p2p/QmZJJx6AdaoEkGLrYG4JeLCKeCKDjnFz2wfHNHxAqFSGA9,/ip4/13.113.101." +
				"219/tcp/12019/p2p/QmQayinFSgMMw5cSpDUiD9pQ2WeP6WNmGxpZ6ou3mdVFJX,/ip4/99.81.170.167/tcp/12019/p" +
				"2p/QmRVbTpEYup8dSaURZfF6ByrMTSKa4UyUzJhSjahFzRqNj --ip 8.8.8.8 --port 9000 --network_type=mainn" +
				"et --dns_zone=asadal.timestope.net --blacklist=./.fch/blacklist.txt --min_peers=6 --max_bls_keys_per_node=" +
				"10 --broadcast_invalid_tx=true --verbosity=3 --is_archival=false --shard_id=-1 --staking=true -" +
				"-aws-config-source file:config.json --p2p.disc.concurrency 5 --p2p.security.max-conn-per-ip 5",
			expConfig: feechainconfig.FeechainConfig{
				Version: tomlConfigVersion,
				General: feechainconfig.GeneralConfig{
					NodeType:   "validator",
					NoStaking:  false,
					ShardID:    -1,
					IsArchival: false,
					DataDir:    "./",
				},
				Network: feechainconfig.NetworkConfig{
					NetworkType: "asadal",
					BootNodes: []string{
						"/ip4/100.26.90.187/tcp/9874/p2p/Qmdfjtk6hPoyrH1zVD9PEH4zfWLo38dP2mDvvKXfh3tnEv",
						"/ip4/54.213.43.194/tcp/9874/p2p/QmZJJx6AdaoEkGLrYG4JeLCKeCKDjnFz2wfHNHxAqFSGA9",
						"/ip4/13.113.101.219/tcp/12019/p2p/QmQayinFSgMMw5cSpDUiD9pQ2WeP6WNmGxpZ6ou3mdVFJX",
						"/ip4/99.81.170.167/tcp/12019/p2p/QmRVbTpEYup8dSaURZfF6ByrMTSKa4UyUzJhSjahFzRqNj",
					},
				},
				DNSSync: feechainconfig.DnsSync{
					Port:       6000,
					Zone:       "asadal.timestope.net",
					Server:     true,
					Client:     true,
					ServerPort: nodeconfig.DefaultDNSPort,
				},
				P2P: feechainconfig.P2pConfig{
					Port:            9000,
					IP:              defaultConfig.P2P.IP,
					KeyFile:         defaultConfig.P2P.KeyFile,
					DiscConcurrency: 5,
					MaxConnsPerIP:   5,
				},
				HTTP: feechainconfig.HttpConfig{
					Enabled:        true,
					IP:             "127.0.0.1",
					Port:           9500,
					AuthPort:       9501,
					RosettaEnabled: false,
					RosettaPort:    9700,
				},
				RPCOpt: feechainconfig.RpcOptConfig{
					DebugEnabled:      false,
					RateLimterEnabled: true,
					RequestsPerSecond: 1000,
				},
				WS: feechainconfig.WsConfig{
					Enabled:  true,
					IP:       "127.0.0.1",
					Port:     9800,
					AuthPort: 9801,
				},
				Consensus: &feechainconfig.ConsensusConfig{
					MinPeers:     6,
					AggregateSig: true,
				},
				BLSKeys: feechainconfig.BlsConfig{
					KeyDir:           "./.blskeys",
					KeyFiles:         []string{},
					MaxKeys:          10,
					PassEnabled:      true,
					PassSrcType:      "auto",
					PassFile:         "",
					SavePassphrase:   false,
					KMSEnabled:       false,
					KMSConfigSrcType: "file",
					KMSConfigFile:    "config.json",
				},
				TxPool: feechainconfig.TxPoolConfig{
					BlacklistFile:  "./.fch/blacklist.txt",
					RosettaFixFile: "",
					AccountSlots:   16,
				},
				Pprof: feechainconfig.PprofConfig{
					Enabled:            false,
					ListenAddr:         "127.0.0.1:6060",
					Folder:             "./profiles",
					ProfileNames:       []string{},
					ProfileIntervals:   []int{600},
					ProfileDebugValues: []int{0},
				},
				Log: feechainconfig.LogConfig{
					Folder:       "./latest",
					FileName:     "validator-8.8.8.8-9000.log",
					RotateSize:   100,
					RotateCount:  0,
					RotateMaxAge: 0,
					Verbosity:    3,
					Context: &feechainconfig.LogContext{
						IP:   "8.8.8.8",
						Port: 9000,
					},
					VerbosePrints: feechainconfig.LogVerbosePrints{
						Config: true,
					},
				},
				Sys: &feechainconfig.SysConfig{
					NtpServer: defaultSysConfig.NtpServer,
				},
				Legacy: &feechainconfig.LegacyConfig{
					TPBroadcastInvalidTxn: &trueBool,
				},
				Prometheus: &feechainconfig.PrometheusConfig{
					Enabled:    true,
					IP:         "0.0.0.0",
					Port:       9900,
					EnablePush: true,
					Gateway:    "https://gateway.mojaik.com",
				},
				Sync: defaultMainnetSyncConfig,
				ShardData: feechainconfig.ShardDataConfig{
					EnableShardData: false,
					DiskCount:       8,
					ShardCount:      4,
					CacheTime:       10,
					CacheSize:       512,
				},
			},
		},
	}
	for i, test := range tests {
		ts := newFlagTestSuite(t, getRootFlags(), applyRootFlags)
		hc, err := ts.run(strings.Split(test.argStr, " "))
		if err != nil {
			t.Fatalf("Test %v: %v", i, err)
		}
		if !reflect.DeepEqual(hc, test.expConfig) {
			t.Errorf("Test %v: unexpected config: \n\t%+v\n\t%+v", i, hc, test.expConfig)
		}
		ts.tearDown()
	}
}

func TestGeneralFlags(t *testing.T) {
	tests := []struct {
		args      []string
		expConfig feechainconfig.GeneralConfig
		expErr    error
	}{
		{
			args: []string{},
			expConfig: feechainconfig.GeneralConfig{
				NodeType:   "validator",
				NoStaking:  false,
				ShardID:    -1,
				IsArchival: false,
				DataDir:    "./",
			},
		},
		{
			args: []string{"--run", "explorer", "--run.legacy", "--run.shard=0",
				"--run.archive=true", "--datadir=./.fch"},
			expConfig: feechainconfig.GeneralConfig{
				NodeType:   "explorer",
				NoStaking:  true,
				ShardID:    0,
				IsArchival: true,
				DataDir:    "./.fch",
			},
		},
		{
			args: []string{"--node_type", "explorer", "--staking", "--shard_id", "0",
				"--is_archival", "--db_dir", "./"},
			expConfig: feechainconfig.GeneralConfig{
				NodeType:   "explorer",
				NoStaking:  false,
				ShardID:    0,
				IsArchival: true,
				DataDir:    "./",
			},
		},
		{
			args: []string{"--staking=false", "--is_archival=false"},
			expConfig: feechainconfig.GeneralConfig{
				NodeType:   "validator",
				NoStaking:  true,
				ShardID:    -1,
				IsArchival: false,
				DataDir:    "./",
			},
		},
		{
			args: []string{"--run", "explorer", "--run.shard", "0"},
			expConfig: feechainconfig.GeneralConfig{
				NodeType:   "explorer",
				NoStaking:  false,
				ShardID:    0,
				IsArchival: false,
				DataDir:    "./",
			},
		},
		{
			args: []string{"--run", "explorer", "--run.shard", "0", "--run.archive=false"},
			expConfig: feechainconfig.GeneralConfig{
				NodeType:   "explorer",
				NoStaking:  false,
				ShardID:    0,
				IsArchival: false,
				DataDir:    "./",
			},
		},
	}
	for i, test := range tests {
		ts := newFlagTestSuite(t, generalFlags, applyGeneralFlags)

		got, err := ts.run(test.args)

		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Fatalf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}
		if !reflect.DeepEqual(got.General, test.expConfig) {
			t.Errorf("Test %v: unexpected config: \n\t%+v\n\t%+v", i, got.General, test.expConfig)
		}
		ts.tearDown()
	}
}

func TestNetworkFlags(t *testing.T) {
	tests := []struct {
		args      []string
		expConfig feechainconfig.FeechainConfig
		expErr    error
	}{
		{
			args: []string{},
			expConfig: feechainconfig.FeechainConfig{
				Network: feechainconfig.NetworkConfig{
					NetworkType: defNetworkType,
					BootNodes:   nodeconfig.GetDefaultBootNodes(defNetworkType),
				},
				DNSSync: getDefaultDNSSyncConfig(defNetworkType)},
		},
		{
			args: []string{"-n", "stn"},
			expConfig: feechainconfig.FeechainConfig{
				Network: feechainconfig.NetworkConfig{
					NetworkType: nodeconfig.Stressnet,
					BootNodes:   nodeconfig.GetDefaultBootNodes(nodeconfig.Stressnet),
				},
				DNSSync: getDefaultDNSSyncConfig(nodeconfig.Stressnet),
			},
		},
		{
			args: []string{"--network", "stk", "--bootnodes", "1,2,3,4", "--dns.zone", "8.8.8.8",
				"--dns.port", "9001", "--dns.server-port", "9002"},
			expConfig: feechainconfig.FeechainConfig{
				Network: feechainconfig.NetworkConfig{
					NetworkType: "pangaea",
					BootNodes:   []string{"1", "2", "3", "4"},
				},
				DNSSync: feechainconfig.DnsSync{
					Port:          9001,
					Zone:          "8.8.8.8",
					LegacySyncing: false,
					Server:        true,
					ServerPort:    9002,
				},
			},
		},
		{
			args: []string{"--network_type", "stk", "--bootnodes", "1,2,3,4", "--dns_zone", "8.8.8.8",
				"--dns_port", "9001"},
			expConfig: feechainconfig.FeechainConfig{
				Network: feechainconfig.NetworkConfig{
					NetworkType: "pangaea",
					BootNodes:   []string{"1", "2", "3", "4"},
				},
				DNSSync: feechainconfig.DnsSync{
					Port:          9001,
					Zone:          "8.8.8.8",
					LegacySyncing: false,
					Server:        true,
					ServerPort:    nodeconfig.GetDefaultDNSPort(nodeconfig.Pangaea),
				},
			},
		},
		{
			args: []string{"--dns=false"},
			expConfig: feechainconfig.FeechainConfig{
				Network: feechainconfig.NetworkConfig{
					NetworkType: defNetworkType,
					BootNodes:   nodeconfig.GetDefaultBootNodes(defNetworkType),
				},
				DNSSync: feechainconfig.DnsSync{
					Port:          nodeconfig.GetDefaultDNSPort(defNetworkType),
					Zone:          nodeconfig.GetDefaultDNSZone(defNetworkType),
					LegacySyncing: true,
					Client:        true,
					Server:        true,
					ServerPort:    nodeconfig.GetDefaultDNSPort(nodeconfig.Pangaea),
				},
			},
		},
	}
	for i, test := range tests {
		neededFlags := make([]cli.Flag, 0)
		neededFlags = append(neededFlags, networkFlags...)
		neededFlags = append(neededFlags, dnsSyncFlags...)
		ts := newFlagTestSuite(t, neededFlags, func(cmd *cobra.Command, config *feechainconfig.FeechainConfig) {
			// This is the network related logic in function getfeechainconfig.FeechainConfig
			nt := getNetworkType(cmd)
			config.Network = getDefaultNetworkConfig(nt)
			config.DNSSync = getDefaultDNSSyncConfig(nt)
			applyNetworkFlags(cmd, config)
			applyDNSSyncFlags(cmd, config)
		})

		got, err := ts.run(test.args)

		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Fatalf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}
		if !reflect.DeepEqual(got.Network, test.expConfig.Network) {
			t.Errorf("Test %v: unexpected network config: \n\t%+v\n\t%+v", i, got.Network, test.expConfig.Network)
		}
		if !reflect.DeepEqual(got.DNSSync, test.expConfig.DNSSync) {
			t.Errorf("Test %v: unexpected dnssync config: \n\t%+v\n\t%+v", i, got.DNSSync, test.expConfig.DNSSync)
		}
		ts.tearDown()
	}
}

var defDataStore = ".dht-127.0.0.1"

func TestP2PFlags(t *testing.T) {
	tests := []struct {
		args      []string
		expConfig feechainconfig.P2pConfig
		expErr    error
	}{
		{
			args:      []string{},
			expConfig: defaultConfig.P2P,
		},
		{
			args: []string{"--p2p.port", "9001", "--p2p.keyfile", "./key.file", "--p2p.dht.datastore",
				defDataStore},
			expConfig: feechainconfig.P2pConfig{
				Port:          9001,
				IP:            nodeconfig.DefaultPublicListenIP,
				KeyFile:       "./key.file",
				DHTDataStore:  &defDataStore,
				MaxConnsPerIP: 10,
			},
		},
		{
			args: []string{"--port", "9001", "--key", "./key.file"},
			expConfig: feechainconfig.P2pConfig{
				Port:          9001,
				IP:            nodeconfig.DefaultPublicListenIP,
				KeyFile:       "./key.file",
				MaxConnsPerIP: 10,
			},
		},
		{
			args: []string{"--p2p.port", "9001", "--p2p.disc.concurrency", "5", "--p2p.security.max-conn-per-ip", "5"},
			expConfig: feechainconfig.P2pConfig{
				Port:            9001,
				IP:              nodeconfig.DefaultPublicListenIP,
				KeyFile:         "./.nkey",
				DiscConcurrency: 5,
				MaxConnsPerIP:   5,
			},
		},
	}
	for i, test := range tests {
		ts := newFlagTestSuite(t, append(p2pFlags, legacyMiscFlags...),
			func(cmd *cobra.Command, config *feechainconfig.FeechainConfig) {
				applyLegacyMiscFlags(cmd, config)
				applyP2PFlags(cmd, config)
			},
		)

		got, err := ts.run(test.args)

		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Fatalf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}
		if !reflect.DeepEqual(got.P2P, test.expConfig) {
			t.Errorf("Test %v: unexpected config: \n\t%+v\n\t%+v", i, got.P2P, test.expConfig)
		}
		ts.tearDown()
	}
}

func TestRPCFlags(t *testing.T) {
	tests := []struct {
		args      []string
		expConfig feechainconfig.HttpConfig
		expErr    error
	}{
		{
			args:      []string{},
			expConfig: defaultConfig.HTTP,
		},
		{
			args: []string{"--http=false"},
			expConfig: feechainconfig.HttpConfig{
				Enabled:        false,
				RosettaEnabled: false,
				IP:             defaultConfig.HTTP.IP,
				Port:           defaultConfig.HTTP.Port,
				AuthPort:       defaultConfig.HTTP.AuthPort,
				RosettaPort:    defaultConfig.HTTP.RosettaPort,
			},
		},
		{
			args: []string{"--http.ip", "8.8.8.8", "--http.port", "9001"},
			expConfig: feechainconfig.HttpConfig{
				Enabled:        true,
				RosettaEnabled: false,
				IP:             "8.8.8.8",
				Port:           9001,
				AuthPort:       defaultConfig.HTTP.AuthPort,
				RosettaPort:    defaultConfig.HTTP.RosettaPort,
			},
		},
		{
			args: []string{"--http.ip", "8.8.8.8", "--http.auth-port", "9001"},
			expConfig: feechainconfig.HttpConfig{
				Enabled:        true,
				RosettaEnabled: false,
				IP:             "8.8.8.8",
				Port:           defaultConfig.HTTP.Port,
				AuthPort:       9001,
				RosettaPort:    defaultConfig.HTTP.RosettaPort,
			},
		},
		{
			args: []string{"--http.ip", "8.8.8.8", "--http.port", "9001", "--http.rosetta.port", "10001"},
			expConfig: feechainconfig.HttpConfig{
				Enabled:        true,
				RosettaEnabled: true,
				IP:             "8.8.8.8",
				Port:           9001,
				AuthPort:       defaultConfig.HTTP.AuthPort,
				RosettaPort:    10001,
			},
		},
		{
			args: []string{"--http.ip", "8.8.8.8", "--http.rosetta.port", "10001"},
			expConfig: feechainconfig.HttpConfig{
				Enabled:        true,
				RosettaEnabled: true,
				IP:             "8.8.8.8",
				Port:           defaultConfig.HTTP.Port,
				AuthPort:       defaultConfig.HTTP.AuthPort,
				RosettaPort:    10001,
			},
		},
		{
			args: []string{"--ip", "8.8.8.8", "--port", "9001", "--public_rpc"},
			expConfig: feechainconfig.HttpConfig{
				Enabled:        true,
				RosettaEnabled: false,
				IP:             nodeconfig.DefaultPublicListenIP,
				Port:           9501,
				AuthPort:       9502,
				RosettaPort:    9701,
			},
		},
	}
	for i, test := range tests {
		ts := newFlagTestSuite(t, append(httpFlags, legacyMiscFlags...),
			func(cmd *cobra.Command, config *feechainconfig.FeechainConfig) {
				applyLegacyMiscFlags(cmd, config)
				applyHTTPFlags(cmd, config)
			},
		)

		hc, err := ts.run(test.args)

		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Fatalf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}

		if !reflect.DeepEqual(hc.HTTP, test.expConfig) {
			t.Errorf("Test %v: unexpected config: \n\t%+v\n\t%+v", i, hc.HTTP, test.expConfig)
		}
		ts.tearDown()
	}
}

func TestWSFlags(t *testing.T) {
	tests := []struct {
		args      []string
		expConfig feechainconfig.WsConfig
		expErr    error
	}{
		{
			args:      []string{},
			expConfig: defaultConfig.WS,
		},
		{
			args: []string{"--ws=false"},
			expConfig: feechainconfig.WsConfig{
				Enabled:  false,
				IP:       defaultConfig.WS.IP,
				Port:     defaultConfig.WS.Port,
				AuthPort: defaultConfig.WS.AuthPort,
			},
		},
		{
			args: []string{"--ws", "--ws.ip", "8.8.8.8", "--ws.port", "9001"},
			expConfig: feechainconfig.WsConfig{
				Enabled:  true,
				IP:       "8.8.8.8",
				Port:     9001,
				AuthPort: defaultConfig.WS.AuthPort,
			},
		},
		{
			args: []string{"--ws", "--ws.ip", "8.8.8.8", "--ws.auth-port", "9001"},
			expConfig: feechainconfig.WsConfig{
				Enabled:  true,
				IP:       "8.8.8.8",
				Port:     defaultConfig.WS.Port,
				AuthPort: 9001,
			},
		},
		{
			args: []string{"--ip", "8.8.8.8", "--port", "9001", "--public_rpc"},
			expConfig: feechainconfig.WsConfig{
				Enabled:  true,
				IP:       nodeconfig.DefaultPublicListenIP,
				Port:     9801,
				AuthPort: 9802,
			},
		},
	}
	for i, test := range tests {
		ts := newFlagTestSuite(t, append(wsFlags, legacyMiscFlags...),
			func(cmd *cobra.Command, config *feechainconfig.FeechainConfig) {
				applyLegacyMiscFlags(cmd, config)
				applyWSFlags(cmd, config)
			},
		)

		hc, err := ts.run(test.args)

		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Fatalf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}

		if !reflect.DeepEqual(hc.WS, test.expConfig) {
			t.Errorf("Test %v: \n\t%+v\n\t%+v", i, hc.WS, test.expConfig)
		}
		ts.tearDown()
	}
}

func TestRPCOptFlags(t *testing.T) {
	tests := []struct {
		args      []string
		expConfig feechainconfig.RpcOptConfig
	}{
		{
			args: []string{"--rpc.debug"},
			expConfig: feechainconfig.RpcOptConfig{
				DebugEnabled:      true,
				RateLimterEnabled: true,
				RequestsPerSecond: 1000,
			},
		},

		{
			args: []string{},
			expConfig: feechainconfig.RpcOptConfig{
				DebugEnabled:      false,
				RateLimterEnabled: true,
				RequestsPerSecond: 1000,
			},
		},

		{
			args: []string{"--rpc.ratelimiter", "--rpc.ratelimit", "2000"},
			expConfig: feechainconfig.RpcOptConfig{
				DebugEnabled:      false,
				RateLimterEnabled: true,
				RequestsPerSecond: 2000,
			},
		},

		{
			args: []string{"--rpc.ratelimiter=false", "--rpc.ratelimit", "2000"},
			expConfig: feechainconfig.RpcOptConfig{
				DebugEnabled:      false,
				RateLimterEnabled: false,
				RequestsPerSecond: 2000,
			},
		},
	}
	for i, test := range tests {
		ts := newFlagTestSuite(t, rpcOptFlags, applyRPCOptFlags)

		hc, _ := ts.run(test.args)

		if !reflect.DeepEqual(hc.RPCOpt, test.expConfig) {
			t.Errorf("Test %v: \n\t%+v\n\t%+v", i, hc.RPCOpt, test.expConfig)
		}

		ts.tearDown()
	}
}

func TestBLSFlags(t *testing.T) {
	tests := []struct {
		args      []string
		expConfig feechainconfig.BlsConfig
		expErr    error
	}{
		{
			args:      []string{},
			expConfig: defaultConfig.BLSKeys,
		},
		{
			args: []string{"--bls.dir", "./.blskeys", "--bls.keys", "key1,key2",
				"--bls.maxkeys", "8", "--bls.pass", "--bls.pass.src", "auto", "--bls.pass.save",
				"--bls.kms", "--bls.kms.src", "shared",
			},
			expConfig: feechainconfig.BlsConfig{
				KeyDir:           "./.blskeys",
				KeyFiles:         []string{"key1", "key2"},
				MaxKeys:          8,
				PassEnabled:      true,
				PassSrcType:      "auto",
				PassFile:         "",
				SavePassphrase:   true,
				KMSEnabled:       true,
				KMSConfigSrcType: "shared",
				KMSConfigFile:    "",
			},
		},
		{
			args: []string{"--bls.pass.file", "xxx.pass", "--bls.kms.config", "config.json"},
			expConfig: feechainconfig.BlsConfig{
				KeyDir:           defaultConfig.BLSKeys.KeyDir,
				KeyFiles:         defaultConfig.BLSKeys.KeyFiles,
				MaxKeys:          defaultConfig.BLSKeys.MaxKeys,
				PassEnabled:      true,
				PassSrcType:      "file",
				PassFile:         "xxx.pass",
				SavePassphrase:   false,
				KMSEnabled:       false,
				KMSConfigSrcType: "file",
				KMSConfigFile:    "config.json",
			},
		},
		{
			args: []string{"--blskey_file", "key1,key2", "--blsfolder", "./.nkeys",
				"--max_bls_keys_per_node", "5", "--blspass", "file:xxx.pass", "--save-passphrase",
				"--aws-config-source", "file:config.json",
			},
			expConfig: feechainconfig.BlsConfig{
				KeyDir:           "./.nkeys",
				KeyFiles:         []string{"key1", "key2"},
				MaxKeys:          5,
				PassEnabled:      true,
				PassSrcType:      "file",
				PassFile:         "xxx.pass",
				SavePassphrase:   true,
				KMSEnabled:       false,
				KMSConfigSrcType: "file",
				KMSConfigFile:    "config.json",
			},
		},
	}
	for i, test := range tests {
		ts := newFlagTestSuite(t, blsFlags, applyBLSFlags)

		hc, err := ts.run(test.args)

		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Fatalf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}
		if !reflect.DeepEqual(hc.BLSKeys, test.expConfig) {
			t.Errorf("Test %v: \n\t%+v\n\t%+v", i, hc.BLSKeys, test.expConfig)
		}

		ts.tearDown()
	}
}

func TestConsensusFlags(t *testing.T) {
	tests := []struct {
		args      []string
		expConfig *feechainconfig.ConsensusConfig
		expErr    error
	}{
		{
			args:      []string{},
			expConfig: nil,
		},
		{
			args: []string{"--consensus.min-peers", "10", "--consensus.aggregate-sig=false"},
			expConfig: &feechainconfig.ConsensusConfig{
				MinPeers:     10,
				AggregateSig: false,
			},
		},
		{
			args: []string{"--delay_commit", "10ms", "--block_period", "5", "--min_peers", "10",
				"--consensus.aggregate-sig=true"},
			expConfig: &feechainconfig.ConsensusConfig{
				MinPeers:     10,
				AggregateSig: true,
			},
		},
	}
	for i, test := range tests {
		ts := newFlagTestSuite(t, consensusFlags, applyConsensusFlags)

		hc, err := ts.run(test.args)

		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Fatalf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}
		if !reflect.DeepEqual(hc.Consensus, test.expConfig) {
			t.Errorf("Test %v: unexpected config \n\t%+v\n\t%+v", i, hc.Consensus, test.expConfig)
		}

		ts.tearDown()
	}
}

func TestTxPoolFlags(t *testing.T) {
	tests := []struct {
		args      []string
		expConfig feechainconfig.TxPoolConfig
		expErr    error
	}{
		{
			args: []string{},
			expConfig: feechainconfig.TxPoolConfig{
				BlacklistFile:  defaultConfig.TxPool.BlacklistFile,
				RosettaFixFile: defaultConfig.TxPool.RosettaFixFile,
				AccountSlots:   defaultConfig.TxPool.AccountSlots,
			},
		},
		{
			args: []string{"--txpool.blacklist", "blacklist.file", "--txpool.rosettafixfile", "rosettafix.file"},
			expConfig: feechainconfig.TxPoolConfig{
				BlacklistFile:  "blacklist.file",
				RosettaFixFile: "rosettafix.file",
				AccountSlots:   16, // default
			},
		},
		{
			args: []string{"--blacklist", "blacklist.file", "--txpool.rosettafixfile", "rosettafix.file"},
			expConfig: feechainconfig.TxPoolConfig{
				BlacklistFile:  "blacklist.file",
				RosettaFixFile: "rosettafix.file",
				AccountSlots:   16, // default
			},
		},
		{
			args: []string{"--txpool.accountslots", "5", "--txpool.blacklist", "blacklist.file", "--txpool.rosettafixfile", "rosettafix.file"},
			expConfig: feechainconfig.TxPoolConfig{
				AccountSlots:   5,
				BlacklistFile:  "blacklist.file",
				RosettaFixFile: "rosettafix.file",
			},
		},
	}
	for i, test := range tests {
		ts := newFlagTestSuite(t, txPoolFlags, applyTxPoolFlags)
		hc, err := ts.run(test.args)

		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Fatalf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}

		if !reflect.DeepEqual(hc.TxPool, test.expConfig) {
			t.Errorf("Test %v: unexpected config\n\t%+v\n\t%+v", i, hc.TxPool, test.expConfig)
		}
		ts.tearDown()
	}
}

func TestPprofFlags(t *testing.T) {
	tests := []struct {
		args      []string
		expConfig feechainconfig.PprofConfig
		expErr    error
	}{
		{
			args:      []string{},
			expConfig: defaultConfig.Pprof,
		},
		{
			args: []string{"--pprof"},
			expConfig: feechainconfig.PprofConfig{
				Enabled:            true,
				ListenAddr:         defaultConfig.Pprof.ListenAddr,
				Folder:             defaultConfig.Pprof.Folder,
				ProfileNames:       defaultConfig.Pprof.ProfileNames,
				ProfileIntervals:   defaultConfig.Pprof.ProfileIntervals,
				ProfileDebugValues: defaultConfig.Pprof.ProfileDebugValues,
			},
		},
		{
			args: []string{"--pprof.addr", "8.8.8.8:9001"},
			expConfig: feechainconfig.PprofConfig{
				Enabled:            true,
				ListenAddr:         "8.8.8.8:9001",
				Folder:             defaultConfig.Pprof.Folder,
				ProfileNames:       defaultConfig.Pprof.ProfileNames,
				ProfileIntervals:   defaultConfig.Pprof.ProfileIntervals,
				ProfileDebugValues: defaultConfig.Pprof.ProfileDebugValues,
			},
		},
		{
			args: []string{"--pprof=false", "--pprof.addr", "8.8.8.8:9001"},
			expConfig: feechainconfig.PprofConfig{
				Enabled:            false,
				ListenAddr:         "8.8.8.8:9001",
				Folder:             defaultConfig.Pprof.Folder,
				ProfileNames:       defaultConfig.Pprof.ProfileNames,
				ProfileIntervals:   defaultConfig.Pprof.ProfileIntervals,
				ProfileDebugValues: defaultConfig.Pprof.ProfileDebugValues,
			},
		},
		{
			args: []string{"--pprof.profile.names", "cpu,heap,mutex"},
			expConfig: feechainconfig.PprofConfig{
				Enabled:            true,
				ListenAddr:         defaultConfig.Pprof.ListenAddr,
				Folder:             defaultConfig.Pprof.Folder,
				ProfileNames:       []string{"cpu", "heap", "mutex"},
				ProfileIntervals:   defaultConfig.Pprof.ProfileIntervals,
				ProfileDebugValues: defaultConfig.Pprof.ProfileDebugValues,
			},
		},
		{
			args: []string{"--pprof.profile.intervals", "0,1"},
			expConfig: feechainconfig.PprofConfig{
				Enabled:            true,
				ListenAddr:         defaultConfig.Pprof.ListenAddr,
				Folder:             defaultConfig.Pprof.Folder,
				ProfileNames:       defaultConfig.Pprof.ProfileNames,
				ProfileIntervals:   []int{0, 1},
				ProfileDebugValues: defaultConfig.Pprof.ProfileDebugValues,
			},
		},
		{
			args: []string{"--pprof.profile.debug", "0,1,0"},
			expConfig: feechainconfig.PprofConfig{
				Enabled:            true,
				ListenAddr:         defaultConfig.Pprof.ListenAddr,
				Folder:             defaultConfig.Pprof.Folder,
				ProfileNames:       defaultConfig.Pprof.ProfileNames,
				ProfileIntervals:   defaultConfig.Pprof.ProfileIntervals,
				ProfileDebugValues: []int{0, 1, 0},
			},
		},
	}
	for i, test := range tests {
		ts := newFlagTestSuite(t, pprofFlags, applyPprofFlags)
		hc, err := ts.run(test.args)

		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Fatalf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}
		if !reflect.DeepEqual(hc.Pprof, test.expConfig) {
			t.Errorf("Test %v: unexpected config\n\t%+v\n\t%+v", i, hc.Pprof, test.expConfig)
		}
		ts.tearDown()
	}
}

func TestLogFlags(t *testing.T) {
	tests := []struct {
		args      []string
		expConfig feechainconfig.LogConfig
		expErr    error
	}{
		{
			args:      []string{},
			expConfig: defaultConfig.Log,
		},
		{
			args: []string{"--log.dir", "latest_log", "--log.max-size", "10", "--log.rotate-count", "3",
				"--log.rotate-max-age", "0", "--log.name", "feechain.log", "--log.verb", "5",
				"--log.verbose-prints", "config"},
			expConfig: feechainconfig.LogConfig{
				Folder:       "latest_log",
				FileName:     "feechain.log",
				RotateSize:   10,
				RotateCount:  3,
				RotateMaxAge: 0,
				Verbosity:    5,
				VerbosePrints: feechainconfig.LogVerbosePrints{
					Config: true,
				},
				Context: nil,
			},
		},
		{
			args: []string{"--log.ctx.ip", "8.8.8.8", "--log.ctx.port", "9001"},
			expConfig: feechainconfig.LogConfig{
				Folder:        defaultConfig.Log.Folder,
				FileName:      defaultConfig.Log.FileName,
				RotateSize:    defaultConfig.Log.RotateSize,
				RotateCount:   defaultConfig.Log.RotateCount,
				RotateMaxAge:  defaultConfig.Log.RotateMaxAge,
				Verbosity:     defaultConfig.Log.Verbosity,
				VerbosePrints: defaultConfig.Log.VerbosePrints,
				Context: &feechainconfig.LogContext{
					IP:   "8.8.8.8",
					Port: 9001,
				},
			},
		},
		{
			args: []string{"--log_folder", "latest_log", "--log_max_size", "10", "--verbosity",
				"5", "--ip", "8.8.8.8", "--port", "9001"},
			expConfig: feechainconfig.LogConfig{
				Folder:        "latest_log",
				FileName:      "validator-8.8.8.8-9001.log",
				RotateSize:    10,
				RotateCount:   0,
				RotateMaxAge:  0,
				Verbosity:     5,
				VerbosePrints: defaultConfig.Log.VerbosePrints,
				Context: &feechainconfig.LogContext{
					IP:   "8.8.8.8",
					Port: 9001,
				},
			},
		},
	}
	for i, test := range tests {
		ts := newFlagTestSuite(t, append(logFlags, legacyMiscFlags...),
			func(cmd *cobra.Command, config *feechainconfig.FeechainConfig) {
				applyLegacyMiscFlags(cmd, config)
				applyLogFlags(cmd, config)
			},
		)
		hc, err := ts.run(test.args)

		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Fatalf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}

		if !reflect.DeepEqual(hc.Log, test.expConfig) {
			t.Errorf("Test %v:\n\t%+v\n\t%+v", i, hc.Log, test.expConfig)
		}
		ts.tearDown()
	}
}

func TestSysFlags(t *testing.T) {
	tests := []struct {
		args      []string
		expConfig *feechainconfig.SysConfig
		expErr    error
	}{
		{
			args: []string{},
			expConfig: &feechainconfig.SysConfig{
				NtpServer: defaultSysConfig.NtpServer,
			},
		},
		{
			args: []string{"--sys.ntp", "0.pool.ntp.org"},
			expConfig: &feechainconfig.SysConfig{
				NtpServer: "0.pool.ntp.org",
			},
		},
	}

	for i, test := range tests {
		ts := newFlagTestSuite(t, sysFlags, applySysFlags)
		hc, err := ts.run(test.args)

		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Fatalf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}

		if !reflect.DeepEqual(hc.Sys, test.expConfig) {
			t.Errorf("Test %v:\n\t%+v\n\t%+v", i, hc.Sys, test.expConfig)
		}
		ts.tearDown()
	}
}

func TestDevnetFlags(t *testing.T) {
	tests := []struct {
		args      []string
		expConfig *feechainconfig.DevnetConfig
		expErr    error
	}{
		{
			args:      []string{},
			expConfig: nil,
		},
		{
			args: []string{"--devnet.num-shard", "3", "--devnet.shard-size", "100",
				"--devnet.fch-node-size", "60"},
			expConfig: &feechainconfig.DevnetConfig{
				NumShards:   3,
				ShardSize:   100,
				FchNodeSize: 60,
			},
		},
		{
			args: []string{"--dn_num_shards", "3", "--dn_shard_size", "100", "--dn_fch_size",
				"60"},
			expConfig: &feechainconfig.DevnetConfig{
				NumShards:   3,
				ShardSize:   100,
				FchNodeSize: 60,
			},
		},
	}
	for i, test := range tests {
		ts := newFlagTestSuite(t, devnetFlags, applyDevnetFlags)
		hc, err := ts.run(test.args)

		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Fatalf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}

		if !reflect.DeepEqual(hc.Devnet, test.expConfig) {
			t.Errorf("Test %v:\n\t%+v\n\t%+v", i, hc.Devnet, test.expConfig)
		}
		ts.tearDown()
	}
}

func TestRevertFlags(t *testing.T) {
	tests := []struct {
		args      []string
		expConfig *feechainconfig.RevertConfig
		expErr    error
	}{
		{
			args:      []string{},
			expConfig: nil,
		},
		{
			args: []string{"--revert.beacon"},
			expConfig: &feechainconfig.RevertConfig{
				RevertBeacon: true,
				RevertTo:     defaultRevertConfig.RevertTo,
				RevertBefore: defaultRevertConfig.RevertBefore,
			},
		},
		{
			args: []string{"--revert.beacon", "--revert.to", "100", "--revert.do-before", "10000"},
			expConfig: &feechainconfig.RevertConfig{
				RevertBeacon: true,
				RevertTo:     100,
				RevertBefore: 10000,
			},
		},
		{
			args: []string{"--revert_beacon", "--do_revert_before", "10000", "--revert_to", "100"},
			expConfig: &feechainconfig.RevertConfig{
				RevertBeacon: true,
				RevertTo:     100,
				RevertBefore: 10000,
			},
		},
	}
	for i, test := range tests {
		ts := newFlagTestSuite(t, revertFlags, applyRevertFlags)
		hc, err := ts.run(test.args)

		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Fatalf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}
		if !reflect.DeepEqual(hc.Revert, test.expConfig) {
			t.Errorf("Test %v:\n\t%+v\n\t%+v", i, hc.Revert, test.expConfig)
		}
		ts.tearDown()
	}
}

func TestDNSSyncFlags(t *testing.T) {
	tests := []struct {
		args      []string
		network   string
		expConfig feechainconfig.DnsSync
		expErr    error
	}{
		{
			args:      []string{},
			network:   "asadal",
			expConfig: getDefaultDNSSyncConfig(nodeconfig.Mainnet),
		},
		{
			args:      []string{"--sync.legacy.server", "--sync.legacy.client"},
			network:   "asadal",
			expConfig: getDefaultDNSSyncConfig(nodeconfig.Mainnet),
		},
		{
			args:    []string{"--sync.legacy.server", "--sync.legacy.client"},
			network: "testnet",
			expConfig: func() feechainconfig.DnsSync {
				cfg := getDefaultDNSSyncConfig(nodeconfig.Mainnet)
				cfg.Client = true
				cfg.Server = true
				return cfg
			}(),
		},
		{
			args:      []string{"--dns.server", "--dns.client"},
			network:   "asadal",
			expConfig: getDefaultDNSSyncConfig(nodeconfig.Mainnet),
		},
	}

	for i, test := range tests {
		ts := newFlagTestSuite(t, dnsSyncFlags, func(command *cobra.Command, config *feechainconfig.FeechainConfig) {
			config.Network.NetworkType = test.network
			applyDNSSyncFlags(command, config)
		})
		hc, err := ts.run(test.args)

		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Fatalf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}
		if !reflect.DeepEqual(hc.DNSSync, test.expConfig) {
			t.Errorf("Test %v:\n\t%+v\n\t%+v", i, hc.DNSSync, test.expConfig)
		}

		ts.tearDown()
	}
}
func TestSyncFlags(t *testing.T) {
	tests := []struct {
		args      []string
		network   string
		expConfig feechainconfig.SyncConfig
		expErr    error
	}{
		{
			args: []string{"--sync", "--sync.downloader", "--sync.concurrency", "10", "--sync.min-peers", "10",
				"--sync.init-peers", "10", "--sync.disc.soft-low-cap", "10",
				"--sync.disc.hard-low-cap", "10", "--sync.disc.hi-cap", "10",
				"--sync.disc.batch", "10",
			},
			network: "asadal",
			expConfig: func() feechainconfig.SyncConfig {
				cfgSync := defaultMainnetSyncConfig
				cfgSync.Enabled = true
				cfgSync.Downloader = true
				cfgSync.Concurrency = 10
				cfgSync.MinPeers = 10
				cfgSync.InitStreams = 10
				cfgSync.DiscSoftLowCap = 10
				cfgSync.DiscHardLowCap = 10
				cfgSync.DiscHighCap = 10
				cfgSync.DiscBatch = 10
				return cfgSync
			}(),
		},
	}
	for i, test := range tests {
		ts := newFlagTestSuite(t, syncFlags, func(command *cobra.Command, config *feechainconfig.FeechainConfig) {
			applySyncFlags(command, config)
		})
		hc, err := ts.run(test.args)

		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Fatalf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}
		if !reflect.DeepEqual(hc.Sync, test.expConfig) {
			t.Errorf("Test %v:\n\t%+v\n\t%+v", i, hc.Sync, test.expConfig)
		}

		ts.tearDown()
	}
}

func TestShardDataFlags(t *testing.T) {
	tests := []struct {
		args      []string
		expConfig feechainconfig.ShardDataConfig
		expErr    error
	}{
		{
			args:      []string{},
			expConfig: defaultConfig.ShardData,
		},
		{
			args: []string{"--sharddata.enable",
				"--sharddata.disk_count", "8",
				"--sharddata.shard_count", "4",
				"--sharddata.cache_time", "10",
				"--sharddata.cache_size", "512",
			},
			expConfig: feechainconfig.ShardDataConfig{
				EnableShardData: true,
				DiskCount:       8,
				ShardCount:      4,
				CacheTime:       10,
				CacheSize:       512,
			},
		},
	}
	for i, test := range tests {
		ts := newFlagTestSuite(t, shardDataFlags, func(command *cobra.Command, config *feechainconfig.FeechainConfig) {
			applyShardDataFlags(command, config)
		})
		hc, err := ts.run(test.args)

		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Fatalf("Test %v: %v", i, assErr)
		}
		if err != nil || test.expErr != nil {
			continue
		}
		if !reflect.DeepEqual(hc.ShardData, test.expConfig) {
			t.Errorf("Test %v:\n\t%+v\n\t%+v", i, hc.ShardData, test.expConfig)
		}

		ts.tearDown()
	}
}

type flagTestSuite struct {
	t *testing.T

	cmd *cobra.Command
	hc  feechainconfig.FeechainConfig
}

func newFlagTestSuite(t *testing.T, flags []cli.Flag, applyFlags func(*cobra.Command, *feechainconfig.FeechainConfig)) *flagTestSuite {
	cli.SetParseErrorHandle(func(err error) { t.Fatal(err) })

	ts := &flagTestSuite{hc: getDefaultFchConfigCopy(defNetworkType)}
	ts.cmd = makeTestCommand(func(cmd *cobra.Command, args []string) {
		applyFlags(cmd, &ts.hc)
	})
	if err := cli.RegisterFlags(ts.cmd, flags); err != nil {
		t.Fatal(err)
	}

	return ts
}

func (ts *flagTestSuite) run(args []string) (feechainconfig.FeechainConfig, error) {
	ts.cmd.SetArgs(args)
	err := ts.cmd.Execute()
	return ts.hc, err
}

func (ts *flagTestSuite) tearDown() {
	cli.SetParseErrorHandle(func(error) {})
}

func makeTestCommand(run func(cmd *cobra.Command, args []string)) *cobra.Command {
	return &cobra.Command{
		Use: "test",
		Run: run,
	}
}

func assertError(gotErr, expErr error) error {
	if (gotErr == nil) != (expErr == nil) {
		return fmt.Errorf("error unexpected [%v] / [%v]", gotErr, expErr)
	}
	if gotErr == nil {
		return nil
	}
	if !strings.Contains(gotErr.Error(), expErr.Error()) {
		return fmt.Errorf("error unexpected [%v] / [%v]", gotErr, expErr)
	}
	return nil
}

func intPtr(i int) *int {
	return &i
}
