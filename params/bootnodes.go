// Copyright 2015 The go-ethereum Authors
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

package params

// MainnetBootnodes are the enode URLs of the P2P boolbchain-devrap nodes running on
// the main lbchain-devchain network.
var MainnetBootnodes = []string{
	// lbchain-devchain Foundation Go Bootnodes
	"enode://2a374a238067c41ed0cfb6c2f0b0361123b1b2e3286a1769b51981cac9200b1245e82e3769509311f04459bcb67a1d97a6ffbab980e99f42e48b3cb7df58b733@47.74.229.190:8080",
}

// TestnetBootnodes are the enode URLs of the P2P boolbchain-devrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{
}

// RinkebyBootnodes are the enode URLs of the P2P boolbchain-devrap nodes running on the
// Rinkeby test network.
var RinkebyBootnodes = []string{
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P boolbchain-devrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
}
