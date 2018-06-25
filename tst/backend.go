// Copyright 2014 The go-ethereum Authors
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

// Package lbchain-dev implements the lbchain-devchain protocol.
package lbchain-dev

import (
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/lbchain-devchain/go-lbchain-dev/accounts"
	"github.com/lbchain-devchain/go-lbchain-dev/common"
	"github.com/lbchain-devchain/go-lbchain-dev/common/hexutil"
	"github.com/lbchain-devchain/go-lbchain-dev/consensus"
	"github.com/lbchain-devchain/go-lbchain-dev/consensus/clique"
	"github.com/lbchain-devchain/go-lbchain-dev/consensus/ethash"
	"github.com/lbchain-devchain/go-lbchain-dev/core"
	"github.com/lbchain-devchain/go-lbchain-dev/core/bloombits"
	"github.com/lbchain-devchain/go-lbchain-dev/core/types"
	"github.com/lbchain-devchain/go-lbchain-dev/core/vm"
	"github.com/lbchain-devchain/go-lbchain-dev/lbchain-dev/downloader"
	"github.com/lbchain-devchain/go-lbchain-dev/lbchain-dev/filters"
	"github.com/lbchain-devchain/go-lbchain-dev/lbchain-dev/gasprice"
	"github.com/lbchain-devchain/go-lbchain-dev/lbchain-devdb"
	"github.com/lbchain-devchain/go-lbchain-dev/event"
	"github.com/lbchain-devchain/go-lbchain-dev/internal/ethapi"
	"github.com/lbchain-devchain/go-lbchain-dev/log"
	"github.com/lbchain-devchain/go-lbchain-dev/miner"
	"github.com/lbchain-devchain/go-lbchain-dev/node"
	"github.com/lbchain-devchain/go-lbchain-dev/p2p"
	"github.com/lbchain-devchain/go-lbchain-dev/params"
	"github.com/lbchain-devchain/go-lbchain-dev/rlp"
	"github.com/lbchain-devchain/go-lbchain-dev/rpc"
)

type LesServer interface {
	Start(srvr *p2p.Server)
	Stop()
	Protocols() []p2p.Protocol
	SetBloomBitsIndexer(bbIndexer *core.ChainIndexer)
}

// lbchain-devchain implements the lbchain-devchain full node service.
type lbchain-devchain struct {
	config      *Config
	chainConfig *params.ChainConfig

	// Channel for shutting down the service
	shutdownChan  chan bool    // Channel for shutting down the lbchain-devereum
	stopDbUpgrade func() error // stop chain db sequential key upgrade

	// Handlers
	txPool          *core.TxPool
	blockchain      *core.BlockChain
	protocolManager *ProtocolManager
	lesServer       LesServer

	// DB interfaces
	chainDb lbchain-devdb.Database // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer             // Bloom indexer operating during block imports

	ApiBackend *lbchain-devApiBackend

	miner     *miner.Miner
	gasPrice  *big.Int
	lbchain-deverbase common.Address

	networkId     uint64
	netRPCService *ethapi.PublicNetAPI

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and lbchain-deverbase)
}

func (s *lbchain-devchain) AddLesServer(ls LesServer) {
	s.lesServer = ls
	ls.SetBloomBitsIndexer(s.bloomIndexer)
}

// New creates a new lbchain-devchain object (including the
// initialisation of the common lbchain-devchain object)
func New(ctx *node.ServiceContext, config *Config) (*lbchain-devchain, error) {
	if config.SyncMode == downloader.LightSync {
		return nil, errors.New("can't run lbchain-dev.lbchain-devchain in light sync mode, use les.Lightlbchain-devchain")
	}
	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	chainDb, err := CreateDB(ctx, config, "chaindata")
	if err != nil {
		return nil, err
	}
	stopDbUpgrade := upgradeDeduplicateData(chainDb)
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	lbchain-dev := &lbchain-devchain{
		config:         config,
		chainDb:        chainDb,
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		engine:         CreateConsensusEngine(ctx, &config.lbchain-devash, chainConfig, chainDb),
		shutdownChan:   make(chan bool),
		stopDbUpgrade:  stopDbUpgrade,
		networkId:      config.NetworkId,
		gasPrice:       config.GasPrice,
		lbchain-deverbase:      config.lbchain-deverbase,
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   NewBloomIndexer(chainDb, params.BloomBitsBlocks),
	}

	log.Info("Initialising lbchain-devchain protocol", "versions", ProtocolVersions, "network", config.NetworkId)

	if !config.SkipBcVersionCheck {
		bcVersion := core.GetBlockChainVersion(chainDb)
		if bcVersion != core.BlockChainVersion && bcVersion != 0 {
			return nil, fmt.Errorf("Blockchain DB version mismatch (%d / %d). Run glbchain-dev upgradedb.\n", bcVersion, core.BlockChainVersion)
		}
		core.WriteBlockChainVersion(chainDb, core.BlockChainVersion)
	}
	var (
		vmConfig    = vm.Config{EnablePreimageRecording: config.EnablePreimageRecording}
		cacheConfig = &core.CacheConfig{Disabled: config.NoPruning, TrieNodeLimit: config.TrieCache, TrieTimeLimit: config.TrieTimeout}
	)
	lbchain-dev.blockchain, err = core.NewBlockChain(chainDb, cacheConfig, lbchain-dev.chainConfig, lbchain-dev.engine, vmConfig)
	if err != nil {
		return nil, err
	}
	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		lbchain-dev.blockchain.SetHead(compat.RewindTo)
		core.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}
	lbchain-dev.bloomIndexer.Start(lbchain-dev.blockchain)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	}
	lbchain-dev.txPool = core.NewTxPool(config.TxPool, lbchain-dev.chainConfig, lbchain-dev.blockchain)

	if lbchain-dev.protocolManager, err = NewProtocolManager(lbchain-dev.chainConfig, config.SyncMode, config.NetworkId, lbchain-dev.eventMux, lbchain-dev.txPool, lbchain-dev.engine, lbchain-dev.blockchain, chainDb); err != nil {
		return nil, err
	}
	lbchain-dev.miner = miner.New(lbchain-dev, lbchain-dev.chainConfig, lbchain-dev.EventMux(), lbchain-dev.engine)
	lbchain-dev.miner.SetExtra(makeExtraData(config.ExtraData))

	lbchain-dev.ApiBackend = &lbchain-devApiBackend{lbchain-dev, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.GasPrice
	}
	lbchain-dev.ApiBackend.gpo = gasprice.NewOracle(lbchain-dev.ApiBackend, gpoParams)

	return lbchain-dev, nil
}

func makeExtraData(extra []byte) []byte {
	if len(extra) == 0 {
		// create default extradata
		extra, _ = rlp.EncodeToBytes([]interface{}{
			uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
			"glbchain-dev",
			runtime.Version(),
			runtime.GOOS,
		})
	}
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		log.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.MaximumExtraDataSize)
		extra = nil
	}
	return extra
}

// CreateDB creates the chain database.
func CreateDB(ctx *node.ServiceContext, config *Config, name string) (lbchain-devdb.Database, error) {
	db, err := ctx.OpenDatabase(name, config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return nil, err
	}
	if db, ok := db.(*lbchain-devdb.LDBDatabase); ok {
		db.Meter("lbchain-dev/db/chaindata/")
	}
	return db, nil
}

// CreateConsensusEngine creates the required type of consensus engine instance for an lbchain-devchain service
func CreateConsensusEngine(ctx *node.ServiceContext, config *ethash.Config, chainConfig *params.ChainConfig, db lbchain-devdb.Database) consensus.Engine {
	// If proof-of-authority is requested, set it up
	if chainConfig.Clique != nil {
		return clique.New(chainConfig.Clique, db)
	}
	// Otherwise assume proof-of-work
	switch {
	case config.PowMode == ethash.ModeFake:
		log.Warn("lbchain-devash used in fake mode")
		return ethash.NewFaker()
	case config.PowMode == ethash.ModeTest:
		log.Warn("lbchain-devash used in test mode")
		return ethash.NewTester()
	case config.PowMode == ethash.ModeShared:
		log.Warn("lbchain-devash used in shared mode")
		return ethash.NewShared()
	default:
		engine := ethash.New(ethash.Config{
			CacheDir:       ctx.ResolvePath(config.CacheDir),
			CachesInMem:    config.CachesInMem,
			CachesOnDisk:   config.CachesOnDisk,
			DatasetDir:     config.DatasetDir,
			DatasetsInMem:  config.DatasetsInMem,
			DatasetsOnDisk: config.DatasetsOnDisk,
		})
		engine.SetThreads(-1) // Disable CPU mining
		return engine
	}
}

// APIs returns the collection of RPC services the lbchain-devereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *lbchain-devchain) APIs() []rpc.API {
	apis := ethapi.GetAPIs(s.ApiBackend)

	// Append any APIs exposed explicitly by the consensus engine
	apis = append(apis, s.engine.APIs(s.BlockChain())...)

	// Append all the local APIs and return
	return append(apis, []rpc.API{
		{
			Namespace: "lbchain-dev",
			Version:   "1.0",
			Service:   NewPubliclbchain-devchainAPI(s),
			Public:    true,
		}, {
			Namespace: "lbchain-dev",
			Version:   "1.0",
			Service:   NewPublicMinerAPI(s),
			Public:    true,
		}, {
			Namespace: "lbchain-dev",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "miner",
			Version:   "1.0",
			Service:   NewPrivateMinerAPI(s),
			Public:    false,
		}, {
			Namespace: "lbchain-dev",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.ApiBackend, false),
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(s.chainConfig, s),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *lbchain-devchain) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *lbchain-devchain) lbchain-deverbase() (eb common.Address, err error) {
	s.lock.RLock()
	lbchain-deverbase := s.lbchain-deverbase
	s.lock.RUnlock()

	if lbchain-deverbase != (common.Address{}) {
		return lbchain-deverbase, nil
	}
	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accounts := wallets[0].Accounts(); len(accounts) > 0 {
			lbchain-deverbase := accounts[0].Address

			s.lock.Lock()
			s.lbchain-deverbase = lbchain-deverbase
			s.lock.Unlock()

			log.Info("lbchain-deverbase automatically configured", "address", lbchain-deverbase)
			return lbchain-deverbase, nil
		}
	}
	return common.Address{}, fmt.Errorf("lbchain-deverbase must be explicitly specified")
}

// set in js console via admin interface or wrapper from cli flags
func (self *lbchain-devchain) Setlbchain-deverbase(lbchain-deverbase common.Address) {
	self.lock.Lock()
	self.lbchain-deverbase = lbchain-deverbase
	self.lock.Unlock()

	self.miner.Setlbchain-deverbase(lbchain-deverbase)
}

func (s *lbchain-devchain) StartMining(local bool) error {
	eb, err := s.lbchain-deverbase()
	if err != nil {
		log.Error("Cannot start mining without lbchain-deverbase", "err", err)
		return fmt.Errorf("lbchain-deverbase missing: %v", err)
	}
	if clique, ok := s.engine.(*clique.Clique); ok {
		wallet, err := s.accountManager.Find(accounts.Account{Address: eb})
		if wallet == nil || err != nil {
			log.Error("lbchain-deverbase account unavailable locally", "err", err)
			return fmt.Errorf("signer missing: %v", err)
		}
		clique.Authorize(eb, wallet.SignHash)
	}
	if local {
		// If local (CPU) mining is started, we can disable the transaction rejection
		// mechanism introduced to speed sync times. CPU mining on mainnet is ludicrous
		// so noone will ever hit this path, whereas marking sync done on CPU mining
		// will ensure that private networks work in single miner mode too.
		atomic.StoreUint32(&s.protocolManager.acceptTxs, 1)
	}
	go s.miner.Start(eb)
	return nil
}

func (s *lbchain-devchain) StopMining()         { s.miner.Stop() }
func (s *lbchain-devchain) IsMining() bool      { return s.miner.Mining() }
func (s *lbchain-devchain) Miner() *miner.Miner { return s.miner }

func (s *lbchain-devchain) AccountManager() *accounts.Manager  { return s.accountManager }
func (s *lbchain-devchain) BlockChain() *core.BlockChain       { return s.blockchain }
func (s *lbchain-devchain) TxPool() *core.TxPool               { return s.txPool }
func (s *lbchain-devchain) EventMux() *event.TypeMux           { return s.eventMux }
func (s *lbchain-devchain) Engine() consensus.Engine           { return s.engine }
func (s *lbchain-devchain) ChainDb() lbchain-devdb.Database            { return s.chainDb }
func (s *lbchain-devchain) IsListening() bool                  { return true } // Always listening
func (s *lbchain-devchain) lbchain-devVersion() int                    { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *lbchain-devchain) NetVersion() uint64                 { return s.networkId }
func (s *lbchain-devchain) Downloader() *downloader.Downloader { return s.protocolManager.downloader }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *lbchain-devchain) Protocols() []p2p.Protocol {
	if s.lesServer == nil {
		return s.protocolManager.SubProtocols
	}
	return append(s.protocolManager.SubProtocols, s.lesServer.Protocols()...)
}

// Start implements node.Service, starting all internal goroutines needed by the
// lbchain-devchain protocol implementation.
func (s *lbchain-devchain) Start(srvr *p2p.Server) error {
	// Start the bloom bits servicing goroutines
	s.startBloomHandlers()

	// Start the RPC service
	s.netRPCService = ethapi.NewPublicNetAPI(srvr, s.NetVersion())

	// Figure out a max peers count based on the server limits
	maxPeers := srvr.MaxPeers
	if s.config.LightServ > 0 {
		if s.config.LightPeers >= srvr.MaxPeers {
			return fmt.Errorf("invalid peer config: light peer count (%d) >= total peer count (%d)", s.config.LightPeers, srvr.MaxPeers)
		}
		maxPeers -= s.config.LightPeers
	}
	// Start the networking layer and the light server if requested
	s.protocolManager.Start(maxPeers)
	if s.lesServer != nil {
		s.lesServer.Start(srvr)
	}
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// lbchain-devchain protocol.
func (s *lbchain-devchain) Stop() error {
	if s.stopDbUpgrade != nil {
		s.stopDbUpgrade()
	}
	s.bloomIndexer.Close()
	s.blockchain.Stop()
	s.protocolManager.Stop()
	if s.lesServer != nil {
		s.lesServer.Stop()
	}
	s.txPool.Stop()
	s.miner.Stop()
	s.eventMux.Stop()

	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}
