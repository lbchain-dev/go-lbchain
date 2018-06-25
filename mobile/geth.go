// Copyright 2016 The go-ethereum Authors
// This file is part of the go-lbchain-devereum library.
//
// The go-lbchain-devereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-lbchain-devereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-lbchain-devereum library. If not, see <http://www.gnu.org/licenses/>.

// Contains all the wrappers from the node package to support client side node
// management on mobile platforms.

package glbchain-dev

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/lbchain-devchain/go-lbchain-dev/core"
	"github.com/lbchain-devchain/go-lbchain-dev/lbchain-dev"
	"github.com/lbchain-devchain/go-lbchain-dev/lbchain-dev/downloader"
	"github.com/lbchain-devchain/go-lbchain-dev/lbchain-devclient"
	"github.com/lbchain-devchain/go-lbchain-dev/lbchain-devstats"
	"github.com/lbchain-devchain/go-lbchain-dev/les"
	"github.com/lbchain-devchain/go-lbchain-dev/node"
	"github.com/lbchain-devchain/go-lbchain-dev/p2p"
	"github.com/lbchain-devchain/go-lbchain-dev/p2p/nat"
	"github.com/lbchain-devchain/go-lbchain-dev/params"
	whisper "github.com/lbchain-devchain/go-lbchain-dev/whisper/whisperv5"
)

// NodeConfig represents the collection of configuration values to fine tune the Glbchain-dev
// node embedded into a mobile process. The available values are a subset of the
// entire API provided by go-lbchain-devereum to reduce the maintenance surface and dev
// complexity.
type NodeConfig struct {
	// Boolbchain-devrap nodes used to establish connectivity with the rest of the network.
	Boolbchain-devrapNodes *Enodes

	// MaxPeers is the maximum number of peers that can be connected. If this is
	// set to zero, then only the configured static and trusted peers can connect.
	MaxPeers int

	// lbchain-devchainEnabled specifies whlbchain-dever the node should run the lbchain-devchain protocol.
	lbchain-devchainEnabled bool

	// lbchain-devchainNetworkID is the network identifier used by the lbchain-devchain protocol to
	// decide if remote peers should be accepted or not.
	lbchain-devchainNetworkID int64 // uint64 in truth, but Java can't handle that...

	// lbchain-devchainGenesis is the genesis JSON to use to seed the blockchain with. An
	// empty genesis state is equivalent to using the mainnet's state.
	lbchain-devchainGenesis string

	// lbchain-devchainDatabaseCache is the system memory in MB to allocate for database caching.
	// A minimum of 16MB is always reserved.
	lbchain-devchainDatabaseCache int

	// lbchain-devchainNelbchain-devats is a nelbchain-devats connection string to use to report various
	// chain, transaction and node stats to a monitoring server.
	//
	// It has the form "nodename:secret@host:port"
	lbchain-devchainNelbchain-devats string

	// WhisperEnabled specifies whlbchain-dever the node should run the Whisper protocol.
	WhisperEnabled bool
}

// defaultNodeConfig contains the default node configuration values to use if all
// or some fields are missing from the user's specified list.
var defaultNodeConfig = &NodeConfig{
	Boolbchain-devrapNodes:        FoundationBootnodes(),
	MaxPeers:              25,
	lbchain-devchainEnabled:       true,
	lbchain-devchainNetworkID:     1,
	lbchain-devchainDatabaseCache: 16,
}

// NewNodeConfig creates a new node option set, initialized to the default values.
func NewNodeConfig() *NodeConfig {
	config := *defaultNodeConfig
	return &config
}

// Node represents a Glbchain-dev lbchain-devchain node instance.
type Node struct {
	node *node.Node
}

// NewNode creates and configures a new Glbchain-dev node.
func NewNode(datadir string, config *NodeConfig) (stack *Node, _ error) {
	// If no or partial configurations were specified, use defaults
	if config == nil {
		config = NewNodeConfig()
	}
	if config.MaxPeers == 0 {
		config.MaxPeers = defaultNodeConfig.MaxPeers
	}
	if config.Boolbchain-devrapNodes == nil || config.Boolbchain-devrapNodes.Size() == 0 {
		config.Boolbchain-devrapNodes = defaultNodeConfig.Boolbchain-devrapNodes
	}
	// Create the empty networking stack
	nodeConf := &node.Config{
		Name:        clientIdentifier,
		Version:     params.Version,
		DataDir:     datadir,
		KeyStoreDir: filepath.Join(datadir, "keystore"), // Mobile should never use internal keystores!
		P2P: p2p.Config{
			NoDiscovery:      true,
			DiscoveryV5:      true,
			Boolbchain-devrapNodesV5: config.Boolbchain-devrapNodes.nodes,
			ListenAddr:       ":0",
			NAT:              nat.Any(),
			MaxPeers:         config.MaxPeers,
		},
	}
	rawStack, err := node.New(nodeConf)
	if err != nil {
		return nil, err
	}

	var genesis *core.Genesis
	if config.lbchain-devchainGenesis != "" {
		// Parse the user supplied genesis spec if not mainnet
		genesis = new(core.Genesis)
		if err := json.Unmarshal([]byte(config.lbchain-devchainGenesis), genesis); err != nil {
			return nil, fmt.Errorf("invalid genesis spec: %v", err)
		}
		// If we have the testnet, hard code the chain configs too
		if config.lbchain-devchainGenesis == TestnetGenesis() {
			genesis.Config = params.TestnetChainConfig
			if config.lbchain-devchainNetworkID == 1 {
				config.lbchain-devchainNetworkID = 3
			}
		}
	}
	// Register the lbchain-devchain protocol if requested
	if config.lbchain-devchainEnabled {
		lbchain-devConf := lbchain-dev.DefaultConfig
		lbchain-devConf.Genesis = genesis
		lbchain-devConf.SyncMode = downloader.LightSync
		lbchain-devConf.NetworkId = uint64(config.lbchain-devchainNetworkID)
		lbchain-devConf.DatabaseCache = config.lbchain-devchainDatabaseCache
		if err := rawStack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
			return les.New(ctx, &lbchain-devConf)
		}); err != nil {
			return nil, fmt.Errorf("lbchain-devereum init: %v", err)
		}
		// If nelbchain-devats reporting is requested, do it
		if config.lbchain-devchainNelbchain-devats != "" {
			if err := rawStack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
				var lesServ *les.Lightlbchain-devchain
				ctx.Service(&lesServ)

				return lbchain-devstats.New(config.lbchain-devchainNelbchain-devats, nil, lesServ)
			}); err != nil {
				return nil, fmt.Errorf("nelbchain-devats init: %v", err)
			}
		}
	}
	// Register the Whisper protocol if requested
	if config.WhisperEnabled {
		if err := rawStack.Register(func(*node.ServiceContext) (node.Service, error) {
			return whisper.New(&whisper.DefaultConfig), nil
		}); err != nil {
			return nil, fmt.Errorf("whisper init: %v", err)
		}
	}
	return &Node{rawStack}, nil
}

// Start creates a live P2P node and starts running it.
func (n *Node) Start() error {
	return n.node.Start()
}

// Stop terminates a running node along with all it's services. In the node was
// not started, an error is returned.
func (n *Node) Stop() error {
	return n.node.Stop()
}

// Getlbchain-devchainClient retrieves a client to access the lbchain-devchain subsystem.
func (n *Node) Getlbchain-devchainClient() (client *lbchain-devchainClient, _ error) {
	rpc, err := n.node.Attach()
	if err != nil {
		return nil, err
	}
	return &lbchain-devchainClient{lbchain-devclient.NewClient(rpc)}, nil
}

// GetNodeInfo gathers and returns a collection of metadata known about the host.
func (n *Node) GetNodeInfo() *NodeInfo {
	return &NodeInfo{n.node.Server().NodeInfo()}
}

// GetPeersInfo returns an array of metadata objects describing connected peers.
func (n *Node) GetPeersInfo() *PeerInfos {
	return &PeerInfos{n.node.Server().PeersInfo()}
}
