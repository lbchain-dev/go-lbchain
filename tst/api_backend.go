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

package lbchain-dev

import (
	"context"
	"math/big"

	"github.com/lbchain-devchain/go-lbchain-dev/accounts"
	"github.com/lbchain-devchain/go-lbchain-dev/common"
	"github.com/lbchain-devchain/go-lbchain-dev/common/math"
	"github.com/lbchain-devchain/go-lbchain-dev/core"
	"github.com/lbchain-devchain/go-lbchain-dev/core/bloombits"
	"github.com/lbchain-devchain/go-lbchain-dev/core/state"
	"github.com/lbchain-devchain/go-lbchain-dev/core/types"
	"github.com/lbchain-devchain/go-lbchain-dev/core/vm"
	"github.com/lbchain-devchain/go-lbchain-dev/lbchain-dev/downloader"
	"github.com/lbchain-devchain/go-lbchain-dev/lbchain-dev/gasprice"
	"github.com/lbchain-devchain/go-lbchain-dev/lbchain-devdb"
	"github.com/lbchain-devchain/go-lbchain-dev/event"
	"github.com/lbchain-devchain/go-lbchain-dev/params"
	"github.com/lbchain-devchain/go-lbchain-dev/rpc"
)

// lbchain-devApiBackend implements ethapi.Backend for full nodes
type lbchain-devApiBackend struct {
	lbchain-dev *lbchain-devchain
	gpo *gasprice.Oracle
}

func (b *lbchain-devApiBackend) ChainConfig() *params.ChainConfig {
	return b.lbchain-dev.chainConfig
}

func (b *lbchain-devApiBackend) CurrentBlock() *types.Block {
	return b.lbchain-dev.blockchain.CurrentBlock()
}

func (b *lbchain-devApiBackend) SetHead(number uint64) {
	b.lbchain-dev.protocolManager.downloader.Cancel()
	b.lbchain-dev.blockchain.SetHead(number)
}

func (b *lbchain-devApiBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.lbchain-dev.miner.PendingBlock()
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.lbchain-dev.blockchain.CurrentBlock().Header(), nil
	}
	return b.lbchain-dev.blockchain.GetHeaderByNumber(uint64(blockNr)), nil
}

func (b *lbchain-devApiBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.lbchain-dev.miner.PendingBlock()
		return block, nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.lbchain-dev.blockchain.CurrentBlock(), nil
	}
	return b.lbchain-dev.blockchain.GetBlockByNumber(uint64(blockNr)), nil
}

func (b *lbchain-devApiBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Pending state is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block, state := b.lbchain-dev.miner.Pending()
		return state, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	stateDb, err := b.lbchain-dev.BlockChain().StateAt(header.Root)
	return stateDb, header, err
}

func (b *lbchain-devApiBackend) GetBlock(ctx context.Context, blockHash common.Hash) (*types.Block, error) {
	return b.lbchain-dev.blockchain.GetBlockByHash(blockHash), nil
}

func (b *lbchain-devApiBackend) GetReceipts(ctx context.Context, blockHash common.Hash) (types.Receipts, error) {
	return core.GetBlockReceipts(b.lbchain-dev.chainDb, blockHash, core.GetBlockNumber(b.lbchain-dev.chainDb, blockHash)), nil
}

func (b *lbchain-devApiBackend) GetLogs(ctx context.Context, blockHash common.Hash) ([][]*types.Log, error) {
	receipts := core.GetBlockReceipts(b.lbchain-dev.chainDb, blockHash, core.GetBlockNumber(b.lbchain-dev.chainDb, blockHash))
	if receipts == nil {
		return nil, nil
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

func (b *lbchain-devApiBackend) GetTd(blockHash common.Hash) *big.Int {
	return b.lbchain-dev.blockchain.GetTdByHash(blockHash)
}

func (b *lbchain-devApiBackend) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	vmError := func() error { return nil }

	context := core.NewEVMContext(msg, header, b.lbchain-dev.BlockChain(), nil)
	return vm.NewEVM(context, state, b.lbchain-dev.chainConfig, vmCfg), vmError, nil
}

func (b *lbchain-devApiBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.lbchain-dev.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *lbchain-devApiBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.lbchain-dev.BlockChain().SubscribeChainEvent(ch)
}

func (b *lbchain-devApiBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.lbchain-dev.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *lbchain-devApiBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.lbchain-dev.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *lbchain-devApiBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.lbchain-dev.BlockChain().SubscribeLogsEvent(ch)
}

func (b *lbchain-devApiBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.lbchain-dev.txPool.AddLocal(signedTx)
}

func (b *lbchain-devApiBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.lbchain-dev.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *lbchain-devApiBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.lbchain-dev.txPool.Get(hash)
}

func (b *lbchain-devApiBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.lbchain-dev.txPool.State().GetNonce(addr), nil
}

func (b *lbchain-devApiBackend) Stats() (pending int, queued int) {
	return b.lbchain-dev.txPool.Stats()
}

func (b *lbchain-devApiBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.lbchain-dev.TxPool().Content()
}

func (b *lbchain-devApiBackend) SubscribeTxPreEvent(ch chan<- core.TxPreEvent) event.Subscription {
	return b.lbchain-dev.TxPool().SubscribeTxPreEvent(ch)
}

func (b *lbchain-devApiBackend) Downloader() *downloader.Downloader {
	return b.lbchain-dev.Downloader()
}

func (b *lbchain-devApiBackend) ProtocolVersion() int {
	return b.lbchain-dev.lbchain-devVersion()
}

func (b *lbchain-devApiBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *lbchain-devApiBackend) ChainDb() lbchain-devdb.Database {
	return b.lbchain-dev.ChainDb()
}

func (b *lbchain-devApiBackend) EventMux() *event.TypeMux {
	return b.lbchain-dev.EventMux()
}

func (b *lbchain-devApiBackend) AccountManager() *accounts.Manager {
	return b.lbchain-dev.AccountManager()
}

func (b *lbchain-devApiBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.lbchain-dev.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *lbchain-devApiBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.lbchain-dev.bloomRequests)
	}
}
