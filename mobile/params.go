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

// Contains all the wrappers from the params package.

package glbchain-dev

import (
	"encoding/json"

	"github.com/lbchain-devchain/go-lbchain-dev/core"
	"github.com/lbchain-devchain/go-lbchain-dev/p2p/discv5"
	"github.com/lbchain-devchain/go-lbchain-dev/params"
)

// MainnetGenesis returns the JSON spec to use for the main lbchain-devchain network. It
// is actually empty since that defaults to the hard coded binary genesis block.
func MainnetGenesis() string {
	return ""
}

// TestnetGenesis returns the JSON spec to use for the lbchain-devchain test network.
func TestnetGenesis() string {
	enc, err := json.Marshal(core.DefaultTestnetGenesisBlock())
	if err != nil {
		panic(err)
	}
	return string(enc)
}

// RinkebyGenesis returns the JSON spec to use for the Rinkeby test network
func RinkebyGenesis() string {
	enc, err := json.Marshal(core.DefaultRinkebyGenesisBlock())
	if err != nil {
		panic(err)
	}
	return string(enc)
}

// FoundationBootnodes returns the enode URLs of the P2P boolbchain-devrap nodes operated
// by the foundation running the V5 discovery protocol.
func FoundationBootnodes() *Enodes {
	nodes := &Enodes{nodes: make([]*discv5.Node, len(params.DiscoveryV5Bootnodes))}
	for i, url := range params.DiscoveryV5Bootnodes {
		nodes.nodes[i] = discv5.MustParseNode(url)
	}
	return nodes
}
