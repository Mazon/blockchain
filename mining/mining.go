package mining

import (
	"blockchain/blockchain"
	"blockchain/chaincfg"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	// maxNonce is the maximum value a nonce can be in a block header.
	maxNonce = ^uint32(0) // 2^32 - 1

	// maxExtraNonce is the maximum value an extra nonce used in a coinbase
	// transaction can be.
	//	maxExtraNonce = ^uint64(0) // 2^64 - 1

	// hpsUpdateSecs is the number of seconds to wait in between each
	// update to the hashes per second monitor.
	//	hpsUpdateSecs = 10

	// hashUpdateSec is the number of seconds each worker waits in between
	// notifying the speed monitor with how many hashes have been completed
	// while they are actively searching for a solution.  This is done to
	// reduce the amount of syncs between the workers that must be done to
	// keep track of the hashes per second.
	hashUpdateSecs = 15
)

var (
	// defaultNumWorkers is the default number of workers to use for mining
	// and is based on the number of processor cores.  This helps ensure the
	// system stays reasonably responsive under heavy load.
	defaultNumWorkers = uint32(runtime.NumCPU())
)

// Config is a descriptor containing the cpu miner configuration.
type Config struct {
	// ChainParams identifies which chain parameters the cpu miner is
	// associated with.
	ChainParams *chaincfg.Params

	// BlockTemplateGenerator identifies the instance to use in order to
	// generate block templates that the miner will attempt to solve.
	BlockTemplateGenerator *BlkTmplGenerator

	// MiningAddrs is a list of payment addresses to use for the generated
	// blocks.  Each generated block will randomly choose one of them.
	//MiningAddrs []btcutil.Address

	// ProcessBlock defines the function to call with any solved blocks.
	// It typically must run the provided block through the same set of
	// rules and handling as any other block coming from the network.
	//ProcessBlock func(*btcutil.Block, blockchain.BehaviorFlags) (bool, error)

	// ConnectedCount defines the function to use to obtain how many other
	// peers the server is connected to.  This is used by the automatic
	// persistent mining routine to determine whether or it should attempt
	// mining.  This is useful because there is no point in mining when not
	// connected to any peers since there would no be anyone to send any
	// found blocks to.
	ConnectedCount func() int32

	// IsCurrent defines the function to use to obtain whether or not the
	// block chain is current.  This is used by the automatic persistent
	// mining routine to determine whether or it should attempt mining.
	// This is useful because there is no point in mining if the chain is
	// not current since any solved blocks would be on a side chain and and
	// up orphaned anyways.
	IsCurrent func() bool
}

// BlockTemplate houses a block that has yet to be solved along with additional
// details about the fees and the number of signature operations for each
// transaction in the block.
type BlockTemplate struct {
	// Block is a block that is ready to be solved by miners.  Thus, it is
	// completely valid with the exception of satisfying the proof-of-work
	// requirement.
	Block *blockchain.Block

	// Fees contains the amount of fees each transaction in the generated
	// template pays in base units.  Since the first transaction is the
	// coinbase, the first entry (offset 0) will contain the negative of the
	// sum of the fees of all other transactions.
	//Fees []int64

	// SigOpCosts contains the number of signature operations each
	// transaction in the generated template performs.
	//SigOpCosts []int64

	// Height is the height at which the block template connects to the main
	// chain.
	Height int32

	// ValidPayAddress indicates whether or not the template coinbase pays
	// to an address or is redeemable by anyone.  See the documentation on
	// NewBlockTemplate for details on which this can be useful to generate
	// templates without a coinbase payment address.
	ValidPayAddress bool

	// WitnessCommitment is a commitment to the witness data (if any)
	// within the block. This field will only be populted once segregated
	// witness has been activated, and the block contains a transaction
	// which has witness data.
	//WitnessCommitment []byte
}

// Miner provides facilities for solving blocks (mining) using the CPU in
// a concurrency-safe manner.  It consists of two main goroutines -- a speed
// monitor and a controller for worker goroutines which generate and solve
// blocks.  The number of goroutines can be set via the SetMaxGoRoutines
// function, but the default is based on the number of processor cores in the
// system which is typically sufficient.
type Miner struct {
	sync.Mutex
	//g          *mining.BlkTmplGenerator
	g          *BlkTmplGenerator
	cfg        Config
	numWorkers uint32
	started    bool
	//	discreteMining    bool
	submitBlockLock   sync.Mutex
	wg                sync.WaitGroup
	workerWg          sync.WaitGroup
	updateNumWorkers  chan struct{}
	queryHashesPerSec chan float64
	updateHashes      chan uint64
	//	speedMonitorQuit  chan struct{}
	quit chan struct{}
}

// BlkTmplGenerator provides a type that can be used to generate block templates
// based on a given mining policy and source of transactions to choose from.
// It also houses additional state required in order to ensure the templates
// are built on top of the current best chain and adhere to the consensus rules.
type BlkTmplGenerator struct {
	//policy      *Policy
	//chainParams *chaincfg.Params
	//txSource TxSource
	chain *blockchain.Blockchain
	//timeSource blockchain.MedianTimeSource
	//sigCache   *txscript.SigCache
	//hashCache *txscript.HashCache
}

// NewBlkTmplGenerator returns a new block template generator for the given
// policy using transactions from the provided transaction source.
//
// The additional state-related fields are required in order to ensure the
// templates are built on top of the current best chain and adhere to the
// consensus rules.
/*func NewBlkTmplGenerator(policy *Policy, params *chaincfg.Params,
txSource TxSource, chain *blockchain.BlockChain,
timeSource blockchain.MedianTimeSource,
sigCache *txscript.SigCache,
hashCache *txscript.HashCache) *BlkTmplGenerator {
*/
func NewBlkTmplGenerator(policy *Policy, params *chaincfg.Params, chain *blockchain.Blockchain) *BlkTmplGenerator {
	return &BlkTmplGenerator{
		chain: chain,
	}
}

// Start begins the mining proces.
//
// This function is safe for concurrent access.
func (m *Miner) Start() {
	m.Lock()
	defer m.Unlock()

	if m.started {
		return
	}

	m.quit = make(chan struct{})
	m.wg.Add(1)
	go m.miningWorkerController()

	m.started = true
	//log.Infof("miner started")
}

// miningWorkerController launches the worker goroutines that are used to
// generate block templates and solve them.  It also provides the ability to
// dynamically adjust the number of running worker goroutines.
//
// It must be run as a goroutine.
func (m *Miner) miningWorkerController() {
	// launchWorkers groups common code to launch a specified number of
	// workers for generating blocks.
	var runningWorkers []chan struct{}
	launchWorkers := func(numWorkers uint32) {
		for i := uint32(0); i < numWorkers; i++ {
			quit := make(chan struct{})
			runningWorkers = append(runningWorkers, quit)

			m.workerWg.Add(1)
			go m.generateBlocks(quit)
		}
	}

	// Launch the current number of workers by default.
	runningWorkers = make([]chan struct{}, 0, m.numWorkers)
	launchWorkers(m.numWorkers)

out:
	for {
		select {
		// Update the number of running workers.
		case <-m.updateNumWorkers:
			// No change.
			numRunning := uint32(len(runningWorkers))
			if m.numWorkers == numRunning {
				continue
			}

			// Add new workers.
			if m.numWorkers > numRunning {
				launchWorkers(m.numWorkers - numRunning)
				continue
			}

			// Signal the most recently created goroutines to exit.
			for i := numRunning - 1; i >= m.numWorkers; i-- {
				close(runningWorkers[i])
				runningWorkers[i] = nil
				runningWorkers = runningWorkers[:i]
			}

		case <-m.quit:
			for _, quit := range runningWorkers {
				close(quit)
			}
			break out
		}
	}

	// Wait until all workers shut down to stop the speed monitor since
	// they rely on being able to send updates to it.
	m.workerWg.Wait()
	//close(m.speedMonitorQuit)
	m.wg.Done()
}

// generateBlocks is a worker that is controlled by the miningWorkerController.
// It is self contained in that it creates block templates and attempts to solve
// them while detecting when it is performing stale work and reacting
// accordingly by generating a new block template.  When a block is solved, it
// is submitted.
//
// It must be run as a goroutine.
func (m *Miner) generateBlocks(quit chan struct{}) {
	//log.Tracef("Starting generate blocks worker")

	// Start a ticker which is used to signal checks for stale work and
	// updates to the speed monitor.
	ticker := time.NewTicker(time.Second * hashUpdateSecs)
	defer ticker.Stop()
out:
	for {
		// Quit when the miner is stopped.
		select {
		case <-quit:
			break out
		default:
			// Non-blocking select to fall through
		}

		// Wait until there is a connection to at least one other peer
		// since there is no way to relay a found block or receive
		// transactions to work on when there are no connected peers.
		//	if m.cfg.ConnectedCount() == 0 {
		//			time.Sleep(time.Second)
		//			continue
		//		}

		// No point in searching for a solution before the chain is
		// synced.  Also, grab the same lock as used for block
		// submission, since the current block will be changing and
		// this would otherwise end up building a new block template on
		// a block that is in the process of becoming stale.
		m.submitBlockLock.Lock()
		curHeight := m.g.BestSnapshot().Height
		if curHeight != 0 && !m.cfg.IsCurrent() {
			m.submitBlockLock.Unlock()
			time.Sleep(time.Second)
			continue
		}

		// Choose a payment address at random.
		rand.Seed(time.Now().UnixNano())
		//payToAddr := m.cfg.MiningAddrs[rand.Intn(len(m.cfg.MiningAddrs))]
		payToAddr := "abc123"

		// Create a new block template using the available transactions
		// in the memory pool as a source of transactions to potentially
		// include in the block.
		template, err := m.g.NewBlockTemplate(payToAddr)
		m.submitBlockLock.Unlock()
		if err != nil {
			errStr := fmt.Sprintf("Failed to create new block "+
				"template: %v", err)
			_ = errStr
			//			log.Errorf(errStr)
			continue
		}

		// Attempt to solve the block.  The function will exit early
		// with false when conditions that trigger a stale block, so
		// a new block template can be generated.  When the return is
		// true a solution was found, so submit the solved block.
		if m.solveBlock(template.Block, curHeight+1, ticker, quit) {
			//			block := btcutil.NewBlock(template.Block)
			// TODO
			//m.submitBlock(block)
		}
	}

	m.workerWg.Done()
	//	log.Tracef("Generate blocks worker done")
}

// solveBlock attempts to find some combination of a nonce, extra nonce, and
// current timestamp which makes the passed block hash to a value less than the
// target difficulty.  The timestamp is updated periodically and the passed
// block is modified with all tweaks during this process.  This means that
// when the function returns true, the block is ready for submission.
//
// This function will return early with false when conditions that trigger a
// stale block such as a new block showing up or periodically when there are
// new transactions and enough time has elapsed without finding a solution.
func (m *Miner) solveBlock(msgBlock *blockchain.Block, blockHeight int32,
	ticker *time.Ticker, quit chan struct{}) bool {

	// Choose a random extra nonce offset for this block template and
	// worker.
	//enOffset, err := wire.RandomUint64()
	//if err != nil {
	//		log.Errorf("Unexpected error while generating random "+
	//			"extra nonce offset: %v", err)
	//		enOffset = 0
	//	}

	// Create some convenience variables.
	header := &msgBlock.Header
	//	targetDifficulty := ""
	//targetDifficulty := blockchain.CompactToBig(header.Bits)

	// Initial state.
	//	lastGenerated := time.Now()
	//lastTxUpdate := m.g.TxSource().LastUpdated()
	hashesCompleted := uint64(0)

	// Note that the entire extra nonce range is iterated and the offset is
	// added relying on the fact that overflow will wrap around 0 as
	// provided by the Go spec.
	//	for extraNonce := uint64(0); extraNonce < maxExtraNonce; extraNonce++ {
	// Update the extra nonce in the block template with the
	// new value by regenerating the coinbase script and
	// setting the merkle root to the new value.
	//		m.g.UpdateExtraNonce(msgBlock, blockHeight, extraNonce+enOffset)

	// Search through the entire nonce range for a solution while
	// periodically checking for early quit and stale block
	// conditions along with updates to the speed monitor.
	for i := uint32(0); i <= maxNonce; i++ {
		select {
		case <-quit:
			return false

		case <-ticker.C:
			m.updateHashes <- hashesCompleted
			hashesCompleted = 0

			// The current block is stale if the best block
			// has changed.
			best := m.g.BestSnapshot()
			_ = best
			//if !header.PrevHash.IsEqual(&best.Hash) {
			//	return false
			//}

			// The current block is stale if the memory pool
			// has been updated since the block template was
			// generated and it has been at least one
			// minute.
			//if lastTxUpdate != m.g.TxSource().LastUpdated() &&
			//	time.Now().After(lastGenerated.Add(time.Minute)) {

			//	return false
			//}

			//m.g.UpdateBlockTime(msgBlock)

		default:
			// Non-blocking select to fall through
		}

		// Update the nonce and hash the block header.  Each
		// hash is a sha256, so
		// increment the number of hashes completed for each
		// attempt accordingly.
		header.Nonce = i
		hash := header.BlockHash()
		_ = hash
		hashesCompleted += 1

		// The block is solved when the new block hash is less
		// than the target difficulty.  Yay!
		//if blockchain.HashToBig(&hash).Cmp(targetDifficulty) <= 0 {
		//		m.updateHashes <- hashesCompleted
		//	return true
		//	}
	}
	//}

	return false
}

// New returns a new instance of a miner for the provided configuration.
// Use Start to begin the mining process.  See the documentation for Miner
// type for more details.
func New(cfg *Config) *Miner {
	return &Miner{
		g:                 cfg.BlockTemplateGenerator,
		cfg:               *cfg,
		numWorkers:        defaultNumWorkers,
		updateNumWorkers:  make(chan struct{}),
		queryHashesPerSec: make(chan float64),
		updateHashes:      make(chan uint64),
	}
}
