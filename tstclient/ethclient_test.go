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

package lbchain-devclient

import "github.com/lbchain-devchain/go-lbchain-dev"

// Verify that Client implements the lbchain-devereum interfaces.
var (
	_ = lbchain-devereum.ChainReader(&Client{})
	_ = lbchain-devereum.TransactionReader(&Client{})
	_ = lbchain-devereum.ChainStateReader(&Client{})
	_ = lbchain-devereum.ChainSyncReader(&Client{})
	_ = lbchain-devereum.ContractCaller(&Client{})
	_ = lbchain-devereum.GasEstimator(&Client{})
	_ = lbchain-devereum.GasPricer(&Client{})
	_ = lbchain-devereum.LogFilterer(&Client{})
	_ = lbchain-devereum.PendingStateReader(&Client{})
	// _ = lbchain-devereum.PendingStateEventer(&Client{})
	_ = lbchain-devereum.PendingContractCaller(&Client{})
)
