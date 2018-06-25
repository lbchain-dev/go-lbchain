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

// Package les implements the Light lbchain-devchain Subprotocol.
package les

import (
	"fmt"
	"sync"
	"time"

	"github.com/lbchain-devchain/go-lbchain-dev/accounts"
	"github.com/lbchain-devchain/go-lbchain-dev/common"
	"github.com/lbchain-devchain/go-lbchain-dev/common/hexutil"
	"github.com/lbchain-devchain/go-lbchain-dev/consensus"
	"github.com/lbchain-devchain/go-lbchain-dev/core"
	"github.com/lbchain-devchain/go-lbchain-dev/core/bloombits"
	"github.com/lbchain-devchain/go-lbchain-dev/core/types"
	"github.com/lbchain-devchain/go-lbchain-dev/lbchain-dev"
	"github.com/lbchain-devchain/go-lbchain-dev/lbchain-dev/downloader"
	"github.com/lbchain-devchain/go-lbchain-dev/lbchain-dev/filters"
	"github.com/lbchain-devchain/go-lbchain-dev/lbchain-dev/gasprice"
	"github.com/lbchain-devchain/go-lbchain-dev/lbchain-devdb"
	"github.com/lbchain-devchain/go-lbchain-dev/event"
	"github.com/lbchain-devchain/go-lbchain-dev/internal/ethapi"
	"github.com/lbchain-devchain/go-lbchain-dev/light"
	"github.com/lbchain-devchain/go-lbchain-dev/log"
	"github.com/lbchain-devchain/go-lbchain-dev/node"
	"github.com/lbchain-devchain/go-lbchain-dev/p2p"
	"github.com/lbchain-devchain/go-lbchain-dev/p2p/discv5"
	"github.com/lbchain-devchain/go-lbchain-dev/params"
	rpc "github.com/lbchain-devchain/go-lbchain-dev/rpc"
)

type Lightlbchain-devchain struct {
	config *lbchain-dev.Config

	odr         *LesOdr
	relay       *LesTxRelay
	chainConfig *params.ChainConfig
	// Channel for shutting down the service
	shutdownChan chan bool
	// Handlers
	peers           *peerSet
	txPool          *light.TxPool
	blockchain      *light.LightChain
	protocolManager *ProtocolManager
	serverPool      *serverPool
	reqDist         *requestDistributor
	retriever       *retrieveManager
	// DB interfaces
	chainDb lbchain-devdb.Database // Block chain database

	bloomRequests                              chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer, chtIndexer, bloomTrieIndexer *core.ChainIndexer

	ApiBackend *LesApiBackend

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	networkId     uint64
	netRPCService *ethapi.PublicNetAPI

	wg sync.WaitGroup
}

func New(ctx *node.ServiceContext, config *lbchain-dev.Config) (*Lightlbchain-devchain, error) {
	chainDb, err := lbchain-dev.CreateDB(ctx, config, "lightchaindata")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)
	if _, isCompat := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !isCompat {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	peers := newPeerSet()
	quitSync := make(chan struct{})

	llbchain-dev := &Lightlbchain-devchain{
		config:           config,
		chainConfig:      chainConfig,
		chainDb:          chainDb,
		eventMux:         ctx.EventMux,
		peers:            peers,
		reqDist:          newRequestDistributor(peers, quitSync),
		accountManager:   ctx.AccountManager,
		engine:           lbchain-dev.CreateConsensusEngine(ctx, &config.lbchain-devash, chainConfig, chainDb),
		shutdownChan:     make(chan bool),
		networkId:        config.NetworkId,
		bloomRequests:    make(chan chan *bloombits.Retrieval),
		bloomIndexer:     lbchain-dev.NewBloomIndexer(chainDb, light.BloomTrieFrequency),
		chtIndexer:       light.NewChtIndexer(chainDb, true),
		bloomTrieIndexer: light.NewBloomTrieIndexer(chainDb, true),
	}

	llbchain-dev.relay = NewLesTxRelay(peers, llbchain-dev.reqDist)
	llbchain-dev.serverPool = newServerPool(chainDb, quitSync, &llbchain-dev.wg)
	llbchain-dev.retriever = newRetrieveManager(peers, llbchain-dev.reqDist, llbchain-dev.serverPool)
	llbchain-dev.odr = NewLesOdr(chainDb, llbchain-dev.chtIndexer, llbchain-dev.bloomTrieIndexer, llbchain-dev.bloomIndexer, llbchain-dev.retriever)
	if llbchain-dev.blockchain, err = light.NewLightChain(llbchain-dev.odr, llbchain-dev.chainConfig, llbchain-dev.engine); err != nil {
		return nil, err
	}
	llbchain-dev.bloomIndexer.Start(llbchain-dev.blockchain)
	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		llbchain-dev.blockchain.SetHead(compat.RewindTo)
		core.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}

	llbchain-dev.txPool = light.NewTxPool(llbchain-dev.chainConfig, llbchain-dev.blockchain, llbchain-dev.relay)
	if llbchain-dev.protocolManager, err = NewProtocolManager(llbchain-dev.chainConfig, true, ClientProtocolVersions, config.NetworkId, llbchain-dev.eventMux, llbchain-dev.engine, llbchain-dev.peers, llbchain-dev.blockchain, nil, chainDb, llbchain-dev.odr, llbchain-dev.relay, quitSync, &llbchain-dev.wg); err != nil {
		return nil, err
	}
	llbchain-dev.ApiBackend = &LesApiBackend{llbchain-dev, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.GasPrice
	}
	llbchain-dev.ApiBackend.gpo = gasprice.NewOracle(llbchain-dev.ApiBackend, gpoParams)
	return llbchain-dev, nil
}

func lesTopic(genesisHash common.Hash, protocolVersion uint) discv5.Topic {
	var name string
	switch protocolVersion {
	case lpv1:
		name = "LES"
	case lpv2:
		name = "LES2"
	default:
		panic(nil)
	}
	return discv5.Topic(name + "@" + common.Bytes2Hex(genesisHash.Bytes()[0:8]))
}

type LightDummyAPI struct{}

// lbchain-deverbase is the address that mining rewards will be send to
func (s *LightDummyAPI) lbchain-deverbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("not supported")
}

// Coinbase is the address that mining rewards will be send to (alias for lbchain-deverbase)
func (s *LightDummyAPI) Coinbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("not supported")
}

// Hashrate returns the POW hashrate
func (s *LightDummyAPI) Hashrate() hexutil.Uint {
	return 0
}

// Mining returns an indication if this node is currently mining.
func (s *LightDummyAPI) Mining() bool {
	return false
}

// APIs returns the collection of RPC services the lbchain-devereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *Lightlbchain-devchain) APIs() []rpc.API {
	return append(ethapi.GetAPIs(s.ApiBackend), []rpc.API{
		{
			Namespace: "lbchain-dev",
			Version:   "1.0",
			Service:   &LightDummyAPI{},
			Public:    true,
		}, {
			Namespace: "lbchain-dev",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "lbchain-dev",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.ApiBackend, true),
			Public:    true,
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *Lightlbchain-devchain) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *Lightlbchain-devchain) BlockChain() *light.LightChain      { return s.blockchain }
func (s *Lightlbchain-devchain) TxPool() *light.TxPool              { return s.txPool }
func (s *Lightlbchain-devchain) Engine() consensus.Engine           { return s.engine }
func (s *Lightlbchain-devchain) LesVersion() int                    { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *Lightlbchain-devchain) Downloader() *downloader.Downloader { return s.protocolManager.downloader }
func (s *Lightlbchain-devchain) EventMux() *event.TypeMux           { return s.eventMux }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *Lightlbchain-devchain) Protocols() []p2p.Protocol {
	return s.protocolManager.SubProtocols
}

// Start implements node.Service, starting all internal goroutines needed by the
// lbchain-devchain protocol implementation.
func (s *Lightlbchain-devchain) Start(srvr *p2p.Server) error {
	s.startBloomHandlers()
	log.Warn("Light client mode is an experimental feature")
	s.netRPCService = ethapi.NewPublicNetAPI(srvr, s.networkId)
	// clients are searching for the first advertised protocol in the list
	protocolVersion := AdvertiseProtocolVersions[0]
	s.serverPool.start(srvr, lesTopic(s.blockchain.Genesis().Hash(), protocolVersion))
	s.protocolManager.Start(s.config.LightPeers)
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// lbchain-devchain protocol.
func (s *Lightlbchain-devchain) Stop() error {
	s.odr.Stop()
	if s.bloomIndexer != nil {
		s.bloomIndexer.Close()
	}
	if s.chtIndexer != nil {
		s.chtIndexer.Close()
	}
	if s.bloomTrieIndexer != nil {
		s.bloomTrieIndexer.Close()
	}
	s.blockchain.Stop()
	s.protocolManager.Stop()
	s.txPool.Stop()

	s.eventMux.Stop()

	time.Sleep(time.Millisecond * 200)
	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}
