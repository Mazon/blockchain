// Copyright (c) 2014-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mining

const (
	// UnminedHeight is the height used for the "block" height field of the
	// contextual transaction information provided in a transaction store
	// when it has not yet been mined into a block.
	UnminedHeight = 0x7fffffff
)

// Policy houses the policy (configuration parameters) which is used to control
// the generation of block templates.  See the documentation for
// NewBlockTemplate for more details on each of these parameters are used.
type Policy struct {
	// BlockMinWeight is the minimum block weight to be used when
	// generating a block template.
	BlockMinWeight uint32

	// BlockMaxWeight is the maximum block weight to be used when
	// generating a block template.
	BlockMaxWeight uint32

	// BlockMinWeight is the minimum block size to be used when generating
	// a block template.
	BlockMinSize uint32

	// BlockMaxSize is the maximum block size to be used when generating a
	// block template.
	BlockMaxSize uint32

	// BlockPrioritySize is the size in bytes for high-priority / low-fee
	// transactions to be used when generating a block template.
	BlockPrioritySize uint32

	// TxMinFreeFee is the minimum fee in Satoshi/1000 bytes that is
	// required for a transaction to be treated as free for mining purposes
	// (block template generation).
	//TxMinFreeFee btcutil.Amount
}
