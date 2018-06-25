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

package les

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
	"github.com/lbchain-devchain/go-lbchain-dev/light"
	"github.com/lbchain-devchain/go-lbchain-dev/params"
	"github.com/lbchain-devchain/go-lbchain-dev/rpc"
)

type LesApiBackend struct {
	lbchain-dev *Lightlbchain-devchain
	gpo *gasprice.Oracle
}

func (b *LesApiBackend) ChainConfig() *params.ChainConfig {
	return b.lbchain-dev.chainConfig
}

func (b *LesApiBackend) CurrentBlock() *types.Block {
	return types.NewBlockWithHeader(b.lbchain-dev.BlockChain().CurrentHeader())
}

func (b *LesApiBackend) SetHead(number uint64) {
	b.lbchain-dev.protocolManager.downloader.Cancel()
	b.lbchain-dev.blockchain.SetHead(number)
}

func (b *LesApiBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	if blockNr == rpc.LatestBlockNumber || blockNr == rpc.PendingBlockNumber {
		return b.lbchain-dev.blockchain.CurrentHeader(), nil
	}

	return b.lbchain-dev.blockchain.GetHeaderByNumberOdr(ctx, uint64(blockNr))
}

func (b *LesApiBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, err
	}
	return b.GetBlock(ctx, header.Hash())
}

func (b *LesApiBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	return light.NewState(ctx, header, b.lbchain-dev.odr), header, nil
}

func (b *LesApiBackend) GetBlock(ctx context.Context, blockHash common.Hash) (*types.Block, error) {
	return b.lbchain-dev.blockchain.GetBlockByHash(ctx, blockHash)
}

func (b *LesApiBackend) GetReceipts(ctx context.Context, blockHash common.Hash) (types.Receipts, error) {
	return light.GetBlockReceipts(ctx, b.lbchain-dev.odr, blockHash, core.GetBlockNumber(b.lbchain-dev.chainDb, blockHash))
}

func (b *LesApiBackend) GetLogs(ctx context.Context, blockHash common.Hash) ([][]*types.Log, error) {
	return light.GetBlockLogs(ctx, b.lbchain-dev.odr, blockHash, core.GetBlockNumber(b.lbchain-dev.chainDb, blockHash))
}

func (b *LesApiBackend) GetTd(blockHash common.Hash) *big.Int {
	return b.lbchain-dev.blockchain.GetTdByHash(blockHash)
}

func (b *LesApiBackend) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	context := core.NewEVMContext(msg, header, b.lbchain-dev.blockchain, nil)
	return vm.NewEVM(context, state, b.lbchain-dev.chainConfig, vmCfg), state.Error, nil
}

func (b *LesApiBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.lbchain-dev.txPool.Add(ctx, signedTx)
}

func (b *LesApiBackend) RemoveTx(txHash common.Hash) {
	b.lbchain-dev.txPool.RemoveTx(txHash)
}

func (b *LesApiBackend) GetPoolTransactions() (types.Transactions, error) {
	return b.lbchain-dev.txPool.GetTransactions()
}

func (b *LesApiBackend) GetPoolTransaction(txHash common.Hash) *types.Transaction {
	return b.lbchain-dev.txPool.GetTransaction(txHash)
}

func (b *LesApiBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.lbchain-dev.txPool.GetNonce(ctx, addr)
}

func (b *LesApiBackend) Stats() (pending int, queued int) {
	return b.lbchain-dev.txPool.Stats(), 0
}

func (b *LesApiBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.lbchain-dev.txPool.Content()
}

func (b *LesApiBackend) SubscribeTxPreEvent(ch chan<- core.TxPreEvent) event.Subscription {
	return b.lbchain-dev.txPool.SubscribeTxPreEvent(ch)
}

func (b *LesApiBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.lbchain-dev.blockchain.SubscribeChainEvent(ch)
}

func (b *LesApiBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.lbchain-dev.blockchain.SubscribeChainHeadEvent(ch)
}

func (b *LesApiBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.lbchain-dev.blockchain.SubscribeChainSideEvent(ch)
}

func (b *LesApiBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.lbchain-dev.blockchain.SubscribeLogsEvent(ch)
}

func (b *LesApiBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.lbchain-dev.blockchain.SubscribeRemovedLogsEvent(ch)
}

func (b *LesApiBackend) Downloader() *downloader.Downloader {
	return b.lbchain-dev.Downloader()
}

func (b *LesApiBackend) ProtocolVersion() int {
	return b.lbchain-dev.LesVersion() + 10000
}

func (b *LesApiBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *LesApiBackend) ChainDb() lbchain-devdb.Database {
	return b.lbchain-dev.chainDb
}

func (b *LesApiBackend) EventMux() *event.TypeMux {
	return b.lbchain-dev.eventMux
}

func (b *LesApiBackend) AccountManager() *accounts.Manager {
	return b.lbchain-dev.accountManager
}

func (b *LesApiBackend) BloomStatus() (uint64, uint64) {
	if b.lbchain-dev.bloomIndexer == nil {
		return 0, 0
	}
	sections, _, _ := b.lbchain-dev.bloomIndexer.Sections()
	return light.BloomTrieFrequency, sections
}

func (b *LesApiBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.lbchain-dev.bloomRequests)
	}
}
