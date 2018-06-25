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

// Contains the metrics collected by the downloader.

package downloader

import (
	"github.com/lbchain-devchain/go-lbchain-dev/metrics"
)

var (
	headerInMeter      = metrics.NewRegisteredMeter("lbchain-dev/downloader/headers/in", nil)
	headerReqTimer     = metrics.NewRegisteredTimer("lbchain-dev/downloader/headers/req", nil)
	headerDropMeter    = metrics.NewRegisteredMeter("lbchain-dev/downloader/headers/drop", nil)
	headerTimeoutMeter = metrics.NewRegisteredMeter("lbchain-dev/downloader/headers/timeout", nil)

	bodyInMeter      = metrics.NewRegisteredMeter("lbchain-dev/downloader/bodies/in", nil)
	bodyReqTimer     = metrics.NewRegisteredTimer("lbchain-dev/downloader/bodies/req", nil)
	bodyDropMeter    = metrics.NewRegisteredMeter("lbchain-dev/downloader/bodies/drop", nil)
	bodyTimeoutMeter = metrics.NewRegisteredMeter("lbchain-dev/downloader/bodies/timeout", nil)

	receiptInMeter      = metrics.NewRegisteredMeter("lbchain-dev/downloader/receipts/in", nil)
	receiptReqTimer     = metrics.NewRegisteredTimer("lbchain-dev/downloader/receipts/req", nil)
	receiptDropMeter    = metrics.NewRegisteredMeter("lbchain-dev/downloader/receipts/drop", nil)
	receiptTimeoutMeter = metrics.NewRegisteredMeter("lbchain-dev/downloader/receipts/timeout", nil)

	stateInMeter   = metrics.NewRegisteredMeter("lbchain-dev/downloader/states/in", nil)
	stateDropMeter = metrics.NewRegisteredMeter("lbchain-dev/downloader/states/drop", nil)
)
